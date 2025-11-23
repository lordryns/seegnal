package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)


var networks = []network{};


func main() {
	var app = app.New()
	var window = app.NewWindow("Seegnal")
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(500, 500))

	var rescanButton = widget.NewButton("Rescan", func() {
	});
	rescanButton.Importance = widget.HighImportance


	var topBar = container.NewHBox(widget.NewLabel("Seegnal 0.1"), layout.NewSpacer(), rescanButton)
	var mainContainer = container.NewStack(container.NewCenter(widget.NewLabel("Welcome to Seegnal!")))

	var wifiList = widget.NewList(
		func() int {
			return len(networks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(networks[i].ssid)
		})

		var selected int
		wifiList.OnSelected = func(id widget.ListItemID) {
			selected = id
			loadNetworkDetailsInContainer(mainContainer, networks[id])
		}

		
		rescanButton.OnTapped = func() {
			go scanAndReflectChangesOnMain(rescanButton, wifiList, selected, mainContainer)
	}


	go scanAndReflectChangesOnMain(rescanButton, wifiList, -1, mainContainer)



		var header = container.NewBorder(container.NewVBox(topBar, widget.NewSeparator()), nil, nil, nil)

		var splitContainer = container.NewHSplit(wifiList, mainContainer)
		splitContainer.SetOffset(0.3)
	window.SetContent(container.NewBorder(header, nil, nil, nil, splitContainer))

	var cmd = cli.Command{Name: "seegnal", Usage: "Lightweight wifi app", Action: func(ctx context.Context, c *cli.Command) error {
		window.ShowAndRun()
		return nil
	}}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Error(fmt.Errorf("Could not start the cli..."))
	}
}


func loadNetworkDetailsInContainer(c *fyne.Container, net network) {
	var tab = container.NewVBox(widget.NewLabel(fmt.Sprintf("SSID: %v", net.ssid)))
	c.Objects = []fyne.CanvasObject{tab}
}



func scanAndReflectChangesOnMain(rescanButton *widget.Button, wifiList *widget.List, selected int, mainContainer *fyne.Container) {
		fyne.Do(func() {
			rescanButton.SetText("Scanning...")
			log.Info("Scanning...")
		})
    net, err := scanForExistingNetworks()

    fyne.Do(func() {
        if err != nil {
			log.Error("Failed to scan for networks, is this an nmcli issue?")
            fyne.CurrentApp().SendNotification(
                fyne.NewNotification("Error", "Failed to scan for networks!"),
            )
            return
        }

	
		if len(net) > 0 {
			fyne.CurrentApp().SendNotification(
                fyne.NewNotification("Found", fmt.Sprintf("%v networks found!", len(net))),
            )
		} else {
			mainContainer.Objects = []fyne.CanvasObject{container.NewCenter(widget.NewLabel("No wifi networks detected!"))}
		}


		var foundSelected int = -1
		if selected > -1 {
			foundSelected = func () int {
				// this is madness (someone stop me)
				for _, n := range networks {
					for j, k := range net {	
						if n == k {
							return j
						}
					}
				}
				return -1
			}()
		}


		networks = net
		if foundSelected > -1 {
			wifiList.Select(foundSelected)
		} else {
			wifiList.UnselectAll()
		}
		wifiList.Refresh()
		log.Infof("Scan completed! Found: %v networks\n", len(networks))

        rescanButton.SetText("Rescan")
    })
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
