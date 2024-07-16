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

// Returns whether to run the app.
func confirmationPrompt() bool {
	var input = make([]byte, 2)

	for {
		fmt.Printf("run program? (y/n)\n");

		_, err := os.Stdin.Read(input)
		if err != nil {
			panic(err)
		}

		switch input[0] {
		case 'y':
			return true

		case 'n':
			return false
		}
	}
}

func draw(bgLineYs []float64,
	bgLine     Sprite,
	drawIntro  int,
	hearts     [2]Sprite,
	intro      [2]Sprite,
	ponyMdl    PonyModel,
	renderer   *sdl.Renderer,
	win        *sdl.Window) {

	renderer.SetDrawColor(49, 229, 184, 255)
	renderer.Clear()

	bgLine.Rect.X = gfxWindowWidth / 2 - bgLine.Rect.W / 2

	for i := 0; i < len(bgLineYs); i++ {
		bgLine.Rect.Y = int32(bgLineYs[i])
		bgLine.Draw()
	}

	ponyMdl.Draw()

	// demo drawing, no que
	hPos := getHeartPositions()
	for i := 0; i < len(hPos); i++ {
		hearts[0].Rect.X = hPos[i].X
		hearts[0].Rect.Y = hPos[i].Y
		hearts[0].Draw()
	}

	switch drawIntro {
	case 2:
		intro[1].Draw()
		fallthrough
	case 1:
		intro[0].Draw()
	}

	renderer.Present()
}

// Returns whether mainloop should stay active.
func handleEvents(gameActive *bool, ponyMdl *PonyModel, wags *int) bool {
	event := sdl.PollEvent()

	for ; event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			*gameActive = false
			return false

		case *sdl.MouseButtonEvent:
			if event.GetType() == sdl.MOUSEBUTTONDOWN {
				if *gameActive {
					if ponyMdl.RumpIdx == 0 {
						ponyMdl.RumpIdx = 1
						ponyMdl.TailIdx = 1
					} else {
						ponyMdl.RumpIdx = 0
						ponyMdl.TailIdx = 0
					}

					*wags++
					if *wags == wagsUntilJoy {
						go startTwiJoy(ponyMdl)
					}
				}
			}
		}
	}

	return true
}

func moveBgLines(bgLineYs *[]float64,
	delta float64,
	untilBgSpawn *float64,
	lineHeight int32) {

	*untilBgSpawn -= delta
	if *untilBgSpawn <= 0 {
		*bgLineYs = append(*bgLineYs, float64(0 - lineHeight))
		*untilBgSpawn = bgLineSpawnTime
	}

	for i := 0; i < len(*bgLineYs); i++ {
		(*bgLineYs)[i] += bgLineVelocity * delta
	}

	if len(*bgLineYs) > bgMaxLines + 1 {
		*bgLineYs = (*bgLineYs)[1:]
	}
}

func startTwiJoy(ponyMdl *PonyModel) {
	joyDelayBegin := time.Now()
	for time.Since(joyDelayBegin).Seconds() * timescale < joyThroughWagsDelay {}

	ponyMdl.EyeIdx = 2
}

// Returns whether mainloop should stay active.
func tick(bgLineYs   *[]float64,
	bgText       Sprite,
	drawIntro    int,
	gameActive   *bool,
	hearts       [2]Sprite,
	intro        [2]Sprite,
	lastTick     *time.Time,
	ponyMdl      *PonyModel,
	renderer     *sdl.Renderer,
	untilBgSpawn *float64,
	wags         *int,
	win          *sdl.Window) bool {
	var (
		delta    float64
	)
	delta = time.Since(*lastTick).Seconds()
	if delta >= (1.0 / tickrate) {
		delta *= timescale

		if *gameActive {
			moveBgLines(bgLineYs,
				delta,
				untilBgSpawn,
				bgText.Rect.H)
		}

		draw((*bgLineYs)[:],
			bgText,
			drawIntro,
			hearts,
			intro,
			*ponyMdl,
			renderer,
			win)

		if handleEvents(gameActive, ponyMdl, wags) == false {
			return false
		}

		*lastTick = time.Now()
	}

	return true
}

