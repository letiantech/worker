//MIT License
//
//Copyright (c) 2019 letiantech
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package analyzer

import (
	"reflect"
	"time"
)

const (
	analyzerDefaultDataSize = 1000
	analyzerDefaultWindow   = 5 * time.Second
)

type Scorer func(interface{}) *Result

type Result struct {
	Number       int64
	TotalScore   int64
	AverageScore int64
}

type Params struct {
	Type             reflect.Type
	Scorer           Scorer
	Window           time.Duration
	Receiver         AnalyzeReceiver
	ReceiverInterval time.Duration
}

type AnalyzeReceiver func(*Result)

type AnalyzeCreator func()

type Analyzer interface {
	AddData(data interface{}) error
	Analyze() *Result
}

type analyzerData struct {
	data interface{}
	tm   int64
}

type analyzer struct {
	params     *Params
	lastResult Result
	data       []*analyzerData
	dataStart  int
	dataEnd    int
}

func NewAnalyzer(params Params) Analyzer {
	a := &analyzer{
		params:    &params,
		data:      make([]*analyzerData, 0, analyzerDefaultDataSize),
		dataStart: 0,
		dataEnd:   0,
	}

	return a
}

func (a *analyzer) AddData(data interface{}) error {

	return nil
}

func (a *analyzer) Analyze() *Result {
	return nil
}
