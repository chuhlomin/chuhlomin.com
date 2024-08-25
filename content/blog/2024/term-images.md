---
date: 2024-08-24
---

# Terminal Images Protocols

![Terminal with an image](term-images.png)

Some terminals are capable of displaying images.
It is rarely needed, but sometimes it can be useful.

[iTerm2](https://iterm2.com/) introduced an
[Inline Images Protocol](https://iterm2.com/documentation-images.html),
that is now supported by [WezTerm](https://wezfurlong.org/wezterm/) as well.
A running program can output a Base-64 encoded image to the terminal
(wrapped in special control sequences), and the terminal will display this image.

```
<ESC>]1337;File=width=60;height=auto:<base-64 encoded image>^G
```

[Kitty](https://sw.kovidgoyal.net/kitty/) terminal came up with
[Terminal Graphics Protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/),
which has more capabilities. For example, you can specify the position of the image
on the screen, display animations, and remove previously displayed images.

```
<ESC>_Ga=T,f=100,t=d,m=1;<base-64 encoded image><ECS>\
```

For Go, I recommend the [`rasterm`](https://github.com/BourgeoisBear/rasterm) library,
which implements both protocols.

#cli #image #go
