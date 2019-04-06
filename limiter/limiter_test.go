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
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	var testParams = []struct {
		speed float64
		count int
		exp   bool
		sleep bool
	}{
		{-1, 100, false, true},
		{20, 40, true, true},
		{0.5, 2, true, true},
		{0.00001, 1, false, false},
	}
	for _, v := range testParams {
		testFunc(t, v.speed, v.count, v.exp, v.sleep)
	}
}

func testFunc(t *testing.T, speed float64, count int, exp bool, sleep bool) {
	l := NewLimiter(speed)
	start := time.Now().UnixNano()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		tm := l.Update()
		if tm != 0 && sleep {
			time.Sleep(tm)
		}
	}
	end := time.Now().UnixNano()
	tm := (end - start) / int64(time.Second)
	if (int(tm) == int(float64(count)/l.Speed())) != exp {
		t.Fatal("failed", tm)
	}
}

func TestLimiterSleep(t *testing.T) {
	l := NewLimiter(10)
	for {
		tm := l.Update()
		if tm != 0 {
			fmt.Println(tm)
			time.Sleep(tm)
			break
		}
	}
	time.Sleep(time.Second / 8)
	l.Update()
	time.Sleep(time.Second / 2)
	l.Update()
}

func TestLimiterConcurrency(t *testing.T) {
	count := 300
	speed := 100.0
	wg := &sync.WaitGroup{}
	wg.Add(count)
	l := NewLimiter(speed)
	start := time.Now().UnixNano()
	for i := 0; i < count; i++ {
		go func() {
			tm := l.Update()
			if tm != 0 {
				time.Sleep(tm)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	end := time.Now().UnixNano()
	tm := (end - start) / int64(time.Second)
	if int(tm) != int(float64(count)/l.Speed()) {
		t.Fatal("failed", tm)
	}
}

func TestDummyLimiter(t *testing.T) {
	dl := DummyLimiter()
	dl.SetSpeed(1)
	dl.Update()
	dl.Speed()
}
