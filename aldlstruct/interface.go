package aldlstruct

type aldler interface {
	GetFrameLength() uint8
	GetDataStruct() map[uint8]string
	GetFlags() map[uint8]string
	TranslateByte(uint8, uint8) string
}
