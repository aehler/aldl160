package serial

import (
	"aldl160/winlog"
	"fmt"
	"os"
	"time"
)

func StartDump() {
	file, err := os.Create(fmt.Sprintf("assets/trip-%s.csv", time.Now().Format("2006-01-02-03-04-05PM"))) // For write access.
	if err != nil {
		winlog.WriteStatus(err.Error())
		return
	}
	defer file.Close()

	for {
		select {
		case <-fe:
			datastr := fmt.Sprintf("%d.%d", time.Now().Second(), time.Now().Nanosecond())
			for _, val := range frame {
				datastr = fmt.Sprintf("%s,%X", datastr, val)
			}
			file.Write([]byte(fmt.Sprintf("%s\r", datastr)))
		case <-fstop:
			return
		}
	}
}

func StopDump() {
	fstop <- struct{}{}
}
