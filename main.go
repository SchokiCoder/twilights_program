// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./geninfo.go
package main

import (
	"errors"
	"path/filepath"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/mix"
	"os"
	"time"
)

var (
	AppLicense    string
	AppLicenseUrl string
	AppName       string
	AppRepository string
	AppVersion    string
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

func draw(bgLineYs     []float64,
	bgLine         Sprite,
	drawIntro      int,
	hearts         [2]Sprite,
	heartLifetimes []float64,
	intro          [2]Sprite,
	ponyMdl        PonyModel,
	renderer       *sdl.Renderer,
	win            *sdl.Window) {

	renderer.SetDrawColor(49, 229, 184, 255)
	renderer.Clear()

	bgLine.Rect.X = gfxWindowWidth / 2 - bgLine.Rect.W / 2

	for i := 0; i < len(bgLineYs); i++ {
		bgLine.Rect.Y = int32(bgLineYs[i])
		bgLine.Draw()
	}

	ponyMdl.Draw()

	hPos := getHeartPositions()
	for i := 0; i < len(heartLifetimes); i++ {
		if heartLifetimes[i] >= heartLifetime - heartBigLifetime {
			hearts[0].Rect.X = hPos[i].X
			hearts[0].Rect.Y = hPos[i].Y
			hearts[0].Draw()
		} else if heartLifetimes[i] >= heartLifetime -
				heartBigLifetime - heartSmallLifetime {
			hearts[1].Rect.X = hPos[i].X
			hearts[1].Rect.Y = hPos[i].Y
			hearts[1].Draw()
		}
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

func getFilepathFromPaths(pathPrefixes []string, path string) string {
	var fullpath string

	for i := 0; i < len(pathPrefixes); i++ {
		fullpath = filepath.Join(pathPrefixes[i], path)

		f, err := os.Open(fullpath)
		defer f.Close()

		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			fmt.Fprintf(os.Stderr,
				"File could not be opened: \"%v\", \"%v\"\n",
				fullpath, err)
		} else {
			return fullpath
		}
	}

	return ""
}

// Returns whether mainloop should stay active.
func handleEvents(gameActive *bool,
	heartCount *int,
	heartQue *int,
	ponyMdl *PonyModel,
	wags *int) bool {
	event := sdl.PollEvent()

	for ; event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			*gameActive = false
			return false

		case *sdl.MouseButtonEvent:
			if event.GetType() == sdl.MOUSEBUTTONDOWN {
				onWag(gameActive,
					heartCount,
					heartQue,
					ponyMdl,
					wags)
			}
		}
	}

	return true
}

func initAudio(appPath string, sounds []*mix.Music) {
	var err error

	err = mix.Init(0)
	if err != nil {
		panic(err)
	}

	err = mix.OpenAudio(48000, sdl.AUDIO_S16, 2, 4096)
	if err != nil {
		panic(err)
	}

	pathPrefixes := []string{
		appPath,
		filepath.Join(appPath, "sounds"),
		filepath.Join(appPath, AppName + "_data", "sounds"),
	}
	paths := []string{
		"Mappy - Main Theme.ogg",
		"Mappy - Bonus Round Fanfare.ogg",
		"Mappy - Bonus Round.ogg",
		"Mappy - Round Clear.ogg",
	}

	for i := 0; i < len(paths); i++ {
		fullpath := getFilepathFromPaths(pathPrefixes, paths[i])
		if fullpath == "" {
			panic(fmt.Sprintf("Sound not found in asset paths: %v\n",
				pathPrefixes))
		}

		sounds[i], err = mix.LoadMUS(fullpath)
		if err != nil {
			panic(err)
		}
	}
}

func initGfx(appPath string,
	hearts []Sprite,
	ponyMdl *PonyModel,
	renderer *sdl.Renderer,
	win *sdl.Window) {

	hearts[0] = newSprite(renderer)
	hearts[0].InitFromAsset(appPath, "heart/big.png")

	hearts[1] = newSprite(renderer)
	hearts[1].InitFromAsset(appPath, "heart/small.png")

	*ponyMdl = newPonyModel(appPath, renderer)

	poo, brain := gfxPonyX, gfxPonyY // try using directly instead :)
	ponyMdl.SetX(int32(poo))
	ponyMdl.SetY(int32(brain))
}

