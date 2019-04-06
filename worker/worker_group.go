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

package worker

import (
	"sync"

	"github.com/letiantech/worker/producer"
)

type Group interface {
	Worker
}

type group struct {
	producer.OutputProducer
	workers []Worker
}

func NewGroup(number int, creator func() Worker) Group {
	workers := make([]Worker, number)
	for i := 0; i < number; i++ {
		workers[i] = creator()
	}
	return &group{
		workers: workers,
	}
}

func (g *group) Start(p producer.OutputProducer) {
	(&sync.Once{}).Do(func() {
		g.OutputProducer = p
		for i := 0; i < len(g.workers); i++ {
			g.workers[i].Start(p)
		}
	})
}
