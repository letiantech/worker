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
	"testing"

	"github.com/letiantech/worker/job"
	"github.com/letiantech/worker/limiter"
)

func TestLimitedPipe(t *testing.T) {
	lmt1 := limiter.NewLimiter(1)
	lmt2 := limiter.NewLimiter(1)
	p := NewLimitedPipe(NewPipe(10), lmt1, lmt2)
	if p == nil {
		t.Fatal("is nil")
	}
	p.ReceiverLimiter()
	p.SenderLimiter()
	s := Sender(p)
	wg := &sync.WaitGroup{}
	testPipe(p, wg)
	_ = s.Send(job.DummyJob())
	wg.Wait()
	p.Close()
}

func TestLimitedPipe_nil(t *testing.T) {
	p := NewLimitedPipe(nil, nil, nil)
	if p != nil {
		t.Fatal("not nil")
	}
	s := NewLimitedSender(nil, nil)
	if s != nil {
		t.Fatal("not nil")
	}
	r := NewLimitedReceiver(nil, nil)
	if r != nil {
		t.Fatal("not nil")
	}
}

func TestLimitedReceiver(t *testing.T) {
	lmt := limiter.NewLimiter(1)
	p := NewPipe(10)
	r := NewLimitedReceiver(p, lmt)
	wg := &sync.WaitGroup{}
	testPipe(r, wg)
	_ = p.Send(job.DummyJob())
	wg.Wait()
	p.Close()
}

func TestLimitedSender(t *testing.T) {
	lmt := limiter.NewLimiter(1)
	p := NewPipe(10)
	s := NewLimitedSender(p, lmt)
	wg := &sync.WaitGroup{}
	testPipe(p, wg)
	_ = s.Send(job.DummyJob())
	wg.Wait()
	p.Close()
}
