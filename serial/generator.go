package serial

import (
	"aldl160/winlog"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

var SelectedPort chan string
var datacounter int64
var frameLength uint8

func Init(fl uint8) {
	ByteReady = make(chan Tc)
	SelectedPort = make(chan string)
	frameLength = fl
	SetFrameLenth(fl)
	fe = make(chan struct{})
	fstop = make(chan struct{})
	ReadBufferSize = 1
}

func ReadDataStream(port string) {

	if port == "" {
		winlog.WriteStatus("Port not selected")
		return
	}

	c := &Config{Name: port, Baud: 115200}
	s, err := OpenPort(c)
	if err != nil {
		winlog.WriteStatus(fmt.Sprintf("Cannot open port: %s", err.Error()))
		return
	}

	winlog.WriteStatus(fmt.Sprintf("Port %s open", port))

	defer func() {
		s.Close()
	}()

	go func() {
		for {
			winlog.WriteStatus(fmt.Sprintf("Received %d", datacounter))
			time.Sleep(time.Second)
		}
	}()

	//	go func() {
	for {
		buf := make([]byte, ReadBufferSize)
		_, err := s.Read(buf)
		if err != nil {
			log.Fatal("Cannot read port ", err)
		}
		datacounter++
		if string(buf[0]) == "0" || string(buf[0]) == "1" {
			ParseByte(buf[0])
		}
	}

	//	}()

}

func GenDataStream() {

	c := &Config{Name: "COM15", Baud: 115200}
	s, err := OpenPort(c)
	if err != nil {
		log.Fatal("Serial open error: ", err)
	}

	defer func() {
		s.Flush()
		s.Close()
	}()

	for {
		file, err := os.Open("assets/LoggedData2.csv") // For read access.
		if err != nil {
			log.Fatal("Error opening data file", err)
		}

		testdata := []string{}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			testdata = append(testdata, scanner.Text())
		}

		for _, d := range testdata {
			_, err := s.Write([]byte(fmt.Sprintf("%s\r", d)))
			if err != nil {
				continue
			}
			time.Sleep(time.Microsecond * 16)
		}
		file.Close()
	}

}
