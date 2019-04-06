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
	"log"
	"sync"

	"github.com/letiantech/worker/pipe"

	"github.com/letiantech/worker/job"
)

type Worker interface {
	Start(p pipe.Receiver, closer ...pipe.Closer)
	pipe.Closer
}

type Creator func() Worker

type worker struct {
	pipe.Receiver
	pipe.Closer
	o *sync.Once
}

// 创建一个Worker，用于异步并发处理Job
func NewWorker() Worker {
	return &worker{
		Closer: pipe.DummyPipe(),
		o:      &sync.Once{},
	}
}

func (w *worker) doJob(j job.Job) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	j.Do()
	j.Finish()
}

// 启动Worker，开始接受和执行 Job
func (w *worker) Start(p pipe.Receiver, closer ...pipe.Closer) {
	w.o.Do(func() {
		w.Receiver = p
		if len(closer) > 0 {
			w.Closer = closer[0]
		}
		go func() {
			for {
				j := w.Recv()
				if j == nil {
					if w.Closer != nil {
						w.Close()
					}
					break
				}
				w.doJob(j)
			}
		}()
	})
}
