# Packaging

## Flathub
...I like flatpaks and flathub...  
This is probably my fault for being naive.  

The submission docs **lead** with this statement:  

> App submissions are extremely welcome

Not just welcome or very welcome.  
No... **extremely** welcome.  
Upon reading that, I figured that even a small, obscure, and silly app,
such as this here, would be welcome.  
I was wrong.  
The PR was practically denied.  

> It'd be easier to accept this if this had some more functionality.  
>   
> And yes, it should not spawn a terminal for confirmation.  

This app is feature complete and everything works as intended.  
So that's it. Game over.  

Again, this is probably my fault for being naive,
but those words of "extremely welcome" as the first to pop up in the docs,
really did wipe my initial worries.  

After spending hours carefully reading the docs and religiously following them,
after crafting the necessary files and learning their formats beforehand,
after days of back and forth over small issues with the PR,
getting still denied at last for something that **wouldn't** be changed hurt.  

...So after that major disappointment,
I had to get up and do the best out of it.  
Flathub isn't the only way after all.  
There are still

## AppImages

I was hoping for another kind of savior, but you'll have to do.  

### docs
The docs suck.  
I tried to create the AppDir manually,
and they tell you different information in different places about how to do
that.  

[Manual packaging](https://docs.appimage.org/packaging-guide/manual.html#ref-manual)
says to:  

> Create an AppDir structure that looks (as a minimum) like this:  
>   
> MyApp.AppDir/  
> MyApp.AppDir/AppRun  
> MyApp.AppDir/myapp.desktop  
> MyApp.AppDir/myapp.png  
> MyApp.AppDir/usr/bin/myapp  
> MyApp.AppDir/usr/lib/libfoo.so.0  

but in [AppDir specification](https://docs.appimage.org/reference/appdir.html#general-description)
it says:  

> myapp.desktop  
>   
> A desktop file located in the root directory, describing the payload application. As AppImage is following the principle one app = one file, one desktop file is enough to describe the entire AppImage. There MUST NOT be more than one desktop file in the root directory. The name of the file doesn’t matter, as long as it carries the .desktop extension. Can be a symlink to subdirectories such as usr/share/applications/...  

So suddenly you _can_ put it in share/applications, where it belongs,
and then make a symlink. Why not just instruct like that outright?  

Also, sometimes the docs just mention something,
and then refuse to elaborate to the needed degree.  
Remember the previous minimum AppDir structure I quoted above.  

...What is AppRun?  

> The AppRun file can be a script or executable. It sets up required environment variables such as $PATH and launches the payload application. You can write your own, but in most cases it is easiest (and most error-proof) to use a precompiled one from this repository.

It says right under the previous part.  
Oh cool, but... what is "this repository"?  
I am in the [docs](https://docs.appimage.org/packaging-guide/manual.html),
and there is no link...  
So what and where?  

But WAIT there is more!  
The docs are outdated too. Woo!  
Eventually through dumb luck with duckduckgo and _probably_ StackOverflow,
I found those magic precompiled AppRun binaries:  
<https://github.com/AppImage/AppImageKit/releases>

...aaand it's obsolete. Awesome!  

### dIsTrO aGnOsTiC pAcKaGeS
Why do devs bother with AppImages?  
It's the _somewhat_ false promise of distro agnostic packages.  

It turns out the very "distro agnostic package format" has packaging tools,
which **are** packaged "distro agnostic", and they themselves
[fail at being distro agnostic](https://github.com/linuxdeploy/linuxdeploy/issues/272).  

It's not just me on Fedora right now with that one package.  
[It's all over the place](https://ludditus.com/2024/10/31/appimage/),
and I already knew the moment I saw a blacklist for libraries mentioned,
upon running linuxdeploy, this wouldn't go well.  

## Static linking

After days of compressed effort in packaging distro agnostic, I give up.  
Why bother with all that, when I can just link statically,
and just like I would with an appImage, hope for the best?  
The tools for that, being just the go compiler and tar, which I already have.  
Put a simple install.sh into the package, and done.  

Except go-sdl2 doesn't allow static compiles or it doesn't work with libasound.  
...
...
...
...
...
You must find the heart of evil and drive a stake through it.  
