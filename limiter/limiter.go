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

package limiter

import (
	"sync/atomic"
	"time"
)

const (
	//the slowest limited speed
	MinLimitedSpeed = 0.001
)

type Limiter interface {
	Update() time.Duration
	SetSpeed(speed float64)
	Speed() float64
}

//Limiter is used to limit speed of pipeline
//by providing a time duration
type limiter struct {
	duration int64
	lastTime int64
	nextTime int64
	speed    int64
}

//create a new limiter and set its speed
func NewLimiter(speed float64) Limiter {
	l := &limiter{}
	l.SetSpeed(speed)
	return l
}

//update the time duration of limiter
func (l *limiter) Update() time.Duration {
	duration := atomic.LoadInt64(&l.duration)
	if duration <= 0 {
		return 0
	}
	now := time.Now().UnixNano()
	atomic.CompareAndSwapInt64(&l.nextTime, 0, now)
	nextTime := atomic.AddInt64(&l.nextTime, duration)
	waitTime := nextTime - now
	if waitTime < -duration {
		waitTime = 0
	} else if waitTime < 0 {
		waitTime = duration + waitTime
	}
	return time.Duration(waitTime)
}

//set speed of limiter
func (l *limiter) SetSpeed(speed float64) {
	if speed <= 0 {
		atomic.StoreInt64(&l.duration, 0)
		return
	}
	if speed < MinLimitedSpeed {
		speed = MinLimitedSpeed
	}
	atomic.StoreInt64(&l.speed, int64(speed*1000))
	if speed < 1 {
		atomic.StoreInt64(&l.duration, int64(time.Second*1000)/int64(1000.0*speed))
	} else {
		atomic.StoreInt64(&l.duration, int64(time.Second)/int64(speed))
	}
}

func (l *limiter) Speed() float64 {
	return float64(atomic.LoadInt64(&l.speed)) / 1000.0
}

type DummyLimiter struct{}

func (dl *DummyLimiter) Update() time.Duration {
	return 0
}
func (dl *DummyLimiter) SetSpeed(speed float64) {

}
func (dl *DummyLimiter) Speed() float64 {
	return 0
}
