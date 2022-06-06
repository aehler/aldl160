package aldlstruct

import "fyne.io/fyne/v2/data/binding"

type aldler interface {
	GetFrameLength() uint8
	GetDataStruct() map[uint8]string
	GetFlags() map[uint8]string
	TranslateByte(uint8, uint8) string
	TripFWD(string2 []binding.String)
	TripBWD(string2 []binding.String)
}
