// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"time"
)

const (
	gfxWindowWidth = 320
	gfxWindowHeight = 200
	tickrate = 60
	timescale = 1.0
)

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

func getIntroColor() sdl.Color {
	return sdl.Color {
		R: 255,
		G: 255,
		B: 255,
	}
}

func draw(drawIntro int,
	introR [2]sdl.Rect,
	introS [2]*sdl.Surface,
	surface *sdl.Surface,
	win *sdl.Window) {

	bgColor := sdl.MapRGB(surface.Format, 54, 254, 204)
	surface.FillRect(nil, bgColor)

	switch drawIntro {
	case 2:
		introS[1].Blit(nil, surface, &introR[1])
		fallthrough
	case 1:
		introS[0].Blit(nil, surface, &introR[0])
	}

	win.UpdateSurface()
}

// Returns whether mainloop should stay active.
func handleEvents(gameActive *bool, wags *int) bool {
	event := sdl.PollEvent()

	for ; event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			*gameActive = false
			return false

		case *sdl.MouseButtonEvent:
			if event.GetType() == sdl.MOUSEBUTTONDOWN {
				if *gameActive {
					fmt.Printf("Wagged\n")
				}

				*wags++
				if *wags == WagsUntilJoy {
					go startTwiJoy()
				}
			}
		}
	}

	return true
}

func startTwiJoy() {
	joyDelayBegin := time.Now()
	for time.Since(joyDelayBegin) < JoyThroughWagsDelay {}

	fmt.Printf("Joy expression started\n")
}

func main() {
	var (
		delta       float64
		drawIntro   int
		err         error
		eyeMovement time.Time
		font        *ttf.Font
		gameActive  bool
		input       []byte
		introR      [2]sdl.Rect
		introS      [2]*sdl.Surface
		lastTick    time.Time
		start       time.Time
		surface     *sdl.Surface
		wags        int
		win         *sdl.Window
	)

	gameActive = false
	input = make([]byte, 2)

confirmation:
	for {
		fmt.Printf("run program? (y/n)\n");

		_, err := os.Stdin.Read(input)
		if err != nil {
			panic("Cannot read stdin")
		}

		switch input[0] {
		case 'y':
			break confirmation

		case 'n':
			return
		}
	}

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	win, err = sdl.CreateWindow("Twilight's Program",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		gfxWindowWidth,
		gfxWindowHeight,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	surface, err = win.GetSurface()
	if err != nil {
		panic(err)
	}

	_ = ttf.Init()
	defer ttf.Quit()

	font, err = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 25)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	introS[0], _ = font.RenderUTF8Solid("YOU ARE NOW", getIntroColor())
	defer introS[0].Free()

	introS[1], _ = font.RenderUTF8Solid("DOG", getIntroColor())
	defer introS[1].Free()

	introR[1].W = introS[1].W
	introR[1].H = introS[1].H
	introR[1].X = gfxWindowWidth / 2 - introR[1].W / 2
	introR[1].Y = gfxWindowHeight / 2 - introR[1].H / 2

	introR[0].W = introS[0].W
	introR[0].H = introS[0].H
	introR[0].X = gfxWindowWidth / 2 - introR[0].W / 2
	introR[0].Y = introR[1].Y - introR[1].H

	start = time.Now()
	drawIntro++

	go func() {
		for time.Since(start) < IntroDogTime {}
		drawIntro++
	}()

	go func() {
		for time.Since(start) < GameStartTime {}
		fmt.Printf("start wag bg animation\n")
		gameActive = true
	}()

	go func() {
		for time.Since(start) < IntroLifetime {}
		eyeMovement = time.Now()
		drawIntro = 0

		for gameActive {
			for time.Since(eyeMovement) < EyeOpenedDuration {}
			fmt.Printf("character: Eye closed\n")
			eyeMovement = time.Now()

			for time.Since(eyeMovement) < EyeClosedDuration {}
			fmt.Printf("character: Eye opened\n")
			eyeMovement = time.Now()
		}
	}()

mainloop:
	for {
		rawDelta := time.Since(lastTick)
		if rawDelta >= (1_000_000_000 / tickrate) {
			delta = float64(rawDelta) / float64(1_000_000_000)
			delta *= timescale

			draw(drawIntro, introR, introS, surface, win)

			if handleEvents(&gameActive, &wags) == false {
				break mainloop
			}

			lastTick = time.Now()
		}
	}
}
