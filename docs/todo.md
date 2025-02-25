# additional improvements

- make heart lifetimes somewhat random?
  possible, but adds quite some complexity for little visual difference

- scale pony art up and then down to smudge it a bit?
  must be done at runtime, to accomodate different resolutions

- make texts properly pixelated?
  already are when upscaled...

- add purple shadow to intro text, casted from center of screen
  I have no idea how to do that.
  Maybe just cheat by hardwiring a purple tinted version of the texts
  where they are needed?

- the source's bg text also had an issue, in which every 4 lines the Y offset
  between lines is too small,
  but recreating this barely notable inconsistency
  on purpose adds complexity that I am not too fond of

- use fixed deltaTime as ebitengine suggests?
  I mean, it's not like this will ever lag anyway

-----

# Update-sized Update

- [x] remove unused go:generate comments

- [x] rework packaging
Fixes the AppImage, the install script, and icon pathing.
Also updates copyrights.

- [x] fix windows packaging
Since I can't reliably compile statically, I have to get the dlls. Fun.
Also rename package names to anything sensible.

- [x] go fmt

- [x] set version to 1.2

# flatpak polish v1.1

+ enforce new formatting for function declarations

Putting the first parameter already on a new line,
and the closing parenthesis with return type and opening bracket on a
new line fixes everything I disliked so far.
There is a clear cut between params and function start,
which often is var declarations.
No more weird alignment because of first param vs the other.

+ add icon
+ set new icon as window icon
+ fix icon not having a heart symbol

+ set version to 1.1

# basics v1.0

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

## implement hearts art (done)

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
+ add heart spawn que and decay

## tasks

+ add a delay for a hearts position to be open for a respawn

This makes it look a bit more random.

+ replace XCF files with just PNGs

So I thought of just dropping the XCFs and just having the PNGs outright,
but there is a reason why everything is layered in one file. Changing one layer
depends on other layers working with that change. That means
every time I change something I would have to import every corresponding PNG
and make it its own layer and- wait, GIMP has that exact funtion... oh.

On the other hoof, I could try to use GIMP's not at all documented scheme API
for the proper implementation of an export script...
Hmmm manual export after every change... or undocumented API,
in a language I never used and after that never would again...

Alright, bye XCFs.

+ add install script
+ add proper pathfinding for assets

+ add endgame statistics

Also fix first tick's delta being ridiculously high.

+ add -a --about arg

+ add font outline for intro text
  (the pink outline)

+ add startup hearts

Also make the heart lifetime math a ton simpler to read,
potentially also fixing it.
I know that NOW it works, the rest doesn't matter.

+ add DejaVuSansMono directly

+ add music
  Split track into idle part und action part.
  The idle part plays as soon as the app starts,
  before and during the "run program?" confirmation.
  The action part starts as soon as the sdl window is opened.

Also add copyright notices to "--about" for font and music files.

+ fix bonus round track not playing

Also fix install not copying images.
Also set main theme to loop.

+ change 1st bg text line to be at proper position

(Uppermost visible pixelrow's Y is at 0.)
Also increase font size to match original.

+ change tail'n'rump start position to match original

+ tweak heart lifetimes and spawn requirement to align more with original
  (original (wag ever 0.125 s, could barely maintain 2 hearts)

Also increase duration constants readability,
and add heart counter to end game score.

+ add demonstrational gif to README

Also update README's build dependencies.

+ fix asset paths being stupidly hardcoded and add support for os agnostic paths

Also remove Unix system install.

+ add mention of sdl2 and go-sdl2 to about info

+ add linux packaging to makefile

Also add License to packages.

+ add run.bat in package for windows to make this easily run in cmd.exe
  not needed lol

+ add ending game with escape key

+ add focussing sdl window after the begin confirmation

+ add end confirmation after score print

Also hide window as soon as possible,
and add file dependencies in Makefile.

+ add arg for timescale
  -t --timescale
+ add arg for tickrate
  -r --tickrate
+ add arg for disabling round clear sound
  -C --no-clear-sound
+ add arg for disabling confirmations
  -c --no-confirmations
+ add fullscreen arg
  -F --fullscreen
+ add help arg
  -h --help

+ add wag on pressing space key

+ fix score board confirmation on windows

It ignores a simple enter to continue, and a simple stdin question.
It needs the full spiel of a question in a loop,
because Windows is very fun to work with, and didn't cause a delay of days.
Shoot me now.

+ improve switch to fullscreen handling

+ add warning against unrecognized arguments

+ fix confirmation input possibly trailing into next loop iterations

+ set version to 1.0
