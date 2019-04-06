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
	"errors"
	"net/http"
)

type HttpJob interface {
	Do()
	Wait() (*http.Response, error)
	Finish()
}

type httpResponse struct {
	res *http.Response
	err error
}

// 用于给Worker异步并发执行的 Http 类型的任务
type httpJob struct {
	res chan *httpResponse
	req *http.Request
}

// 创建一个Http任务
func NewHttpJob(req *http.Request) HttpJob {
	j := &httpJob{
		req: req,
		res: make(chan *httpResponse, 1),
	}
	return j
}

// 等待http任务执行完毕
func (j *httpJob) Wait() (*http.Response, error) {
	res, ok := <-j.res
	if !ok {
		return nil, errors.New("job already finished")
	}
	return res.res, res.err
}

// 完成http任务，回收资源
func (j *httpJob) Finish() {
	close(j.res)
}

// 执行http任务
func (j *httpJob) Do() {
	ret := &httpResponse{}
	ret.res, ret.err = http.DefaultClient.Do(j.req)
	j.res <- ret
}
