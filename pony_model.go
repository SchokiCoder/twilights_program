// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

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
