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
	"net/http"
	"sync"
	"testing"

	"github.com/letiantech/worker/job"
	"github.com/letiantech/worker/limiter"
)

func TestLimitedProducer(t *testing.T) {
	lmt := limiter.NewLimiter(1)
	p := NewLimitedProducer(NewProducer(10), lmt, nil)
	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
	j := job.NewHttpJob(req)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		j1 := p.GetJob()
		j1.Do()
		j1.Finish()
		wg.Done()
	}()
	_ = p.PushJob(j)
	wg.Wait()
}

func TestLimitedOutputProducer(t *testing.T) {
	lmt := limiter.NewLimiter(1)
	p := NewProducer(10)
	lop := NewLimitedOutputProducer(p, lmt)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		j1 := lop.GetJob()
		j1.Do()
		j1.Finish()
		wg.Done()
	}()
	_ = p.PushJob(&job.DummyJob{})
	wg.Wait()
	p.Close()
}

func TestLimitedInputProducer(t *testing.T) {
	lmt := limiter.NewLimiter(1)
	p := NewProducer(10)
	lip := NewLimitedInputProducer(p, lmt)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		j1 := p.GetJob()
		j1.Do()
		j1.Finish()
		wg.Done()
	}()
	_ = lip.PushJob(&job.DummyJob{})
	wg.Wait()
}
