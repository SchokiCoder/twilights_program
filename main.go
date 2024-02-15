// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"image/color"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"os"
)

type twiTheme struct {}

func (twiTheme) Color(c fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameBackground:
		return &color.RGBA {54, 254, 204, 255}
	case theme.ColorNameShadow:
		return &color.RGBA {R: 0xcc, G: 0xcc, B: 0xcc, A: 0xcc}
	default:
		return color.White
	}
}

func (twiTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.LightTheme().Font(style)
}

func (twiTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (twiTheme) Size(_ fyne.ThemeSizeName) float32 {
	return 30.0
}

/*
type TwiColor struct {}

type TwiSettings struct {}

type TwiApp struct {}

func (tc TwiColor) RGBA() (r, g, b, a uint32) {
	return 54, 254, 204, 255
	//background: Color::from_rgb(0.24, 0.643, 0.565),
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
	var a = app.NewWithID("twilights_program")
	var active = true;
	var win = a.NewWindow("Twilight's Program")
	var intro = widget.NewLabel("YOU ARE NOW")
	var cont = container.NewVBox(intro)
	var input = make([]byte, 2)

	a.Settings().SetTheme(twiTheme {})
	win.SetContent(cont)
	win.SetFixedSize(true)
	win.Resize(fyne.Size {320, 200})

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