func initText(appPath string,
	bgLineYs *[]float64,
	bgText         *Sprite,
	fonts          []*ttf.Font,
	intro          []Sprite,
	renderer       *sdl.Renderer) {
	var err error

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	pathPrefixes := []string{
		appPath,
		filepath.Join(appPath, "fonts"),
		filepath.Join(appPath, AppName + "_data", "fonts"),
	}

	fullpath := getFilepathFromPaths(pathPrefixes, "DejaVuSansMono.ttf")
	if fullpath == "" {
		panic(fmt.Sprintf("Font not found in asset paths: %v\n",
			pathPrefixes))
	}

	for i := 0; i < len(fonts); i++ {
		fonts[i], err = ttf.OpenFont(fullpath, gfxFontSize)
		if err != nil {
			panic(err)
		}
	}
	fonts[1].SetOutline(gfxTextOutlineSize)

	intro[0] = newSprite(renderer)
	intro[0].InitFromText("YOU ARE NOW", getIntroColors(), fonts[:])

	intro[1] = newSprite(renderer)
	intro[1].InitFromText("DOG", getIntroColors(), fonts[:])

	intro[1].Rect.X = gfxWindowWidth / 2 - intro[1].Rect.W / 2
	intro[1].Rect.Y = gfxWindowHeight / 2 - intro[1].Rect.H / 2

	intro[0].Rect.X = gfxWindowWidth / 2 - intro[0].Rect.W / 2
	intro[0].Rect.Y = intro[1].Rect.Y - intro[1].Rect.H

	*bgText = newSprite(renderer)
	bgText.InitFromText("wag wag wag wag",
		[]sdl.Color{getBgTextColor()},
		fonts[:1])
	*bgLineYs = append(*bgLineYs, gfxFirstBgLineYOffset)
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

	if len(*bgLineYs) > gfxBgMaxLines + 1 {
		*bgLineYs = (*bgLineYs)[1:]
	}
}

func onWag(gameActive *bool,
	heartCount    *int,
	heartQue      *int,
	ponyMdl       *PonyModel,
	wags          *int) {

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

		if int(float64(*wags) / wagsForHeart) > *heartCount {
			*heartCount++
			*heartQue++
		}
	}
}

func quitAudio(sounds []*mix.Music) {
	for i := 0; i < len(sounds); i++ {
		sounds[i].Free()
	}

	mix.CloseAudio()

	mix.Quit()
}

func quitGfx(hearts []Sprite,
	ponyMdl     *PonyModel,
	renderer    *sdl.Renderer) {

	for i := 0; i < len(hearts); i++ {
		hearts[i].Free()
	}
	ponyMdl.Free()
	renderer.Destroy()
}

func quitText(bgText *Sprite,
	fonts        []*ttf.Font,
	intro        []Sprite,) {

	for i := 0; i < len(intro); i++ {
		intro[i].Free()
	}
	bgText.Free()

	for i := 0; i < len(fonts); i++ {
		fonts[i].Close()
	}

	ttf.Quit()
}

func startTwiJoy(ponyMdl *PonyModel) {
	joyDelayBegin := time.Now()
	for time.Since(joyDelayBegin).Seconds() * timescale < joyThroughWagsDelay {}

	ponyMdl.EyeIdx = 2
}

// Returns whether mainloop should stay active.
func tick(bgLineYs     *[]float64,
	bgText         Sprite,
	drawIntro      int,
	gameActive     *bool,
	heartCount     *int,
	heartQue       *int,
	hearts         [2]Sprite,
	heartLifetimes []float64,
	intro          [2]Sprite,
	lastTick       *time.Time,
	ponyMdl        *PonyModel,
	renderer       *sdl.Renderer,
	untilBgSpawn   *float64,
	uptime         *float64,
	wags           *int,
	win            *sdl.Window) bool {
	var (
		delta float64
	)
	delta = time.Since(*lastTick).Seconds()
	if delta >= (1.0 / tickrate) {
		delta *= timescale
		*uptime += delta

		if *gameActive {
			moveBgLines(bgLineYs,
				delta,
				untilBgSpawn,
				bgText.Rect.H)

			for i := 0; i < len(heartLifetimes); i++ {
				heartLifetimes[i] -= delta

				if heartLifetimes[i] <= 0.0 && *heartQue > 0 {
					heartLifetimes[i] = heartLifetime
					*heartQue--
				}
			}
		}

		draw((*bgLineYs)[:],
			bgText,
			drawIntro,
			hearts,
			heartLifetimes,
			intro,
			*ponyMdl,
			renderer,
			win)

		if handleEvents(gameActive,
			heartCount,
			heartQue,
			ponyMdl,
			wags) == false {
			return false
		}

		*lastTick = time.Now()
	}

	return true
}

