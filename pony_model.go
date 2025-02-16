// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024 - 2025  Andy Frank Schoknecht

package main

import (
	"path/filepath"
	"github.com/veandco/go-sdl2/sdl"
)

type PonyModel struct {
	Body    Sprite
	Eye     [3]Sprite
	EyeIdx  int
	Rump    [2]Sprite
	RumpIdx int
	Tail    [2]Sprite
	TailIdx int
}

func newPonyModel(appPath string, renderer *sdl.Renderer) PonyModel {
	var ret PonyModel

	ret.Body = newSprite(renderer)
	ret.Body.InitFromAsset(appPath, filepath.Join("pony", "body.png"))

	ret.Eye[0] = newSprite(renderer)
	ret.Eye[0].InitFromAsset(appPath, filepath.Join("pony", "eye.png"))

	ret.Eye[1] = newSprite(renderer)
	ret.Eye[1].InitFromAsset(appPath, filepath.Join("pony", "eye_blink.png"))

	ret.Eye[2] = newSprite(renderer)
	ret.Eye[2].InitFromAsset(appPath, filepath.Join("pony", "eye_joy.png"))

	ret.Rump[0] = newSprite(renderer)
	ret.Rump[0].InitFromAsset(appPath, filepath.Join("pony", "rump_down.png"))

	ret.Rump[1] = newSprite(renderer)
	ret.Rump[1].InitFromAsset(appPath, filepath.Join("pony", "rump_up.png"))

	ret.Tail[0] = newSprite(renderer)
	ret.Tail[0].InitFromAsset(appPath, filepath.Join("pony", "tail_down.png"))

	ret.Tail[1] = newSprite(renderer)
	ret.Tail[1].InitFromAsset(appPath, filepath.Join("pony", "tail_up.png"))

	ret.RumpIdx = 1
	ret.TailIdx = 1

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

	for i := 0; i < len(pm.Eye); i++ {
		pm.Eye[i].Rect.X = x
	}

	for i := 0; i < len(pm.Rump); i++ {
		pm.Rump[i].Rect.X = x
	}

	for i := 0; i < len(pm.Tail); i++ {
		pm.Tail[i].Rect.X = x
	}
}

func (pm *PonyModel) SetY(y int32) {
	pm.Body.Rect.Y = y

	for i := 0; i < len(pm.Eye); i++ {
		pm.Eye[i].Rect.Y = y
	}

	for i := 0; i < len(pm.Rump); i++ {
		pm.Rump[i].Rect.Y = y
	}

	for i := 0; i < len(pm.Tail); i++ {
		pm.Tail[i].Rect.Y = y
	}
}

func (pm *PonyModel) Free() {
	pm.Body.Free()

	for i := 0; i < len(pm.Eye); i++ {
		pm.Eye[i].Free()
	}

	for i := 0; i < len(pm.Rump); i++ {
		pm.Rump[i].Free()
	}

	for i := 0; i < len(pm.Tail); i++ {
		pm.Tail[i].Free()
	}
}
