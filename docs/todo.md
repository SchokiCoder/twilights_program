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

- add font outline for intro text
  (the purple outline)
- properly place text

- consider wag speed cap (max = as seen in source)
- add "wag wag wag..." background text

- add graphics
  see issue 1

- add music ?

- compare with source material
- set version 0.69.0
