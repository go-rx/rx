package rx

import "time"

func Ticker(period time.Duration) Observable[int] {
	return Func(func(subscriber Writer[int]) (err error) {
		if period == 0 {
			return
		}
		count := 0
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
