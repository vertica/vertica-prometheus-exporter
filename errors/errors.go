package errors

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
)

// WithContext is an error associated with a logging context string (e.g. `job="foo", instance="bar"`). It is formatted
// as:
//
//	fmt.Sprintf("[%s] %s", Context(), RawError())
type WithContext interface {
	error

	Context() string
	RawError() string
}

// withContext implements WithContext.
type withContext struct {
	context string
	err     string
}

// New creates a new WithContext.
func New(context, err string) WithContext {
	return &withContext{context, err}
}

// Errorf formats according to a format specifier and returns a new WithContext.
func Errorf(context, format string, a ...interface{}) WithContext {
	return &withContext{context, fmt.Sprintf(format, a...)}
}

// Wrap returns a WithContext wrapping err. If err is nil, it returns nil. If err is a WithContext, it is returned
// unchanged.
func Wrap(context string, err error) WithContext {
	if err == nil {
		return nil
	}
	if w, ok := err.(WithContext); ok {
		return w
	}
	return &withContext{context, err.Error()}
}

// Wrapf returns a WithContext that prepends a formatted message to err.Error(). If err is nil, it returns nil. If err
// is a WithContext, the returned WithContext will have the message prepended but the same context as err (presumed to
// be more specific).
func Wrapf(context string, err error, format string, a ...interface{}) WithContext {
	if err == nil {
		return nil
	}
	prefix := format
	if len(a) > 0 {
		prefix = fmt.Sprintf(format, a...)
	}
	if w, ok := err.(WithContext); ok {
		return &withContext{w.Context(), prefix + ": " + w.RawError()}
	}
	return &withContext{context, prefix + ": " + err.Error()}
}

// Error implements error.
func (w *withContext) Error() string {
	if len(w.context) == 0 {
		return w.err
	}
	return "[" + w.context + "] " + w.err
}

// Context implements WithContext.
func (w *withContext) Context() string {
	return w.context
}

// RawError implements WithContext.
func (w *withContext) RawError() string {
	return w.err
}
