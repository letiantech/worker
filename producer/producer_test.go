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
)

func TestProducer(t *testing.T) {
	p := NewProducer(10)
	ip := InputProducer(p)
	op := OutputProducer(p)
	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
	j := job.NewHttpJob(req)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		j1 := op.GetJob()
		j1.Do()
		j1.Finish()
		wg.Done()
	}()
	_ = ip.PushJob(j)
	wg.Wait()
}
