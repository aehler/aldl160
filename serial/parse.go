package serial

import (
	"aldl160/winlog"
	"errors"
	"fmt"
)

type Tc struct {
	Pos  uint8
	Data uint8
}

var (
	bc             uint8   = 0
	b1             []uint8 = make([]uint8, 9)
	bitSequence    []byte  = make([]byte, 50)
	bscounter      int     = 0
	pb             byte    = 0
	frame          []uint8 = make([]uint8, 50)
	fc             uint8   = 0
	frameSync      bool    = false
	onesc          uint8   = 0
	byte1          uint8   = 0
	ByteReady      chan Tc
	fe             chan struct{}
	fstop          chan struct{}
	ReadBufferSize int
)

/*
Parsing algorythm (timing independant)
Consider 1-4 zeros in a row followed by 1 as a start zero in HIGH bit sequence
Consider long row of 1 (50+) as a startup sequence (occuring on EPROM warmup)
Consider 34-36 signals as a bit (occurs every 6,25 milliseconds at 160 baud)
Consider 23-28 zeros folloed by 8-12 ones as a 0 bit
Consider 2-3 zeros as a start LOW in a bit
Consider following 32-34 ones as a logical 1
!!Thus we get logical HIGH bit when we get 2-4 zeros followed by more than 20 and less then 40 ones
We get logical LOW bit when we get 22-24 zeroes followed by 8-12 ones
*/

func SetFrameLenth(l uint8) {
	frame = make([]uint8, l)
}

func ParseByte(s byte) {
	var bit *uint8
	if bscounter < len(bitSequence) {
		if pb != s {
			if string(pb) == "1" && string(s) == "0" { //Transiting from 1 to 0 is a breakpoint
				bscounter = 0
				var err error
				bit, err = approxBit(bitSequence[:len(bitSequence)-1])
				if err != nil {
					bc = 0
					fc = 0
					winlog.WriteStatus(err.Error())
					winlog.Log(err.Error())
				}
				b1[bc] = *bit
				bitSequence = make([]byte, 50)
				bscounter = 0
			}
			//Look for 9 ones in a row
			if !frameSync && bit != nil {
				if *bit == 1 {
					onesc++
				} else {
					onesc = 0
				}
				if onesc == 9 {
					//winlog.Log("\r\nSync frame found\r\n")
					frameSync = true
					fc = 0
					onesc = 0
				}
			}

			//Now get 9-bit sequences until found another 9 bits in a row
			if frameSync && bit != nil {
				byte1 = byte1 << 1
				byte1 = byte1 | *bit
				bc++
				if bc == 9 { //Next byte
					bc = 0
					frame[fc] = byte1
					ByteReady <- Tc{fc, byte1}
					fc++
					winlog.Log(fmt.Sprintf("%d > %d\r\n", fc, byte1))
				}
			}
			if fc == frameLength {
				//fe <- struct{}{}
				fc = 0
				frameSync = false
			}
			//if bit != nil {
			//	winlog.Log(*bit)
			//}
		}

		pb = s
		bitSequence[bscounter] = s
		bscounter++
		return
	}
	bscounter = 0 //overflow secure
}

func approxBit(bs []byte) (*uint8, error) {
	if len(bs) < 19 {
		return nil, errors.New("Datastream too short")
	}

	var res uint8 = 0

	for i, val := range bs[15:18] {
		if val != bs[i+15] {
			return nil, errors.New("Datastream corrupt")
		}
		if string(val) == "1" { //inverse ALDL logic here
			res = 0
		} else {
			res = 1
		}
	}
	fmt.Println(string(bs), " >> ", res)
	return &res, nil
}

/*
00000000
00001010
10000011
01100100
10101010
00000000
11110010
00000000
00011010
10000000
01100110
00000100
00101000
00100010
00100000
01111000
10110000
00000000
10000000
00000000
01111010
00000000
00000000
*/
