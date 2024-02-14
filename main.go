// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
)

/*
type TwiColor struct {}

type TwiTheme struct {}

type TwiSettings struct {}

type TwiApp struct {}

func (tc TwiColor) RGBA() (r, g, b, a uint32) {
	return 54, 254, 204, 255
	//background: Color::from_rgb(0.24, 0.643, 0.565),
}

func (tt TwiTheme) Color() TwiColor {
	return TwiColor {}
}

func (ts TwiSettings) Theme() TwiTheme {
	return TwiTheme {}
}

func (ts TwiSettings) Scale() float32 {
	return 2.0
}

func New() TwiApp {
	return TwiApp {}
}

func (t TwiApp) Settings() TwiSettings {
	return TwiSettings {}
}
*/

func main() {
	var a = app.New()
	var active = true;
	var win = a.NewWindow("Twilight's Program")
	var intro = widget.NewLabel("YOU ARE NOW")
	var cont = container.NewVBox(intro)
	var input = make([]byte, 2)

	win.SetFixedSize(true)
	win.Resize(fyne.Size {320, 200})
	win.SetContent(cont)

	for active {
		fmt.Printf("run program? (y/n)\n");

		_, err := os.Stdin.Read(input)
		if err != nil {
			panic("Cannot read stdin")
		}

		switch input[0] {
		case 'y':
			win.ShowAndRun()
			active = false
		
		case 'n':
			active = false
		
		default:
		}
	}
}
