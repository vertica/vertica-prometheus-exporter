package vertica_prometheus_exporter

// (c) Copyright [2018-2022] Micro Focus or one of its affiliates.
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// 
// MIT license brought forward from the sql-exporter repo by burningalchemist
// 
// MIT License
// 
// Copyright (c) 2017 Alin Sinpalean
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTIO

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
	"github.com/vertica/vertica-prometheus-exporter/config"
	"github.com/vertica/vertica-prometheus-exporter/errors"
)

// Collector is a self-contained group of SQL queries and metric families to collect from a specific database. It is
// conceptually similar to a prometheus.Collector.
type Collector interface {
	// Collect is the equivalent of prometheus.Collector.Collect() but takes a context to run in and a database to run on.
	Collect(context.Context, *sql.DB, chan<- Metric)
}

// collector implements Collector. It wraps a collection of queries, metrics and the database to collect them from.
type collector struct {
	config     *config.CollectorConfig
	queries    []*Query
	logContext string
}

// NewCollector returns a new Collector with the given configuration and database. The metrics it creates will all have
// the provided const labels applied.
func NewCollector(logContext string, cc *config.CollectorConfig, constLabels []*dto.LabelPair) (Collector, errors.WithContext) {
	logContext = fmt.Sprintf("%s, collector=%q", logContext, cc.Name)

	// Leading comma appears when target name is undefined, which is a side-effect of running in single target mode.
	// Let's trim to avoid confusions.
	if strings.HasPrefix(logContext, ",") {
		logContext = strings.TrimLeft(logContext, ", ")
	}

	// Maps each query to the list of metric families it populates.
	queryMFs := make(map[*config.QueryConfig][]*MetricFamily, len(cc.Metrics))

	// Instantiate metric families.
	for _, mc := range cc.Metrics {
		mf, err := NewMetricFamily(logContext, mc, constLabels)
		if err != nil {
			return nil, err
		}
		mfs, found := queryMFs[mc.Query()]
		if !found {
			mfs = make([]*MetricFamily, 0, 2)
		}
		queryMFs[mc.Query()] = append(mfs, mf)
	}

	// Instantiate queries.
	queries := make([]*Query, 0, len(cc.Metrics))
	for qc, mfs := range queryMFs {
		q, err := NewQuery(logContext, qc, mfs...)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}

	c := collector{
		config:     cc,
		queries:    queries,
		logContext: logContext,
	}
	if c.config.MinInterval > 0 {
		log.Infof("[%s] Non-zero min_interval (%s), using cached collector.", logContext, c.config.MinInterval)
		return newCachingCollector(&c), nil
	}
	return &c, nil
}

// Collect implements Collector.
func (c *collector) Collect(ctx context.Context, conn *sql.DB, ch chan<- Metric) {
	var wg sync.WaitGroup
	wg.Add(len(c.queries))
	for _, q := range c.queries {
		go func(q *Query) {
			defer wg.Done()
			q.Collect(ctx, conn, ch)
		}(q)
	}
	// Only return once all queries have been processed
	wg.Wait()
}

// newCachingCollector returns a new Collector wrapping the provided raw Collector.
func newCachingCollector(rawColl *collector) Collector {
	cc := &cachingCollector{
		rawColl:     rawColl,
		minInterval: time.Duration(rawColl.config.MinInterval),
		cacheSem:    make(chan time.Time, 1),
	}
	cc.cacheSem <- time.Time{}
	return cc
}

// Collector with a cache for collected metrics. Only used when min_interval is non-zero.
type cachingCollector struct {
	// Underlying collector, which is being cached.
	rawColl *collector
	// Convenience copy of rawColl.config.MinInterval.
	minInterval time.Duration

	// Used as a non=blocking semaphore protecting the cache. The value in the channel is the time of the cached metrics.
	cacheSem chan time.Time
	// Metrics saved from the last Collect() call.
	cache []Metric
}

// Collect implements Collector.
func (cc *cachingCollector) Collect(ctx context.Context, conn *sql.DB, ch chan<- Metric) {
	if ctx.Err() != nil {
		ch <- NewInvalidMetric(errors.Wrap(cc.rawColl.logContext, ctx.Err()))
		return
	}

	collTime := time.Now()
	select {
	case cacheTime := <-cc.cacheSem:
		// Have the lock.
		if age := collTime.Sub(cacheTime); age > cc.minInterval {
			// Cache contents are older than minInterval, collect fresh metrics, cache them and pipe them through.
			log.Infof("[%s] Collecting fresh metrics: min_interval=%.3fs cache_age=%.3fs",
				cc.rawColl.logContext, cc.minInterval.Seconds(), age.Seconds())
			cacheChan := make(chan Metric, capMetricChan)
			cc.cache = make([]Metric, 0, len(cc.cache))
			go func() {
				cc.rawColl.Collect(ctx, conn, cacheChan)
				close(cacheChan)
			}()
			for metric := range cacheChan {
				cc.cache = append(cc.cache, metric)
				ch <- metric
			}
			cacheTime = collTime
		} else {
			log.Infof("[%s] Returning cached metrics: min_interval=%.3fs cache_age=%.3fs",
				cc.rawColl.logContext, cc.minInterval.Seconds(), age.Seconds())
			for _, metric := range cc.cache {
				ch <- metric
			}
		}
		// Always replace the value in the semaphore channel.
		cc.cacheSem <- cacheTime

	case <-ctx.Done():
		// Context closed, record an error and return
		// TODO: increment an error counter
		ch <- NewInvalidMetric(errors.Wrap(cc.rawColl.logContext, ctx.Err()))
	}
}
