// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./geninfo.go
package main

import (
	"errors"
	"fmt"
	"os"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	PathAssetsSys  string
	PathAssetsUser string
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

func (s *Sprite) InitFromAsset(assetpath string) {
	var (
		err   error
		found bool
		path  string
		path_prefixes = []string{
			"./assets",
			PathAssetsUser,
			PathAssetsSys,
		}
	)

	for i := 0; i < len(path_prefixes); i++ {
		path = path_prefixes[i] + "/" + assetpath

		f, err := os.Open(path)
		defer f.Close()

		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			fmt.Fprintf(os.Stderr,
				"Asset file could not be opened: \"%v\", \"%v\"\n",
				path, err)
		} else {
			found = true
			break
		}
	}

	if found == false {
		panic(fmt.Sprintf("Asset not found in asset paths: %v\n",
			path_prefixes))
	}

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
