package worker

import (
	"log"
	"runtime"
	"time"

	// Note: the more common https://github.com/shirou/gopsutil was giving build errors.
	"github.com/mikoim/go-loadavg"
)

type LoadLimiter interface {
	LoadLimit() error
}

func NewLoadLimiter(name string, maxRatio float64, throttleDuration time.Duration, throttleMax time.Duration) LoadLimiter {
	return &loadLimiter{
		name:             name,
		maxLoad:          maxRatio * float64(runtime.NumCPU()),
		throttleDuration: throttleDuration,
		throttleMax:      throttleMax,
	}
}

type loadLimiter struct {
	name             string
	maxLoad          float64
	throttleDuration time.Duration
	throttleMax      time.Duration
}

func (l *loadLimiter) LoadLimit() error {
	load, err := loadavg.Parse()
	if err != nil {
		return err
	}

	throttleTime := l.throttleDuration
	for load.LoadAverage1 > l.maxLoad {
		log.Printf("(%s) Load %.2f above threshold %.2f, throttling for %s.", l.name, load.LoadAverage1, l.maxLoad, throttleTime)
		time.Sleep(throttleTime)

		throttleTime *= 2
		if throttleTime > l.throttleMax {
			throttleTime = l.throttleMax
		}
	}

	return nil
}
