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
	gfxScale = 3.2
	gfxWindowWidth = 200.0
	gfxWindowHeight = 150.0
	gfxPonyXPercent = 0.305555555
	gfxPonyYPercent = 0.222222222
	tickrate = 60.0
	timescale = 1.0
)

/*
The durations and times are based on the video being 24 frames/second,
in which 1 frame lasts 41_666_666 nanos.
*/
const (
	BgLineTravelTime = 0.625
	BgMaxLines = 4.0

	IntroDogTime = 1.041666666
	GameStartTime = IntroDogTime + 1.791666666
	IntroLifetime = IntroDogTime + 1.875

	EyeOpenedDuration = 0.208333333
	EyeClosedDuration = 0.125

	JoyThroughWagsDelay = 0.041666666

	WagsUntilJoy = 5
)

const (
	// pixel / second
	BgLineVelocity = gfxWindowHeight / BgLineTravelTime

	BgLineSpawnTime = BgLineTravelTime / BgMaxLines

	gfxPonyX = gfxWindowWidth * gfxPonyXPercent
	gfxPonyY = gfxWindowHeight * gfxPonyYPercent
)

func getBgTextColor() sdl.Color {
	return sdl.Color {
		R: 25,
		G: 255,
		B: 0,
	}
}

func getIntroColor() sdl.Color {
	return sdl.Color {
		R: 255,
		G: 255,
		B: 255,
	}
}

type PonyModel struct {
	Body    Sprite
	Eye     [3]Sprite
	EyeIdx  int
	Rump    [2]Sprite
	RumpIdx int
	Tail    [2]Sprite
	TailIdx int
}

func newPonyModel(renderer *sdl.Renderer) PonyModel {
	var ret PonyModel

	ret.Body = newSprite(renderer)
	ret.Body.InitFromBMP("pkg/pony_body.bmp")

	ret.Eye[0] = newSprite(renderer)
	ret.Eye[0].InitFromBMP("pkg/pony_eye.bmp")

	ret.Eye[1] = newSprite(renderer)
	ret.Eye[1].InitFromBMP("pkg/pony_eye_blink.bmp")

	ret.Eye[2] = newSprite(renderer)
	ret.Eye[2].InitFromBMP("pkg/pony_eye_joy.bmp")

	ret.Rump[0] = newSprite(renderer)
	ret.Rump[0].InitFromBMP("pkg/pony_rump_down.bmp")

	ret.Rump[1] = newSprite(renderer)
	ret.Rump[1].InitFromBMP("pkg/pony_rump_up.bmp")

	ret.Tail[0] = newSprite(renderer)
	ret.Tail[0].InitFromBMP("pkg/pony_tail_down.bmp")

	ret.Tail[1] = newSprite(renderer)
	ret.Tail[1].InitFromBMP("pkg/pony_tail_up.bmp")

	return ret
}

func (pm PonyModel) Draw() {
	pm.Body.Draw()
	pm.Eye[pm.EyeIdx].Draw()
	pm.Rump[pm.RumpIdx].Draw()
	pm.Tail[pm.TailIdx].Draw()
}

func (pm *PonyModel) SetX(x int32) {
	pm.Body.Rect.X = x
	pm.Eye[0].Rect.X = x
	pm.Eye[1].Rect.X = x
	pm.Eye[2].Rect.X = x
	pm.Rump[0].Rect.X = x
	pm.Rump[1].Rect.X = x
	pm.Tail[0].Rect.X = x
	pm.Tail[1].Rect.X = x
}

func (pm *PonyModel) SetY(y int32) {
	pm.Body.Rect.Y = y
	pm.Eye[0].Rect.Y = y
	pm.Eye[1].Rect.Y = y
	pm.Eye[2].Rect.Y = y
	pm.Rump[0].Rect.Y = y
	pm.Rump[1].Rect.Y = y
	pm.Tail[0].Rect.Y = y
	pm.Tail[1].Rect.Y = y
}

func (pm *PonyModel) Free() {
	pm.Body.Free()
	pm.Eye[0].Free()
	pm.Eye[1].Free()
	pm.Eye[2].Free()
	pm.Rump[0].Free()
	pm.Rump[1].Free()
	pm.Tail[0].Free()
	pm.Tail[1].Free()
}

func draw(bgLineYs []float64,
	bgLine Sprite,
	drawIntro int,
	intro [2]Sprite,
	ponyMdl PonyModel,
	renderer *sdl.Renderer,
	win *sdl.Window) {

	renderer.SetDrawColor(49, 229, 184, 255)
	renderer.Clear()

	bgLine.Rect.X = gfxWindowWidth / 2 - bgLine.Rect.W / 2

	for i := 0; i < len(bgLineYs); i++ {
		bgLine.Rect.Y = int32(bgLineYs[i])
		bgLine.Draw()
	}

	ponyMdl.Draw()

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
					if *wags == WagsUntilJoy {
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
		*untilBgSpawn = BgLineSpawnTime
	}

	for i := 0; i < len(*bgLineYs); i++ {
		(*bgLineYs)[i] += BgLineVelocity * delta
	}

	if len(*bgLineYs) > BgMaxLines + 1 {
		*bgLineYs = (*bgLineYs)[1:]
	}
}

func startTwiJoy(ponyMdl *PonyModel) {
	joyDelayBegin := time.Now()
	for time.Since(joyDelayBegin).Seconds() * timescale < JoyThroughWagsDelay {}

	ponyMdl.EyeIdx = 2
}

func main() {
	var (
		bgLineYs     []float64
		bgText       Sprite
		delta        float64
		drawIntro    int
		err          error
		eyeMovement  time.Time
		font         *ttf.Font
		gameActive   bool
		input        []byte
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
	untilBgSpawn = BgLineSpawnTime
	drawIntro++

	go func() {
		for time.Since(start).Seconds() * timescale < IntroDogTime {}
		drawIntro++
	}()

	go func() {
		for time.Since(start).Seconds() * timescale < GameStartTime {}
		gameActive = true
	}()

	go func() {
		for time.Since(start).Seconds() * timescale < IntroLifetime {}
		eyeMovement = time.Now()
		drawIntro = 0

		for gameActive && ponyMdl.EyeIdx != 2 {
			for time.Since(eyeMovement).Seconds() * timescale <
				EyeOpenedDuration {}

			if ponyMdl.EyeIdx != 2 {
				ponyMdl.EyeIdx = 1
			}
			eyeMovement = time.Now()

			for time.Since(eyeMovement).Seconds() * timescale <
				EyeClosedDuration {}

			if ponyMdl.EyeIdx != 2 {
				ponyMdl.EyeIdx = 0
			}
			eyeMovement = time.Now()
		}
	}()

mainloop:
	for {
		delta = time.Since(lastTick).Seconds()
		if delta >= (1.0 / tickrate) {
			delta *= timescale

			if gameActive {
				moveBgLines(&bgLineYs,
					delta,
					&untilBgSpawn,
					bgText.Rect.H)
			}

			draw(bgLineYs[:],
				bgText,
				drawIntro,
				intro,
				ponyMdl,
				renderer,
				win)

			if handleEvents(&gameActive, &ponyMdl, &wags) == false {
				break mainloop
			}

			lastTick = time.Now()
		}
	}
}
