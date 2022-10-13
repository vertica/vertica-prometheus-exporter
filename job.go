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
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vertica/vertica-prometheus-exporter/config"
	"github.com/vertica/vertica-prometheus-exporter/errors"
)

// Job is a collection of targets with the same collectors applied.
type Job interface {
	Targets() []Target
}

// job implements Job. It wraps the corresponding JobConfig and a set of Targets.
type job struct {
	config     *config.JobConfig
	targets    []Target
	logContext string
}

// NewJob returns a new Job with the given configuration.
func NewJob(jc *config.JobConfig, gc *config.GlobalConfig) (Job, errors.WithContext) {
	j := job{
		config:     jc,
		targets:    make([]Target, 0, 10),
		logContext: fmt.Sprintf("job=%q", jc.Name),
	}

	for _, sc := range jc.StaticConfigs {
		for tname, dsn := range sc.Targets {
			constLabels := prometheus.Labels{
				"job":      jc.Name,
				"instance": tname,
			}
			for name, value := range sc.Labels {
				// Shouldn't happen as there are sanity checks in config, but check nonetheless.
				if _, found := constLabels[name]; found {
					return nil, errors.Errorf(j.logContext, "duplicate label %q", name)
				}
				constLabels[name] = value
			}
			t, err := NewTarget(j.logContext, tname, string(dsn), jc.Collectors(), constLabels, gc)
			if err != nil {
				return nil, err
			}
			j.targets = append(j.targets, t)
		}
	}

	return &j, nil
}

func (j *job) Targets() []Target {
	return j.targets
}
