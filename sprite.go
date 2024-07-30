// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./geninfo.go
package main

import (
	"fmt"
	"path/filepath"
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

func (s *Sprite) InitFromAsset(appPath string, assetPath string) {
	var (
		err      error
		fullpath string
		pathPrefixes = []string{
			appPath,
			filepath.Join(appPath, "images"),
			filepath.Join(appPath, AppName + "_data", "images"),
		}
	)

	fullpath = getFilepathFromPaths(pathPrefixes, assetPath)
	if fullpath == "" {
		panic(fmt.Sprintf("Image not found in asset paths: %v\n",
			pathPrefixes))
	}

	s.surface, err = img.Load(fullpath)
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

func (s *Sprite) InitFromText(text string,
	colors []sdl.Color,
	fonts []*ttf.Font) {
	var (
		err error
		allS []*sdl.Surface
	)

	for i := len(fonts) - 1; i >= 0; i-- {
		temp, err := fonts[i].RenderUTF8Solid(text, colors[i])
		if err != nil {
			panic(err)
		}

		allS = append(allS, temp)
	}

	for i := len(allS) - 1; i >= 1; i-- {
		rect := sdl.Rect{
			X: gfxTextOutlineSize,
			Y: gfxTextOutlineSize,
			W: allS[0].W,
			H: allS[0].H,
		}
		err = allS[i].Blit(nil, allS[0], &rect)
		if err != nil {
			panic(err)
		}
	}

	s.surface = allS[0]
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