func main() {
	var (
		appPath        string
		bgLineYs       []float64
		bgText         Sprite
		drawIntro      int
		err            error
		fonts          [2]*ttf.Font
		gameActive     bool
		heartCount     int
		heartQue       int
		hearts         [2]Sprite
		heartLifetimes [8]float64
		intro          [2]Sprite
		lastTick       time.Time
		mainloopActive bool
		ponyMdl        PonyModel
		renderer       *sdl.Renderer
		sounds         [4]*mix.Music
		start          time.Time
		untilBgSpawn   float64
		uptime         float64
		wags           int
		win            *sdl.Window
	)

	appPath, err = os.Executable()
	if err != nil {
		panic(err)
	}
	appPath = filepath.Dir(appPath)

	mainloopActive = true
	gameActive = false

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-a":
			fallthrough
		case "--about":
			fmt.Printf("The sound files used are not created by me.\n"+
				"They have been composed by Nobuyuki Ohnogi, "+
				"for \"Mappy\", developed by Namco in 1983.\n"+
				"\n"+
				"The font used \"DejaVuSansMono\" is not mine.\n"+
				"For more info visit:\n"+
				"https://dejavu-fonts.github.io.\n"+
				"\n"+
				"The source code of \"%v\" %v is available, "+
				"licensed under the %v at:\n"+
				"%v\n\n"+
				"If you did not receive a copy of the license, "+
				"see below:\n"+
				"%v\n",
				AppName, AppVersion,
				AppLicense,
				AppRepository,
				AppLicenseUrl)
			return
		}
	}

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	initAudio(appPath, sounds[:])
	defer quitAudio(sounds[:])

	sounds[0].Play(-1)

	if confirmationPrompt() == false {
		return
	}

	sounds[1].Play(0)

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

	initGfx(appPath, hearts[:], &ponyMdl, renderer, win)
	defer quitGfx(hearts[:], &ponyMdl, renderer)

	initText(appPath, &bgLineYs, &bgText, fonts[:], intro[:], renderer)
	defer quitText(&bgText, fonts[:], intro[:])

	heartLifetimes[1] = heartLifetime
	heartLifetimes[2] = heartLifetime - heartBigLifetime - 0.001
	drawIntro++

	start = time.Now()
	lastTick = time.Now()
	untilBgSpawn = bgLineSpawnTime

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

	go func() {
		for gameActive == false {}
		if mainloopActive {
			sounds[2].Play(0)
		}
	}()

	// autoclicker matching original wagspeed
	/*go func() {
		var lastWag time.Time

		for gameActive == false {}

		lastWag = time.Now()
		for gameActive {
			if time.Since(lastWag).Seconds() * timescale > (1.0 / 24.0 * 3.0) {
				onWag(&gameActive,
					&heartCount,
					&heartQue,
					&ponyMdl,
					&wags)
				lastWag = time.Now()
			}
		}
	}()*/

	for mainloopActive {
		mainloopActive = tick(&bgLineYs,
			bgText,
			drawIntro,
			&gameActive,
			&heartCount,
			&heartQue,
			hearts,
			heartLifetimes[:],
			intro,
			&lastTick,
			&ponyMdl,
			renderer,
			&untilBgSpawn,
			&uptime,
			&wags,
			win)
	}

	sounds[3].Play(0)

	hadJoy := func() string {
		if wags >= wagsUntilJoy {
			return "All"
		} else {
			return "No"
		}
	}()
	fmt.Printf(`Within %.2f seconds,
Twiggy wagged %v times,
and produced %v hearts of joy.
%v ponies had joy in the making of this film.
`, uptime - gameStartTime, wags, heartCount, hadJoy)

	for mix.PlayingMusic() {}
}
