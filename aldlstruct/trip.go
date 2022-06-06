package aldlstruct

import (
	"aldl160/winlog"
	"bufio"
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"os"
	"strconv"
	"strings"
)

var testdata = []string{}
var fpos = 0

func ReadTrip(f string) {

	file, err := os.Open(f) // For read access.
	if err != nil {
		winlog.WriteStatus(fmt.Sprintf("Error opening data file %s", err.Error()))
	}

	testdata = []string{}
	fpos = 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		testdata = append(testdata, scanner.Text())
	}

	testdata = strings.Split(testdata[0], "\r")

	file.Close()

	winlog.WriteStatus(fmt.Sprintf("Read %d entries. Use navigation buttons to look through data", len(testdata)))
}

func (a aldl) TripFWD(str []binding.String) {

	fpos++

	if len(testdata) <= 0 || fpos > len(testdata) {
		winlog.WriteStatus("No further data")
		fpos = len(testdata)
		return
	}

	a.populate(str, testdata[fpos])
}

func (a aldl) TripBWD(str []binding.String) {

	fpos--

	if len(testdata) <= 0 || fpos <= 0 {
		winlog.WriteStatus("No previous data")
		fpos = 0
		return
	}

	a.populate(str, testdata[fpos])

}

func (a aldl) populate(str []binding.String, fs string) {

	data := strings.Split(fs, ",")
	for i, val := range data {
		switch i {
		case 0:
			str[1].Set(fmt.Sprintf("Record: %d at %s ms", fpos, val))
			continue
		case 1:
			trv, _ := strconv.ParseUint(val, 16, 8)
			ns := a.TranslateByte(uint8(i), uint8(trv))
			str[0].Set(ns)
			continue
		case int(a.FrameLength):
			break
		default:
			trv, _ := strconv.ParseUint(val, 16, 8)
			//fmt.Printf("%d, %s, %d\n", i, val, trv)
			ns := a.TranslateByte(uint8(i), uint8(trv))
			str[i].Set(ns)
		}
	}
}
