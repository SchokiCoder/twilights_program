// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import "github.com/veandco/go-sdl2/sdl"

const (
	gfxScale = 3.2
	gfxWindowWidth = 200.0
	gfxWindowHeight = 150.0
	gfxPonyXPercent = 0.305555555
	gfxPonyYPercent = 0.222222222

	gfxPonyX = gfxWindowWidth * gfxPonyXPercent
	gfxPonyY = gfxWindowHeight * gfxPonyYPercent

	gfxTextOutlineSize = 1
)

const (
	tickrate = 60.0
	timescale = 1.0
)

/*
The durations and times are based on the video being 24 frames/second,
in which 1 frame lasts 41_666_666 nanos.
*/
const (
	bgLineTravelTime = 0.625
	bgMaxLines = 4.0

	introDogTime = 1.041666666
	gameStartTime = introDogTime + 1.791666666
	introLifetime = introDogTime + 1.875

	eyeOpenedDuration = 0.208333333
	eyeClosedDuration = 0.125

	joyThroughWagsDelay = 0.041666666

	heartBigLifetime = 0.416666666
	heartSmallLifetime = 0.666666666
	heartLifetime = 0.833333333

	wagsForHeart = 2
	wagsUntilJoy = 5
)

const (
	// pixel / second
	bgLineVelocity = gfxWindowHeight / bgLineTravelTime

	bgLineSpawnTime = bgLineTravelTime / bgMaxLines
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

func getIntroOutlineColor() sdl.Color {
	return sdl.Color {
		R: 255,
		G: 0,
		B: 255,
	}
}

func getIntroColors() []sdl.Color {
	return []sdl.Color {
		getIntroColor(),
		getIntroOutlineColor(),
	}
}

func getHeartPositionPercentages() [8]sdl.FPoint {
	return [8]sdl.FPoint {
		sdl.FPoint {
			X: 0.165972222,
			Y: 0.127777778,
		},
		sdl.FPoint {
			X: 0.565972222,
			Y: 0.240740741,
		},
		sdl.FPoint {
			X: 0.156944444,
			Y: 0.7,
		},
		sdl.FPoint {
			X: 0.74375,
			Y: 0.483333333,
		},
		sdl.FPoint {
			X: 0.672222222,
			Y: 0.85,
		},
		sdl.FPoint {
			X: 0.130555556,
			Y: 0.506481481,
		},
		sdl.FPoint {
			X: 0.452083333,
			Y: 0.076851852,
		},
		sdl.FPoint {
			X: 0.842361111,
			Y: 0.175,
		},
	}
}

func getHeartPositions() [8]sdl.Point {
	var ret [8]sdl.Point
	var cents = getHeartPositionPercentages()

	for i := 0; i < len(cents); i++ {
		ret[i].X = int32(cents[i].X * gfxWindowWidth)
		ret[i].Y = int32(cents[i].Y * gfxWindowHeight)
	}

	return ret
}
