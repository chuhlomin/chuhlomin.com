---
date: 2022-11-25
image: ga.png
---

# Google Analytics Opt-out

<figure>
<img src="./ga.jpg" atl="Google Analytics (Иллюстрация)">
<figcaption>Картинка сгенерирована <a href="https://www.midjourney.com">Midjourney</a></figcaption>
</figure>

Если вас беспокоит, что Гугл аналитика собирает данные о вас,
то вот официальный способ отказаться от этого:

https://tools.google.com/dlpage/gaoptout

По сути это расширение браузера, которое добавляет тег на страницу:

```html
<script type="text/javascript">window["_gaUserPrefs"] = { ioo : function() { return true; } }</script>
```

#privacy #advice
