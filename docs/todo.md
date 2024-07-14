# basics

## tasks

+ add cli
+ add window
+ add proper query loop

Also add theme method to application.

+ rewrite in Go with fyne

See docs/rewrite.md for the why.

+ add theme

+ add timed text logic for intro
+ properly time intro

+ implement timer GameStartTime
+ implement timer IntroLifetime
+ implement blink logic
  (implement duration EyeOpenedDuration
  implement duration EyeClosedDuration)

+ add wag ability logic
+ add joy logic

+ replace fyne with sdl
  cuz fyne is too restrictive
  problems with that:
  	+ get fonts reliably: make a font bitmap
  	+ can ttf do font shadow: ttf.SetOutline()

+ add intro text

For now, with just some font as mock up.

+ add bg text with animation

For now, with just some font as mock up.
Also fix wrong aspect ratio (now 4:3), and tweak font size.

+ add pony art
+ implement pony art
+ fix blink logic still running after joy has been reached

+ fix timescale not applying to intro, blink and joy delay timers

Also largely convert all nanosecond(int) based values to second(float) ones.

+ fix wags being counted during intro

+ scale everything down to pony size by default and add scaling for pony only

Everything is now scaled 1.0 towards native pony image size.
Also add missing err checks for blit function calls.
Scaling for pony only, becaus using BlitScaled on texts causes a panic.

+ add proper sprites

Instead of just blitting surfaces, create and render textures.

+ add gfx scaling
+ properly position pony art

Also improve delta calculation, generally use more floats,
and use renderer.Clear instead of FillRect.

+ add hearts art

## implement hearts art (wip)

### Source material observations

3 set positions around pony, with 2 variants.
A heart livespan or spawntime seems independent of wags,
as there is sometimes no delay between spawn and sometimes there is.
Each heart lives big for 8 or 10 frames and small for 4 or 6 frames.
A new heart is spawned every 10 or 12 frames.
It seems all random.

### Implementation

Changing hearts to be feedback for the player doing well would fit better for
an interactive experience, and would allow me to avoid this randomness of the
source material.
2 wags will be required for the spawn of a heart.
There will be set positions for hearts,
and heart spawns will be queued for those positions,
trying to at first spawn at the center positions,
then moving to outer spawn positions.
Wagging faster means having more hearts maintained on the screen.

### subtasks

+ add heart position values
- add heart que
- add heart decay

## tasks

- add gimp-layer-to-bmp-export-script

- scale pony art up and then down to smudge it a bit?
  must be done at runtime, to accomodate different resolutions

- consider wag speed cap (max = as seen in source)

- make texts properly pixelated
- start 1st bg text line at proper position
  (uppermost visible pixelrow is at Y: 0)

- add graphics
  see issue 1

- add font outline for intro text
  (the purple outline)

- directly ship with current font ?

- add music ?
  Split track into idle part und action part.
  The idle part plays as soon as the app starts,
  before and during the "run program?" confirmation.
  The action part starts as soon as the sdl window is opened.

- compare with source material
- set version 0.69.0
