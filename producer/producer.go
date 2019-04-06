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

package producer

import (
	"sync"
	"sync/atomic"

	"github.com/letiantech/worker/job"

	"github.com/pkg/errors"
)

const (
	MinProducerSize = 10
)

var ErrProducerClosed = errors.New("producer is closed")

type Producer interface {
	PushJob(j job.Job) error
	GetJob() job.Job
	Close()
}

type InputProducer interface {
	PushJob(j job.Job) error
	Close()
}

type OutputProducer interface {
	GetJob() job.Job
	Close()
}

type producer struct {
	jobs    chan job.Job
	running int32
	wg      *sync.WaitGroup
	o       *sync.Once
}

func NewProducer(size int) Producer {
	if size < MinProducerSize {
		size = MinProducerSize
	}
	p := &producer{
		jobs:    make(chan job.Job, size),
		running: 1,
		wg:      &sync.WaitGroup{},
		o:       &sync.Once{},
	}
	return p
}

// 将一个Job推送到Producer的队列里
func (p *producer) PushJob(j job.Job) error {
	if atomic.LoadInt32(&p.running) > 0 {
		p.wg.Add(1)
		p.jobs <- j
		p.wg.Done()
		return nil
	}
	return ErrProducerClosed
}

func (p *producer) GetJob() job.Job {
	j := <-p.jobs
	return j
}

func (p *producer) Close() {
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
