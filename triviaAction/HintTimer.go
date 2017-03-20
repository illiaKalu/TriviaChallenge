package triviaAction

import (
	"time"
)

var restartChan chan bool

func init() {
	restartChan = make (chan bool)
}

func StartTimer(hub *Hub) {
    ticker := time.NewTicker(10 * time.Second).C

	for {
		select {

			case  <- ticker:
				{
					hub.Hint <- true
				}
			case <- restartChan:
				{
					ticker = time.NewTicker(10 * time.Second).C
					break
				}
		}
	}
}
func RestartTimer() {
	restartChan <- true
}
