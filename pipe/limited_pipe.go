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
	"time"

	"github.com/letiantech/worker/job"
	"github.com/letiantech/worker/limiter"
)

type LimitedPipe interface {
	Closer
	LimitedReceiver
	LimitedSender
}

type LimitedReceiver interface {
	Receiver
	ReceiverLimiter() limiter.Limiter
}

type LimitedSender interface {
	Sender
	SenderLimiter() limiter.Limiter
}

type limitedPipe struct {
	Closer
	LimitedReceiver
	LimitedSender
}

func NewLimitedPipe(p Pipe, SenderLimiter limiter.Limiter, ReceiverLimiter limiter.Limiter) LimitedPipe {
	if p == nil {
		return nil
	}
	return &limitedPipe{
		Closer:          p,
		LimitedReceiver: NewLimitedReceiver(p, ReceiverLimiter),
		LimitedSender:   NewLimitedSender(p, SenderLimiter),
	}
}

type limitedReceiver struct {
	Receiver
	lmt limiter.Limiter
}

func NewLimitedReceiver(p Pipe, lmt limiter.Limiter) LimitedReceiver {
	if p == nil {
		return nil
	}
	return &limitedReceiver{
		Receiver: p,
		lmt:      lmt,
	}
}

func (p *limitedReceiver) ReceiverLimiter() limiter.Limiter {
	return p.lmt
}

func (p *limitedReceiver) Recv() job.Job {
	if p.lmt != nil {
		tm := p.lmt.Update()
		if tm != 0 {
			time.Sleep(tm)
		}
	}
	return p.Receiver.Recv()
}

type limitedSender struct {
	Sender
	lmt limiter.Limiter
}

func NewLimitedSender(p Pipe, lmt limiter.Limiter) LimitedSender {
	if p == nil {
		return nil
	}
	return &limitedSender{
		Sender: p,
		lmt:    lmt,
	}
}

func (p *limitedSender) SenderLimiter() limiter.Limiter {
	return p.lmt
}

func (p *limitedSender) Send(j job.Job) error {
	if p.lmt != nil {
		tm := p.lmt.Update()
		if tm != 0 {
			time.Sleep(tm)
		}
	}
	return p.Sender.Send(j)
}