func main() {
	var (
		bgLineYs     []float64
		bgText       Sprite
		drawIntro    int
		err          error
		font         *ttf.Font
		gameActive   bool
		hearts       [2]Sprite
		intro        [2]Sprite
		lastTick     time.Time
		ponyMdl      PonyModel
		renderer     *sdl.Renderer
		start        time.Time
		untilBgSpawn float64
		wags         int
		win          *sdl.Window
	)

	gameActive = false

	if confirmationPrompt() == false {
		return
	}

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	win, err = sdl.CreateWindow("Twilight's Program",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(gfxWindowWidth * gfxScale),
		int32(gfxWindowHeight * gfxScale),
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	renderer, err = sdl.CreateRenderer(win, -1, 0)
	if err != nil {
		panic(err)
	}

	renderer.SetLogicalSize(gfxWindowWidth, gfxWindowHeight)

	hearts[0] = newSprite(renderer)
	hearts[0].InitFromBMP("pkg/heart_big.bmp")
	defer hearts[0].Free()

	hearts[1] = newSprite(renderer)
	hearts[1].InitFromBMP("pkg/heart_small.bmp")
	defer hearts[1].Free()

	ponyMdl = newPonyModel(renderer)
	defer ponyMdl.Free()

	poo, brain := gfxPonyX, gfxPonyY // try using directly instead :)
	ponyMdl.SetX(int32(poo))
	ponyMdl.SetY(int32(brain))

	_ = ttf.Init()
	defer ttf.Quit()

	font, err = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
		20)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	intro[0] = newSprite(renderer)
	intro[0].InitFromText("YOU ARE NOW", getIntroColor(), font)
	defer intro[0].Free()

	intro[1] = newSprite(renderer)
	intro[1].InitFromText("DOG", getIntroColor(), font)
	defer intro[1].Free()

	intro[1].Rect.X = gfxWindowWidth / 2 - intro[1].Rect.W / 2
	intro[1].Rect.Y = gfxWindowHeight / 2 - intro[1].Rect.H / 2

	intro[0].Rect.X = gfxWindowWidth / 2 - intro[0].Rect.W / 2
	intro[0].Rect.Y = intro[1].Rect.Y - intro[1].Rect.H

	bgText = newSprite(renderer)
	bgText.InitFromText("wag wag wag wag", getBgTextColor(), font)
	bgLineYs = append(bgLineYs, 0)

	start = time.Now()
	untilBgSpawn = bgLineSpawnTime
	drawIntro++

	go func() {
		for time.Since(start).Seconds() * timescale < introDogTime {}
		drawIntro++
	}()

	go func() {
		for time.Since(start).Seconds() * timescale < gameStartTime {}
		gameActive = true
	}()

	go func() {
		var eyeMovement  time.Time

		for time.Since(start).Seconds() * timescale < introLifetime {}
		eyeMovement = time.Now()
		drawIntro = 0

		for gameActive && ponyMdl.EyeIdx != 2 {
			for time.Since(eyeMovement).Seconds() * timescale <
				eyeOpenedDuration {}

			if ponyMdl.EyeIdx != 2 {
				ponyMdl.EyeIdx = 1
			}
			eyeMovement = time.Now()

			for time.Since(eyeMovement).Seconds() * timescale <
				eyeClosedDuration {}

			if ponyMdl.EyeIdx != 2 {
				ponyMdl.EyeIdx = 0
			}
			eyeMovement = time.Now()
		}
	}()

mainloop:
	for {
		stayActive := tick(&bgLineYs,
			bgText,
			drawIntro,
			&gameActive,
			hearts,
			intro,
			&lastTick,
			&ponyMdl,
			renderer,
			&untilBgSpawn,
			&wags,
			win)

		if stayActive == false {
			break mainloop
		}
	}
}
