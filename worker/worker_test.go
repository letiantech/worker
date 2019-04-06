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
	"testing"
	"time"

	"github.com/letiantech/worker/job"
	"github.com/letiantech/worker/pipe"
)

func TestWorker(t *testing.T) {

	total := 100
	w := NewWorker()
	p := pipe.NewPipe(10)
	w.Start(p, p)
	jobs := make([]job.HttpJob, total)
	start := time.Now().UnixNano()
	for i := 0; i < total; i++ {
		jobs[i] = job.DummyHttpJob()
		_ = p.Send(jobs[i])
	}
	failedCount := 0
	for i := 0; i < total; i++ {
		_, err := jobs[i].Wait()
		if err != nil {
			failedCount++
		}
	}
	w.Close()
	end := time.Now().UnixNano()
	costTime := float32(end-start) / float32(time.Millisecond)
	t.Log("total: ", total, ", failed: ", failedCount)
	t.Log("cost time: ", costTime, "ms")
}
