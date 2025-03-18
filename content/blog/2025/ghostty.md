---
date: 2025-03-18
refs:
  - blog/2020/iterm2.md
  - blog/2021/shaders.md
---

# Ghostty

![Ghostty](ghostty.png)

By now, many of engineers have probably heard about Ghostty, a relatively new terminal emulator that's been gaining popularity. It stands out for its exceptional speed, simplicity, and sensible default settings.

https://ghostty.org

Ghostty doesn't have a graphical settings interface. Instead, it's configured through a configuration file.

<details>
<summary>Here is my config</summary>
<pre><code>font-family = "TX-02 Condensed"
font-size = 15

copy-on-select = false
window-save-state = always
shell-integration-features = no-cursor

cursor-opacity = 0.8
cursor-style = block
cursor-style-blink = true
cursor-color = #d6d6d6

keybind = super+backspace=text:\x15

theme = Monokai Pro Spectrum
background = #222222
foreground = #f7f1ff
selection-background = #525053
selection-foreground = #f7f1ff
cursor-color = #bab6c0
palette = 0=#222222
palette = 1=#fc618d
palette = 2=#7bd88f
palette = 3=#fce566
palette = 4=#fd9353
palette = 5=#948ae3
palette = 6=#5ad4e6
palette = 7=#f7f1ff
palette = 8=#69676c
palette = 9=#fc618d
palette = 10=#7bd88f
palette = 11=#fce566
palette = 12=#fd9353
palette = 13=#948ae3
palette = 14=#5ad4e6
palette = 15=#f7f1ff

custom-shader = <path to shaders dir>/bloom.glsl
custom-shader = <path to shaders dir>/bettercrt.glsl
</code></pre>

</details>

If you're looking for help with your configuration, there's an [unofficial config tool](https://ghostty.zerebos.com/) and the [official reference documentation](https://ghostty.org/docs/config/reference).

Ghostty supports custom shaders, which can be stacked.
In configuration above I'm using two shaders: bloom and bettercrt.

![Ghostty Shaders](ghostty_screenshot.png)

https://github.com/chuhlomin/ghostty-shaders

#cli #app
