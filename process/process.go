package process

import (
	"time"

	"multicard/models"
)

const INTERVAL_PERIOD time.Duration = 24 * time.Hour

const HOUR_TO_TICK int = 0
const MINUTE_TO_TICK int = 0
const SECOND_TO_TICK int = 0

func Run() chan bool {
	stop := make(chan bool)

	ticker := updateTicker()
	go func() {
		for {
			select {
			case <- ticker.C:
				// todayYear, todayMonth, todayDay := time.Now().Date()
				// cards := models.Card{}.FechtAll()
				// for _, card := range cards {
				// 	dueDateYear, dueDateMonth, dueDateDay := card.DueDate.Date()
				// 	if dueDateYear == todayYear && dueDateMonth == todayMonth && dueDateDay == todayDay {
				// 		card.ResetCredit()
				// 	}
				// }

				nowYear, nowMonth, nowDay := time.Now().Date()
				cards := models.Card{}.FechtAll()
				for _, card := range cards {
					expirationYear, expirationMonth, _ := card.ExpirationDate.Date()
					if nowYear == expirationYear && nowMonth == expirationMonth {
						card.Delete()
						// TODO: check and update, if necessary, wallet user limit after deleting the card
						continue
					}

					dueDay := card.DueDay.Day()
					if dueDay == nowDay {
						card.ResetCredit()
					}
				}
				
				ticker = updateTicker()
			case <- stop:
				ticker.Stop()
				return
			}
		}
	}()

	return stop
}

func updateTicker() *time.Ticker {
    nextTick := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), HOUR_TO_TICK, MINUTE_TO_TICK, SECOND_TO_TICK, 0, time.Local)
    if !nextTick.After(time.Now()) {
        nextTick = nextTick.Add(INTERVAL_PERIOD)
    }

	diff := nextTick.Sub(time.Now())
	if diff <= 0 {
		diff = 500 * time.Millisecond
	}

	return time.NewTicker(diff)
}