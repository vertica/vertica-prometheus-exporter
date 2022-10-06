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
	"fmt"
	"sort"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/vertica/vertica-prometheus-exporter/config"
	"github.com/vertica/vertica-prometheus-exporter/errors"
	"google.golang.org/protobuf/proto"
)

// MetricDesc is a descriptor for a family of metrics, sharing the same name, help, labes, type.
type MetricDesc interface {
	Name() string
	Help() string
	ValueType() prometheus.ValueType
	ConstLabels() []*dto.LabelPair
	Labels() []string
	LogContext() string
}

//
// MetricFamily
//

// MetricFamily implements MetricDesc for SQL metrics, with logic for populating its labels and values from sql.Rows.
type MetricFamily struct {
	config      *config.MetricConfig
	constLabels []*dto.LabelPair
	labels      []string
	logContext  string
}

// NewMetricFamily creates a new MetricFamily with the given metric config and const labels (e.g. job and instance).
func NewMetricFamily(logContext string, mc *config.MetricConfig, constLabels []*dto.LabelPair) (*MetricFamily, errors.WithContext) {
	logContext = fmt.Sprintf("%s, metric=%q", logContext, mc.Name)

	if len(mc.Values) == 0 {
		return nil, errors.New(logContext, "no value column defined")
	}
	if len(mc.Values) > 1 && mc.ValueLabel == "" {
		return nil, errors.New(logContext, "multiple values but no value label")
	}
	if len(mc.KeyLabels) > config.MaxInt32 {
		return nil, errors.New(logContext, "key_labels list is too large")
	}

	labels := make([]string, 0, len(mc.KeyLabels)+1)
	labels = append(labels, mc.KeyLabels...)
	if mc.ValueLabel != "" {
		labels = append(labels, mc.ValueLabel)
	}

	// Create a copy of original slice to avoid modifying constLabels
	sortedLabels := append(constLabels[:0:0], constLabels...)

	for k, v := range mc.StaticLabels {
		sortedLabels = append(sortedLabels, &dto.LabelPair{
			Name:  proto.String(k),
			Value: proto.String(v),
		})
	}
	sort.Sort(labelPairSorter(sortedLabels))

	return &MetricFamily{
		config:      mc,
		constLabels: sortedLabels,
		labels:      labels,
		logContext:  logContext,
	}, nil
}

// Collect is the equivalent of prometheus.Collector.Collect() but takes a Query output map to populate values from.
func (mf MetricFamily) Collect(row map[string]interface{}, ch chan<- Metric) {
	labelValues := make([]string, len(mf.labels))
	for i, label := range mf.config.KeyLabels {
		labelValues[i] = row[label].(string)
	}
	for _, v := range mf.config.Values {
		if mf.config.ValueLabel != "" {
			labelValues[len(labelValues)-1] = v
		}
		value := row[v].(float64)
		ch <- NewMetric(&mf, value, labelValues...)
	}
}

// Name implements MetricDesc.
func (mf MetricFamily) Name() string {
	return mf.config.Name
}

// Help implements MetricDesc.
func (mf MetricFamily) Help() string {
	return mf.config.Help
}

// ValueType implements MetricDesc.
func (mf MetricFamily) ValueType() prometheus.ValueType {
	return mf.config.ValueType()
}

// ConstLabels implements MetricDesc.
func (mf MetricFamily) ConstLabels() []*dto.LabelPair {
	return mf.constLabels
}

// Labels implements MetricDesc.
func (mf MetricFamily) Labels() []string {
	return mf.labels
}

// LogContext implements MetricDesc.
func (mf MetricFamily) LogContext() string {
	return mf.logContext
}

//
// automaticMetricDesc
//

// automaticMetric is a MetricDesc for automatically generated metrics (e.g. `up` and `scrape_duration`).
type automaticMetricDesc struct {
	name        string
	help        string
	valueType   prometheus.ValueType
	labels      []string
	constLabels []*dto.LabelPair
	logContext  string
}

