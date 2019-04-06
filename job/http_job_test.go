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

package job

import (
	"net/http"
	"sync"
	"testing"
)

type DummyHttpJob struct {
	ch chan struct{}
}

func (dhj *DummyHttpJob) Do() {
	(&sync.Once{}).Do(func() {
		dhj.ch = make(chan struct{})
	})
}

func (dhj *DummyHttpJob) Finish() {
	(&sync.Once{}).Do(func() {
		close(dhj.ch)
	})
}

func (dhj *DummyHttpJob) Wait() (*http.Response, error) {
	<-dhj.ch
	return nil, nil
}

func TestHttpJob(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
	job := NewHttpJob(req)
	go func() {
		job.Do()
		job.Finish()
	}()
	resp, err := job.Wait()
	if err != nil {
		t.Log(err)
	}
	t.Log("ok", resp.StatusCode)
}
