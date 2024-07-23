// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/ttf"
)

type Sprite struct {
	renderer *sdl.Renderer
	surface  *sdl.Surface
	texture  *sdl.Texture
	Rect     sdl.Rect
}

func newSprite(renderer *sdl.Renderer) Sprite {
	var ret = Sprite {
		renderer: renderer, 
	}

	return ret
}

func (s *Sprite) InitFromFile(path string) {
	var err error

	s.surface, err = img.Load(path)
	if err != nil {
		panic(err)
	}

	s.texture, err = s.renderer.CreateTextureFromSurface(s.surface)
	if err != nil {
		panic(err)
	}

	s.Rect.W = s.surface.W
	s.Rect.H = s.surface.H
}

func (s *Sprite) InitFromText(text string, color sdl.Color, font *ttf.Font) {
	var err error

	s.surface, err = font.RenderUTF8Solid(text, color)
	if err != nil {
		panic(err)
	}

	s.texture, err = s.renderer.CreateTextureFromSurface(s.surface)
	if err != nil {
		panic(err)
	}

	s.Rect.W = s.surface.W
	s.Rect.H = s.surface.H
}

func (s *Sprite) Draw() {
	err := s.renderer.Copy(s.texture, nil, &s.Rect)
	if err != nil {
		panic(err)
	}
}

func (s *Sprite) Free() {
	s.surface.Free()
	s.texture.Destroy()
}
