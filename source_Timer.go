package rx

import (
	"time"
)

func Timer(delay, period time.Duration) Observable[int] {
	return Func(func(subscriber Writer[int]) (err error) {
		count := 0
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
			if !subscriber.Write(count) {
				return
			}
			count++
		case <-subscriber.Dying():
			if !timer.Stop() {
				<-timer.C
			}
			return
		}
		if period == 0 {
			return
		}
		ticker := time.NewTicker(period)
		for {
			select {
			case <-ticker.C:
				if !subscriber.Write(count) {
					return
				}
				count++
			case <-subscriber.Dying():
				ticker.Stop()
				return
			}
		}
	})
}
