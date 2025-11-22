package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var app = app.New()
	var window = app.NewWindow("Seegnal")
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(500, 500))

	var rescanButton = widget.NewButton("Rescan", func() {
	});
	rescanButton.Importance = widget.HighImportance
	

	go func() {
		fyne.Do(func() {
			rescanButton.SetText("Scanning...")
		})
    net, err := scanForExistingNetworks()

    fyne.Do(func() {
        if err != nil {
            fyne.CurrentApp().SendNotification(
                fyne.NewNotification("Error", "Failed to scan for networks!"),
            )
            return
        }

		fmt.Println(net)
        rescanButton.SetText("Rescan")
    })
}()


	var topBar = container.NewHBox(widget.NewLabel("Seegnal 0.1"), layout.NewSpacer(), rescanButton)


	var data = []string{"one", "two", "three"};
	var wifiList = widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		})


		var header = container.NewBorder(container.NewVBox(topBar, widget.NewSeparator()), nil, nil, nil)

		var mainContainer = container.NewVBox()
		var splitContainer = container.NewHSplit(wifiList, mainContainer)
		splitContainer.SetOffset(0.3)
	window.SetContent(container.NewBorder(header, nil, nil, nil, splitContainer))

	window.ShowAndRun()
}


type network struct {
	ssid string
	strength int
	security string
}

func scanForExistingNetworks() ([]network, error)  {
	var networks []network
	var c = exec.Command("nmcli", "-t", "-f", "SSID,SIGNAL,SECURITY", "device", "wifi", "list")
	var out, err = c.Output();

	if err != nil {
		return networks, err 
	}


	var s_slc = strings.Split(string(out), "\n")
	for _, s := range s_slc {
		var sp = strings.Split(s, ":")

		if len(sp) > 1 {
			var si, serr  = strconv.Atoi(sp[1])
			networks = append(networks, network{sp[0], func () int { if serr != nil {return 0} else {return si}}() , sp[2]}) 
		}
	}

	return networks, nil
}
