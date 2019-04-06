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
)

func testPipe(r Receiver, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		j1 := r.Recv()
		if j1 == nil {
			return
		}
		j1.Do()
		j1.Finish()
	}()
}

func TestPipe_Size1(t *testing.T) {
	p := NewPipe(1)
	s := Sender(p)
	c := Closer(p)
	wg := &sync.WaitGroup{}
	testPipe(p, wg)
	_ = s.Send(job.DummyJob())
	c.Close()
	wg.Wait()
}

func TestPipe_Size10(t *testing.T) {
	p := NewPipe(10)
	s := Sender(p)
	c := Closer(p)
	wg := &sync.WaitGroup{}
	testPipe(p, wg)
	_ = s.Send(job.DummyJob())
	wg.Wait()
	c.Close()
}

func TestPipe_Close(t *testing.T) {
	p := NewPipe(10)
	s := Sender(p)
	c := Closer(p)
	wg := &sync.WaitGroup{}
	testPipe(p, wg)
	c.Close()
	err := s.Send(job.DummyJob())
	if err == ErrPipeClosed {
	}
	wg.Wait()
}

func TestDummyPipe(t *testing.T) {
	dp := DummyPipe()
	j := job.DummyJob()
	_ = dp.Send(j)
	dp.Recv()
	dp.Close()
}
