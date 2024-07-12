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
	gfxScale = 2.0
	gfxStdWindowWidth = 200
	gfxStdWindowHeight = 150
	tickrate = 60
	timescale = 1.0
)

/*
The durations and times are based on the video being 24 frames/second,
in which 1 frame lasts 41_666_666 nanos.
*/
const (
	BgLineTravelTime = 0.625
	BgMaxLines = 4

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
	BgLineVelocity = float64(gfxStdWindowHeight) * gfxScale / BgLineTravelTime

	BgLineSpawnTime = BgLineTravelTime / float64(BgMaxLines)
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

func (pm PonyModel) Draw(scale float64, surface *sdl.Surface) {
	var (
		err  error
		rect sdl.Rect
	)

	rect = sdl.Rect {X: pm.X, Y: pm.Y,
		W: int32(float64(pm.Body.W) * scale),
		H: int32(float64(pm.Body.H) * scale)}
	err = pm.Body.BlitScaled(nil, surface, &rect)
	if err != nil {
		panic(err)
	}

	rect = sdl.Rect {X: pm.X, Y: pm.Y,
		W: int32(float64(pm.Eye[pm.EyeIdx].W) * scale),
		H: int32(float64(pm.Eye[pm.EyeIdx].H) * scale)}
	err = pm.Eye[pm.EyeIdx].BlitScaled(nil, surface, &rect)
	if err != nil {
		panic(err)
	}

	rect = sdl.Rect {X: pm.X, Y: pm.Y,
		W: int32(float64(pm.Rump[pm.RumpIdx].W) * scale),
		H: int32(float64(pm.Rump[pm.RumpIdx].H) * scale)}
	err = pm.Rump[pm.RumpIdx].BlitScaled(nil, surface, &rect)
	if err != nil {
		panic(err)
	}

	rect = sdl.Rect {X: pm.X, Y: pm.Y,
		W: int32(float64(pm.Tail[pm.TailIdx].W) * scale),
		H: int32(float64(pm.Tail[pm.TailIdx].H) * scale)}
	err = pm.Tail[pm.TailIdx].BlitScaled(nil, surface, &rect)
	if err != nil {
		panic(err)
	}
}

func draw(bgLineYs []float64,
	bgLine *sdl.Surface,
	drawIntro int,
	introR [2]sdl.Rect,
	introS [2]*sdl.Surface,
	ponyMdl PonyModel,
	scale float64,
	surface *sdl.Surface,
	win *sdl.Window,
	winW int32) {
	var err error

	bgColor := sdl.MapRGB(surface.Format, 49, 229, 184)
	surface.FillRect(nil, bgColor)

	var rect = sdl.Rect {
		X: winW / 2 - bgLine.W / 2,
		Y: 0,
		W: int32(float64(bgLine.W) * scale),
		H: int32(float64(bgLine.H) * scale),
	}
	for i := 0; i < len(bgLineYs); i++ {
		rect.Y = int32(bgLineYs[i])
		err = bgLine.Blit(nil, surface, &rect)
		if err != nil {
			panic(err)
		}
	}

	ponyMdl.Draw(gfxScale, surface)

	switch drawIntro {
	case 2:
		err = introS[1].Blit(nil, surface, &introR[1])
		if err != nil {
			panic(err)
		}
		fallthrough
	case 1:
		err = introS[0].Blit(nil, surface, &introR[0])
		if err != nil {
			panic(err)
		}
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
		bgText       *sdl.Surface
		delta        float64
		drawIntro    int
		err          error
		eyeMovement  time.Time
		font         *ttf.Font
		gameActive   bool
		input        []byte
		introR       [2]sdl.Rect
		introS       [2]*sdl.Surface
		lastTick     time.Time
		ponyMdl      PonyModel
		start        time.Time
		surface      *sdl.Surface
		untilBgSpawn float64
		wags         int
		win          *sdl.Window
		winW         int32
		winH         int32
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

	winW = int32(float64(gfxStdWindowWidth) * gfxScale)
	winH = int32(float64(gfxStdWindowHeight) * gfxScale)
	win, err = sdl.CreateWindow("Twilight's Program",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winW,
		winH,
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

	font, err = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
		20)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	introS[0], _ = font.RenderUTF8Solid("YOU ARE NOW", getIntroColor())
	defer introS[0].Free()

	introS[1], _ = font.RenderUTF8Solid("DOG", getIntroColor())
	defer introS[1].Free()

	introR[1].W = int32(float64(introS[1].W) * gfxScale)
	introR[1].H = int32(float64(introS[1].H) * gfxScale)
	introR[1].X = winW / 2 - introR[1].W / 2
	introR[1].Y = winH / 2 - introR[1].H / 2

	introR[0].W = int32(float64(introS[0].W) * gfxScale)
	introR[0].H = int32(float64(introS[0].H) * gfxScale)
	introR[0].X = winW / 2 - introR[0].W / 2
	introR[0].Y = introR[1].Y - introR[1].H

	bgText, _ = font.RenderUTF8Solid("wag wag wag wag", getBgTextColor())
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
		rawDelta := time.Since(lastTick)
		if rawDelta >= (1_000_000_000 / tickrate) {
			delta = float64(rawDelta) / float64(1_000_000_000)
			delta *= timescale

			if gameActive {
				moveBgLines(&bgLineYs,
					delta,
					&untilBgSpawn,
					bgText.H)
			}

			draw(bgLineYs[:],
				bgText,
				drawIntro,
				introR,
				introS,
				ponyMdl,
				gfxScale,
				surface,
				win,
				winW)

			if handleEvents(&gameActive, &ponyMdl, &wags) == false {
				break mainloop
			}

			lastTick = time.Now()
		}
	}
}
