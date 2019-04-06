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

package pipe

import (
	"sync"
	"sync/atomic"

	"github.com/letiantech/worker/job"

	"errors"
)

const (
	MinSize = 10
)

var ErrPipeClosed = errors.New("pipe is closed")

type Closer interface {
	Close()
}

type Pipe interface {
	Sender
	Receiver
	Closer
}

type Sender interface {
	Send(j job.Job) error
}

type Receiver interface {
	Recv() job.Job
}

type pipe struct {
	jobs    chan job.Job
	running int32
	wg      *sync.WaitGroup
	o       *sync.Once
}

func NewPipe(size int) Pipe {
	if size < MinSize {
		size = MinSize
	}
	p := &pipe{
		jobs:    make(chan job.Job, size),
		running: 1,
		wg:      &sync.WaitGroup{},
		o:       &sync.Once{},
	}
	return p
}

// 将一个Job推送到Pipe的队列里
func (p *pipe) Send(j job.Job) error {
	if atomic.LoadInt32(&p.running) > 0 {
		p.wg.Add(1)
		p.jobs <- j
		p.wg.Done()
		return nil
	}
	return ErrPipeClosed
}

func (p *pipe) Recv() job.Job {
	j := <-p.jobs
	return j
}

func (p *pipe) Close() {
	p.o.Do(func() {
		if atomic.CompareAndSwapInt32(&p.running, 1, 0) {
			atomic.StoreInt32(&p.running, 0)
			go func() {
				p.wg.Wait()
				close(p.jobs)
			}()
		}
	})
}

func DummyPipe() Pipe {
	return &dummy{
		j: job.DummyJob(),
	}
}

type dummy struct {
	j job.Job
}

func (dp *dummy) Send(j job.Job) error {
	return nil
}
func (dp *dummy) Recv() job.Job {
	return dp.j
}
func (dp *dummy) Close() {
	dp.j = nil
}
