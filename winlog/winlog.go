package winlog

import (
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var e binding.String
var c *widget.Entry

func Init(errMsg binding.String, console *widget.Entry) {
	e = errMsg
	c = console
}

func WriteStatus(s string) {
	e.Set(s)
}

func Log(s interface{}) {
	pd := c.Text
	if len(pd) > 1000 {
		pd = pd[len(pd)-999:]
	}
	c.SetText(fmt.Sprintf("%s%v", pd, s))
}
