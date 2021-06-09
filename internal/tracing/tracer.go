// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracing

import (
	"context"
	"runtime"
	"strings"

	"github.com/opentracing/opentracing-go"
	"k8s.io/apimachinery/pkg/types"
)

var mod string

type SpanOptions struct {
	operationName  string
	customResource *types.NamespacedName
}

func (o *SpanOptions) OperationName() string {
	if o.operationName == "" {
		pc, _, _, _ := runtime.Caller(2)
		details := runtime.FuncForPC(pc)
		name := details.Name()
		return strings.Replace(name, mod, "", 1)
	}
	return o.operationName
}

type SpanOptionFunc func(*SpanOptions)

func WithOperationName(operation string) SpanOptionFunc {
	return func(o *SpanOptions) {
		o.operationName = operation
	}
}

func WithCustomResource(cr types.NamespacedName) SpanOptionFunc {
	return func(o *SpanOptions) {
		o.customResource = &cr
	}
}

func StartSpanFromContext(ctx context.Context, options ...SpanOptionFunc) (*Span, context.Context) {
	opt := new(SpanOptions)

	for _, fn := range options {
		if fn == nil {
			continue
		}
		fn(opt)
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, opt.OperationName())
	if opt.customResource != nil {
		span.SetTag("kubernetes.resource", opt.customResource.String())
	}
	return &Span{Span: span}, ctx
}

func SpanFromContext(ctx context.Context) *Span {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}

	return &Span{Span: span}
}
