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
	"time"

	"github.com/letiantech/worker/job"
	"github.com/letiantech/worker/limiter"
)

type LimitedProducer interface {
	Producer
	OutputLimiter() limiter.Limiter
	InputLimiter() limiter.Limiter
}

type LimitedOutputProducer interface {
	OutputProducer
	OutputLimiter() limiter.Limiter
}

type LimitedInputProducer interface {
	InputProducer
	InputLimiter() limiter.Limiter
}

type limitedProducer struct {
	Producer
	il limiter.Limiter
	ol limiter.Limiter
}

func NewLimitedProducer(p Producer, inputLimiter limiter.Limiter, outputLimiter limiter.Limiter) LimitedProducer {
	if p == nil {
		return nil
	}
	return &limitedProducer{
		Producer: p,
		il:       inputLimiter,
		ol:       outputLimiter,
	}
}

func (p *limitedProducer) InputLimiter() limiter.Limiter {
	return p.il
}

func (p *limitedProducer) PushJob(j job.Job) error {
	if p.il != nil {
		tm := p.il.Update()
		if tm != 0 {
			time.Sleep(tm)
		}
	}
	return p.Producer.PushJob(j)
}

func (p *limitedProducer) OutputLimiter() limiter.Limiter {
	return p.ol
}

func (p *limitedProducer) GetJob() job.Job {
	if p.ol != nil {
		tm := p.ol.Update()
		if tm != 0 {
			time.Sleep(tm)
		}
	}
	return p.Producer.GetJob()
}

type limitedOutputInputProducer struct {
	OutputProducer
	lmt limiter.Limiter
}

func NewLimitedOutputProducer(p Producer, lmt limiter.Limiter) LimitedOutputProducer {
	if p == nil {
		return nil
	}
	return &limitedOutputInputProducer{
		OutputProducer: p,
		lmt:            lmt,
	}
}

func (p *limitedOutputInputProducer) OutputLimiter() limiter.Limiter {
	return p.lmt
}

func (p *limitedOutputInputProducer) GetJob() job.Job {
	if p.lmt != nil {
		tm := p.lmt.Update()
		if tm != 0 {
			time.Sleep(tm)
		}
	}
	return p.OutputProducer.GetJob()
}

type limitedInputProducer struct {
	InputProducer
	lmt limiter.Limiter
}

func NewLimitedInputProducer(p Producer, lmt limiter.Limiter) LimitedInputProducer {
	if p == nil {
		return nil
	}
	return &limitedInputProducer{
		InputProducer: p,
		lmt:           lmt,
	}
}

func (p *limitedInputProducer) InputLimiter() limiter.Limiter {
	return p.lmt
}

func (p *limitedInputProducer) PushJob(j job.Job) error {
	if p.lmt != nil {
		tm := p.lmt.Update()
		if tm != 0 {
			time.Sleep(tm)
		}
	}
	return p.InputProducer.PushJob(j)
}