// NewAutomaticMetricDesc creates a MetricDesc for automatically generated metrics.
func NewAutomaticMetricDesc(
	logContext, name, help string, valueType prometheus.ValueType, constLabels []*dto.LabelPair, labels ...string) MetricDesc {
	return &automaticMetricDesc{
		name:        name,
		help:        help,
		valueType:   valueType,
		constLabels: constLabels,
		labels:      labels,
		logContext:  logContext,
	}
}

// Name implements MetricDesc.
func (a automaticMetricDesc) Name() string {
	return a.name
}

// Help implements MetricDesc.
func (a automaticMetricDesc) Help() string {
	return a.help
}

// ValueType implements MetricDesc.
func (a automaticMetricDesc) ValueType() prometheus.ValueType {
	return a.valueType
}

// ConstLabels implements MetricDesc.
func (a automaticMetricDesc) ConstLabels() []*dto.LabelPair {
	return a.constLabels
}

// Labels implements MetricDesc.
func (a automaticMetricDesc) Labels() []string {
	return a.labels
}

// LogContext implements MetricDesc.
func (a automaticMetricDesc) LogContext() string {
	return a.logContext
}

//
// Metric
//

// A Metric models a single sample value with its meta data being exported to Prometheus.
type Metric interface {
	Desc() MetricDesc
	Write(out *dto.Metric) errors.WithContext
}

// NewMetric returns a metric with one fixed value that cannot be changed.
//
// NewMetric panics if the length of labelValues is not consistent with desc.labels().
func NewMetric(desc MetricDesc, value float64, labelValues ...string) Metric {
	if len(desc.Labels()) != len(labelValues) {
		panic(fmt.Sprintf("[%s] expected %d labels, got %d", desc.LogContext(), len(desc.Labels()), len(labelValues)))
	}
	return &constMetric{
		desc:       desc,
		val:        value,
		labelPairs: makeLabelPairs(desc, labelValues),
	}
}

// constMetric is a metric with one fixed value that cannot be changed.
type constMetric struct {
	desc       MetricDesc
	val        float64
	labelPairs []*dto.LabelPair
}

// Desc implements Metric.
func (m *constMetric) Desc() MetricDesc {
	return m.desc
}

// Write implements Metric.
func (m *constMetric) Write(out *dto.Metric) errors.WithContext {
	out.Label = m.labelPairs
	switch t := m.desc.ValueType(); t {
	case prometheus.CounterValue:
		out.Counter = &dto.Counter{Value: proto.Float64(m.val)}
	case prometheus.GaugeValue:
		out.Gauge = &dto.Gauge{Value: proto.Float64(m.val)}
	default:
		return errors.Errorf(m.desc.LogContext(), "encountered unknown type %v", t)
	}
	return nil
}

func makeLabelPairs(desc MetricDesc, labelValues []string) []*dto.LabelPair {
	labels := desc.Labels()
	constLabels := desc.ConstLabels()

	totalLen := len(labels) + len(constLabels)
	if totalLen == 0 {
		// Super fast path.
		return nil
	}
	if len(labels) == 0 {
		// Moderately fast path.
		return constLabels
	}
	labelPairs := make([]*dto.LabelPair, 0, totalLen)
	for i, label := range labels {
		labelPairs = append(labelPairs, &dto.LabelPair{
			Name:  proto.String(label),
			Value: proto.String(labelValues[i]),
		})
	}
	labelPairs = append(labelPairs, constLabels...)
	sort.Sort(labelPairSorter(labelPairs))
	return labelPairs
}

// labelPairSorter implements sort.Interface.
// It provides a sortable version of a slice of dto.LabelPair pointers.

type labelPairSorter []*dto.LabelPair

func (s labelPairSorter) Len() int {
	return len(s)
}

func (s labelPairSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s labelPairSorter) Less(i, j int) bool {
	return s[i].GetName() < s[j].GetName()
}

type invalidMetric struct {
	err errors.WithContext
}

// NewInvalidMetric returns a metric whose Write method always returns the provided error.
func NewInvalidMetric(err errors.WithContext) Metric {
	return invalidMetric{err}
}

func (m invalidMetric) Desc() MetricDesc { return nil }

func (m invalidMetric) Write(*dto.Metric) errors.WithContext { return m.err }
