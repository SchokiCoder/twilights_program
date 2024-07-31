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

	gfxFontSize = 22
	gfxTextOutlineSize = 1
	gfxFirstBgLineYOffset = -9
	gfxBgMaxLines = 4.0
)

const (
	stdTickrate = 60.0
	stdTimescale = 1.0
)

// The durations and times are based on the video being 24 frames/second.
const (
	bgLineTravelTime = 1.0 / 24.0 * 15.0

	introDogTime = 1.0 / 24.0 * 25.0
	gameStartTime = introDogTime + (1.0 / 24.0 * 43.0)
	introLifetime = introDogTime + (1.0 / 24.0 * 45.0)

	eyeOpenedDuration = 1.0 / 24.0 * 5.0
	eyeClosedDuration = 1.0 / 24.0 * 3.0

	joyThroughWagsDelay = 1.0 / 24.0 * 10.0

	heartBigLifetime = 1.0 / 24.0 * 10.0
	heartSmallLifetime = 1.0 / 24.0 * 4.0
	heartGoneLifetime = 1.0 / 24.0 * 6.0
	heartLifetime = heartBigLifetime + heartSmallLifetime + heartGoneLifetime
)

const (
	wagsForHeart = 2.5
	wagsUntilJoy = 5
)

const (
	// pixel / second
	bgLineVelocity = gfxWindowHeight / bgLineTravelTime

	bgLineSpawnTime = bgLineTravelTime / gfxBgMaxLines
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
