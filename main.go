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
	"time"
)

type twiTheme struct {}

func (twiTheme) Color(c fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameHover:
		fallthrough
	case theme.ColorNameButton:
		fallthrough
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
type twiColor struct {}

func (tc twiColor) RGBA() (r, g, b, a uint32) {
	return 54, 254, 204, 255
	//background: Color::from_rgb(0.24, 0.643, 0.565),
}
*/

/*
The durations and times are based on the video being 24 frames/second,
in which 1 frame lasts 41_666_666 nanos.
*/
const (
	IntroDogTime = 1_041_666_666
	GameStartTime = IntroDogTime + 1_791_666_666
	IntroLifetime = IntroDogTime + 1_875_000_000

	EyeOpenedDuration = 208_333_333
	EyeClosedDuration = 125_000_000

	JoyThroughWagsDelay = 41_666_666

	WagsUntilJoy = 5
)

func main() {
	var a          fyne.App
	var btnWag     *widget.Button
	var cont       *fyne.Container
	var gameActive bool
	var lblWag     *widget.Label
	var win        fyne.Window
	var lblIntro   *widget.Label
	var input      []byte

	btnWagOnTapped := func() {
		if gameActive {
			fmt.Printf("Wagged\n")
		}
	}

	a = app.NewWithID("twilights_program")
	btnWag = widget.NewButton("Wag! (should be invis soon)", btnWagOnTapped)
	gameActive = false
	lblWag = widget.NewLabel("wag wag wag wag wag wag wag wag wag wag wag")
	win = a.NewWindow("Twilight's Program")
	lblIntro = widget.NewLabel("YOU ARE NOW")
	cont = container.NewVBox(btnWag, lblIntro, lblWag)
	input = make([]byte, 2)

	a.Settings().SetTheme(twiTheme {})
	win.SetContent(cont)
	win.SetFixedSize(true)
	win.Resize(fyne.Size {Width: 320, Height: 200})

	appOnStarted := func() {
		var start, eyeMovement time.Time
		start = time.Now()

		go func() {
			for time.Since(start) < IntroDogTime {}
			lblIntro.SetText(fmt.Sprintf("%v\n%v", lblIntro.Text, "DOG"))
		}()

		go func() {
			for time.Since(start) < GameStartTime {}
			lblWag.SetText("remove my.SetText() and animate me please")
			gameActive = true
		}()

		go func() {
			for time.Since(start) < IntroLifetime {}
			eyeMovement = time.Now()
			lblIntro.SetText("")

			for gameActive {
				for time.Since(eyeMovement) < EyeOpenedDuration {}
				fmt.Printf("Eye closed, please add gfx\n")
				eyeMovement = time.Now()

				for time.Since(eyeMovement) < EyeClosedDuration {}
				fmt.Printf("Eye opened, please add gfx\n")
				eyeMovement = time.Now()
			}
		}()
	}
	a.Lifecycle().SetOnStarted(appOnStarted)

	appOnStopped := func() {
		gameActive = false
	}
	a.Lifecycle().SetOnStopped(appOnStopped)

	mainloop:
	for {
		fmt.Printf("run program? (y/n)\n");

		_, err := os.Stdin.Read(input)
		if err != nil {
			panic("Cannot read stdin")
		}

		switch input[0] {
		case 'y':
			win.ShowAndRun()
			fallthrough
		case 'n':
			break mainloop
		}
	}
}
