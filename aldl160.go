package main

import (
	"aldl160/aldlstruct"
	"aldl160/serial"
	"aldl160/winlog"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"sort"
	"time"
)

var portlist = []string{
	"COM1",
	"COM2",
	"COM3",
	"COM4",
	"COM5",
	"COM6",
	"COM7",
	"COM8",
	"COM9",
	"COM10",
	"COM11",
	"COM12",
	"COM13",
	"COM14",
	"COM15",
	"COM16",
}

func main() {

	//go serial.GenDataStream()

	a := app.New()
	w := a.NewWindow("ALDL160")

	errMsg := binding.NewString()
	console := widget.NewMultiLineEntry()
	winlog.Init(errMsg, console)

	combo := widget.NewSelect(portlist, func(value string) {
		serial.SelectedPort <- value
	})
	sform := &widget.Form{}
	sform.Append("Port", combo)

	btn1 := widget.NewButton("Start", func() {
		go serial.StartDump()
	})

	btn2 := widget.NewButton("Stop", func() {
		serial.StopDump()
	})

	sform.Append("Start processing", container.New(layout.NewFormLayout(), btn1, btn2))

	tripFile := widget.NewEntry()

	bfs := container.NewVBox(tripFile, widget.NewButton("Apply", func() {
		winlog.Log(fmt.Sprintf("Data from file %s", tripFile.Text))
		aldlstruct.ReadTrip(tripFile.Text)
	}))

	sform.Append("Parse recorded trip", bfs)

	aldl, err := aldlstruct.NewALDL("a136")
	if err != nil {
		winlog.WriteStatus(err.Error())
		return
	}

	serial.Init(aldl.GetFrameLength())

	str := make([]binding.String, aldl.GetFrameLength())

	keys := make([]int, 0, len(aldl.GetDataStruct()))
	for k := range aldl.GetDataStruct() {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	var gridLayout []fyne.CanvasObject
	for _, k := range keys {
		gridLayout = append(gridLayout, canvas.NewText(aldl.GetDataStruct()[uint8(k)], color.Gray{16}))
		str[k] = binding.NewString()
		s := widget.NewLabelWithData(str[k])
		gridLayout = append(gridLayout, s)
	}

	keys = make([]int, 0, len(aldl.GetFlags()))
	for k := range aldl.GetFlags() {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	var gridLayout2 []fyne.CanvasObject
	for _, k := range keys {
		//cv := canvas.NewText(aldl.GetFlags()[uint8(k)], color.Gray{16})
		//gridLayout2 = append(gridLayout2, cv)
		str[k] = binding.NewString()
		s := widget.NewLabelWithData(str[k])
		gridLayout2 = append(gridLayout2, s)
	}

	go func(c chan (serial.Tc)) {
		//winlog.Log("Waiting for data")
		for {
			select {
			case val := <-c:
				if val.Pos == 0 {
					winlog.Log("\r\n")
				}
				ns := aldl.TranslateByte(val.Pos, val.Data)
				str[val.Pos].Set(ns)
				//winlog.Log(fmt.Sprintf("%d,%d,%s ", val.Pos, val.Data, ns))
			default:
				//winlog.WriteStatus("Waiting for data")
				time.Sleep(time.Duration(time.Microsecond * 5))
			}

		}
	}(serial.ByteReady)

	grid := container.New(layout.NewGridLayout(4), gridLayout...)
	grid2 := container.New(layout.NewGridLayout(2), gridLayout2...)

	form := &widget.Form{}
	form.Append("Console", console)

	grid2.Resize(fyne.Size{500, 500})

	tabs := container.NewAppTabs(
		container.NewTabItem("Data", grid),
		container.NewTabItem("Errors", grid2),
		container.NewTabItem("Settings", container.NewVBox(sform)),
		container.NewTabItem("Console", console),
	)

	btnFwd := widget.NewButton(">>", func() {
		aldl.TripFWD(str)
	})

	btnBwd := widget.NewButton("<<", func() {
		aldl.TripBWD(str)
	})

	playBar := container.NewBorder(
		container.New(layout.NewFormLayout(), btnBwd, btnFwd),
		widget.NewLabelWithData(errMsg),
		nil,
		nil)

	statusBar := container.NewBorder(tabs, playBar, nil, nil)

	winlog.Log("Application started")

	go func() {
		for {
			select {
			case c := <-serial.SelectedPort:
				serial.ReadDataStream(c)
			default:
				time.Sleep(time.Millisecond * 500)
			}

		}

	}()

	w.SetContent(statusBar)
	w.ShowAndRun()

}
