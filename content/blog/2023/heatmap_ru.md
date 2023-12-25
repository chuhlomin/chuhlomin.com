---
date: 2023-12-25
image: heatmap.png
---

# Тепловая карта

Тепловая карта это один из самых старых способов визуализации данных.
Исследуемые данные распологаются в таблице, а цвет ячейки определяется значением в ней: чем больше значение, тем ярче цвет.

Подобная карта используется на GitHub для визуализации активности пользователя.

Кстати, если данные по тому же принципу визулизируется на карте, то это называется [хороплетом](https://ru.wikipedia.org/wiki/Фоновая_картограмма).

Из интереса, я сделал такую карту на D3 для постов в этом блоге:

https://chuhlomin.com/experiments/heatmap/?lang=ru

<iframe src="/experiments/heatmap/?lang=ru"
    width="100%"
    height="auto"
    style="border: none;"
    onload="const viewBox = this.contentWindow.document.querySelector('svg').getAttribute('viewBox');
    const [x, y, width, height] = viewBox.split(',');
    this.style.aspectRatio = width / height;"
></iframe>

#visualization #blog
