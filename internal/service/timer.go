package service

import (
	"time"
)

type Timer struct {
	duration   time.Duration
	ticker     *time.Ticker
	paused     bool
	done       chan bool
	pause      chan bool
	timePassed time.Duration
	startTime  time.Time
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		duration:   duration,
		paused:     true,
		done:       make(chan bool),
		pause:      make(chan bool),
		timePassed: 0,
	}
}

func (t *Timer) Start() {
	if t.paused {
		if t.duration == 0 {
			t.duration = 1
		}
		t.ticker = time.NewTicker(t.duration.Abs())
		t.startTime = time.Now()
		t.timePassed = 0
		t.paused = false
		go func() {
			defer t.ticker.Stop()

			for {
				select {
				case <-t.ticker.C:
					t.done <- true
					return
				}
			}
		}()
	}
}

func (t *Timer) Paused() <-chan bool {
	return t.pause
}

func (t *Timer) Pause() {
	if !t.IsPaused() {
		t.paused = true
		t.timePassed = time.Duration(time.Now().Unix()-t.startTime.Unix()) * time.Second
		t.pause <- true
		t.ticker.Stop()
	}
}

func (t *Timer) IsPaused() bool {
	return t.paused
}

func (t *Timer) Done() <-chan bool {
	return t.done
}
