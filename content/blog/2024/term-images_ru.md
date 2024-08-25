---
date: 2024-08-24
---

# Протоколы вывода изображений в терминале

![Терминал с изображением](term-images.png)

Некоторые терминалы умеют выводить изображения.
Это редко бывает нужно, но иногда может пригодиться.

В [iTerm2](https://iterm2.com/) впервые появился протокол
[Inline Images Protocol](https://iterm2.com/documentation-images.html),
который теперь поддерживается и [WezTerm](https://wezfurlong.org/wezterm/).
Запущенная программа может вывести Base-64 закодированное изображение
в терминал (обёрнутое в специальные управляющие последовательности),
и терминал отобразит это изображение.

```
<ESC>]1337;File=width=60;height=auto:<base-64 encoded image>^G
```

Для терминала [Kitty](https://sw.kovidgoyal.net/kitty/) разработан протокол
[Terminal Graphics Protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/),
у которого больше возможностей. Например, можно указать положение изображения
на экране, выводить анимацию и удалять ранее выведенные изображения.

```
<ESC>_Ga=T,f=100,t=d,m=1;<base-64 encoded image><ECS>\
```

Для Go Рекомендую библиотеку [`rasterm`](https://github.com/BourgeoisBear/rasterm), которая реализует оба протокола.

#cli #image #go
