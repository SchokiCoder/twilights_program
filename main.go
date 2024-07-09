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
	gfxWindowHeight = 240
	tickrate = 60
	timescale = 1.0
)

/*
The durations and times are based on the video being 24 frames/second,
in which 1 frame lasts 41_666_666 nanos.
*/
const (
	BgLineTravelTime = 625_000_000
	BgMaxLines = 4

	IntroDogTime = 1_041_666_666
	GameStartTime = IntroDogTime + 1_791_666_666
	IntroLifetime = IntroDogTime + 1_875_000_000

	EyeOpenedDuration = 208_333_333
	EyeClosedDuration = 125_000_000

	JoyThroughWagsDelay = 41_666_666

	WagsUntilJoy = 5
)

const (
	// pixel / second
	BgLineVelocity = gfxWindowHeight * 1_000_000_000 / BgLineTravelTime

	BgLineSpawnTime = BgLineTravelTime / BgMaxLines
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
	Body    *sdl.Surface
	Eye     [3]*sdl.Surface
	EyeIdx  int
	Rump    [2]*sdl.Surface
	RumpIdx int
	Tail    [2]*sdl.Surface
	TailIdx int
	X       int32
	Y       int32
}

func newPonyModel() PonyModel {
	var err error
	var ret PonyModel

	ret.Body, err = sdl.LoadBMP("pkg/pony_body.bmp")
	if err != nil {
		panic(err)
	}

	ret.Eye[0], err = sdl.LoadBMP("pkg/pony_eye.bmp")
	if err != nil {
		panic(err)
	}

	ret.Eye[1], err = sdl.LoadBMP("pkg/pony_eye_blink.bmp")
	if err != nil {
		panic(err)
	}

	ret.Eye[2], err = sdl.LoadBMP("pkg/pony_eye_joy.bmp")
	if err != nil {
		panic(err)
	}

	ret.Rump[0], err = sdl.LoadBMP("pkg/pony_rump_down.bmp")
	if err != nil {
		panic(err)
	}

	ret.Rump[1], err = sdl.LoadBMP("pkg/pony_rump_up.bmp")
	if err != nil {
		panic(err)
	}

	ret.Tail[0], err = sdl.LoadBMP("pkg/pony_tail_down.bmp")
	if err != nil {
		panic(err)
	}

	ret.Tail[1], err = sdl.LoadBMP("pkg/pony_tail_up.bmp")
	if err != nil {
		panic(err)
	}

	return ret
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

func (pm PonyModel) Draw(surface *sdl.Surface) {
	var rect sdl.Rect

	rect = sdl.Rect {X: pm.X, Y: pm.Y, W: pm.Body.W, H: pm.Body.H}
	pm.Body.Blit(nil, surface, &rect)

	rect = sdl.Rect {X: pm.X, Y: pm.Y, W: pm.Eye[pm.EyeIdx].W, H: pm.Eye[pm.EyeIdx].H}
	pm.Eye[pm.EyeIdx].Blit(nil, surface, &rect)

	rect = sdl.Rect {X: pm.X, Y: pm.Y, W: pm.Rump[pm.RumpIdx].W, H: pm.Rump[pm.RumpIdx].H}
	pm.Rump[pm.RumpIdx].Blit(nil, surface, &rect)

	rect = sdl.Rect {X: pm.X, Y: pm.Y, W: pm.Tail[pm.TailIdx].W, H: pm.Tail[pm.TailIdx].H}
	pm.Tail[pm.TailIdx].Blit(nil, surface, &rect)
}

func draw(bgLineYs []float64,
	bgLine *sdl.Surface,
	drawIntro int,
	introR [2]sdl.Rect,
	introS [2]*sdl.Surface,
	ponyMdl PonyModel,
	surface *sdl.Surface,
	win *sdl.Window) {

	bgColor := sdl.MapRGB(surface.Format, 49, 229, 184)
	surface.FillRect(nil, bgColor)

	var rect = sdl.Rect {
		X: gfxWindowWidth / 2 - bgLine.W / 2,
		Y: 0,
		W: bgLine.W,
		H: bgLine.H,
	}
	for i := 0; i < len(bgLineYs); i++ {
		rect.Y = int32(bgLineYs[i])
		bgLine.Blit(nil, surface, &rect)
	}

	ponyMdl.Draw(surface)

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
				}

				*wags++
				if *wags == WagsUntilJoy {
					go startTwiJoy(ponyMdl)
				}
			}
		}
	}

	return true
}

func moveBgLines(bgLineYs *[]float64,
	delta float64,
	lastBgSpawn *time.Time,
	lineHeight int32) {

	if time.Since(*lastBgSpawn) >= BgLineSpawnTime {
		*bgLineYs = append(*bgLineYs, float64(0 - lineHeight))
		*lastBgSpawn = time.Now()
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
	for time.Since(joyDelayBegin) < JoyThroughWagsDelay {}

	ponyMdl.EyeIdx = 2
}

func main() {
	var (
		bgLineYs    []float64
		bgText      *sdl.Surface
		delta       float64
		drawIntro   int
		err         error
		eyeMovement time.Time
		font        *ttf.Font
		gameActive  bool
		input       []byte
		introR      [2]sdl.Rect
		introS      [2]*sdl.Surface
		lastBgSpawn time.Time
		lastTick    time.Time
		ponyMdl     PonyModel
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

	ponyMdl = newPonyModel()
	defer ponyMdl.Free()

	_ = ttf.Init()
	defer ttf.Quit()

	font, err = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf", 32)
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

	bgText, _ = font.RenderUTF8Solid("wag wag wag wag", getBgTextColor())
	bgLineYs = append(bgLineYs, 0)

	start = time.Now()
	lastBgSpawn = start
	drawIntro++

	go func() {
		for time.Since(start) < IntroDogTime {}
		drawIntro++
	}()

	go func() {
		for time.Since(start) < GameStartTime {}
		gameActive = true
	}()

	go func() {
		for time.Since(start) < IntroLifetime {}
		eyeMovement = time.Now()
		drawIntro = 0

		for gameActive && ponyMdl.EyeIdx != 2 {
			for time.Since(eyeMovement) < EyeOpenedDuration {}
			if ponyMdl.EyeIdx != 2 {
				ponyMdl.EyeIdx = 1
			}
			eyeMovement = time.Now()

			for time.Since(eyeMovement) < EyeClosedDuration {}
			if ponyMdl.EyeIdx != 2 {
				ponyMdl.EyeIdx = 0
			}
			eyeMovement = time.Now()
		}
	}()

mainloop:
	for {
		rawDelta := time.Since(lastTick)
		if rawDelta >= (1_000_000_000 / tickrate) {
			delta = float64(rawDelta) / float64(1_000_000_000)
			delta *= timescale

			if gameActive {
				moveBgLines(&bgLineYs,
					delta,
					&lastBgSpawn,
					bgText.H)
			}

			draw(bgLineYs[:],
				bgText,
				drawIntro,
				introR,
				introS,
				ponyMdl,
				surface,
				win)

			if handleEvents(&gameActive, &ponyMdl, &wags) == false {
				break mainloop
			}

			lastTick = time.Now()
		}
	}
}
