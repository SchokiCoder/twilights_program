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
- add mock up hearts

- scale pony art up and then down to smudge it a bit?
  do this at runtime !

- timescale and delta are not used everywhere
  The timers fully ignore the timescale.
  time.Since() and co must practically be banned.

- consider wag speed cap (max = as seen in source)

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
