// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"time"
)

const (
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

func draw(surface *sdl.Surface) {
	bgColor := sdl.MapRGB(surface.Format, 54, 254, 204)
	surface.FillRect(nil, bgColor)
}

func main() {
	var (
		delta       float64
		err         error
		eyeMovement time.Time
		gameActive  bool
		input       []byte
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
		320,
		200,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	surface, err = win.GetSurface()
	if err != nil {
		panic(err)
	}

	start = time.Now()

	go func() {
		for time.Since(start) < IntroDogTime {}
		fmt.Printf("add 'DOG' to intro\n")
	}()

	go func() {
		for time.Since(start) < GameStartTime {}
		fmt.Printf("start wag bg animation\n")
		gameActive = true
	}()

	go func() {
		for time.Since(start) < IntroLifetime {}
		eyeMovement = time.Now()
		fmt.Printf("hide intro text\n")

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
			event := sdl.PollEvent()

			draw(surface)

			for ; event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					gameActive = false
					break mainloop

				case *sdl.MouseButtonEvent:
					if gameActive {
						fmt.Printf("Wagged\n")
					}

					wags++
					if wags == WagsUntilJoy {
						go func() {
							joyDelayBegin := time.Now()
							for time.Since(joyDelayBegin) < JoyThroughWagsDelay {}
							fmt.Printf("Joy expression started\n")
						}()
					}
				}
			}

			lastTick = time.Now()
		}
	}
}
