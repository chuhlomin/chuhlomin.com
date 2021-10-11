---
created: 2020-05-27
---

# Dark Mode

В macOS, Windows и других операционных системах, с относительно недавних пор, поддерживается «темный режим».
В случае с приложениями адаптация проходит гладко: они, как правило, используют системные библиотеки и получили поддержку «темных тем» без лишних усилий.

В случае с сайтами не все так просто.
Если вы делаете сайт, то можете воспользоваться особым CSS media query:
https://drafts.csswg.org/mediaqueries-5/#descdef-media-prefers-color-scheme

```css
body { color: black; background: white; }

@media (prefers-color-scheme: dark) {
  body { color: white; background: black; }
}
```

Но не у всех дойдут руки до обновления сайтов. В этом случае помогут расширения для браузеров, которые «подберут» цвета за владельцев сайта.

![DarkMode and DarkReader](darkmode.jpeg "DarkMode and DarkReader")

👌 Safari: DarkMode https://apps.apple.com/us/app/dark-mode-for-safari/id1397180934

🆗 Firefox (и другие): Dark Reader https://darkreader.org

#interface
