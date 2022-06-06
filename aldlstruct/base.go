package aldlstruct

import "fmt"

type aldl struct {
	FrameLength    uint8
	dataStruct     map[uint8]string
	parseFunctions []func(uint8, uint8) string
	flags          map[uint8]string
}

var aldlMap = map[string]aldler{
	"a136": new136(),
}

func NewALDL(code string) (aldler, error) {
	return aldlMap[code], nil
}

func getBits(data uint8) []bool {
	return []bool{
		data&128 > 0,
		data&64 > 0,
		data&32 > 0,
		data&16 > 0,
		data&8 > 0,
		data&4 > 0,
		data&2 > 0,
		data&1 > 0,
	}
}

func toBit(pos, data uint8) string {
	return fmt.Sprintf("%b", data)
}

func toHex(pos, data uint8) string {
	return fmt.Sprintf("%X", data)
}
