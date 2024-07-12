# basics

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

- fix current issue of text being misplaced
  It is not misplaced, X and Y are correct,
  but it looks wrong due to scaling being impossible.
  Using BlitScaled on text surfaces causes panics:
  "Blit combination not supported",
  which is probably some scary sdl-magic not working.
  God, help me.

- add gimp layer to bmp export script

- scale pony art up and then down to smudge it a bit
  must be done at runtime, to accomodate different resolutions
- properly position and resize pony art

- add hearts art
- implement hearts art

- consider wag speed cap (max = as seen in source)

- make texts properly pixelated
- start 1st bg text line at proper position
  (uppermost visible pixelrow is at Y: 0)

- add graphics
  see issue 1

- add font outline for intro text
  (the purple outline)

- add music ?
  Split track into idle part und action part.
  The idle part plays as soon as the app starts,
  before and during the "run program?" confirmation.
  The action part starts as soon as the sdl window is opened.

- compare with source material
- set version 0.69.0
