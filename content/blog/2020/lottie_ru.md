---
date: 2020-04-27
---

# Lottie

![Lottie logo](lottie.png "Lottie logo")

Продолжим про Телеграм стикеры.
Вместо WebP для анимированных стикеров был выбран формат под названием Lottie, разработанный в AirBnB (помните, раньше еще можно было путешествовать в другие страны?).

https://airbnb.io/lottie/

Основное отличие – это векторная анимация, а следовательно, как правило, имеет меньший размер. Примеры анимации: https://lottiefiles.com

Основной способ создания таких анимаций – в платном Adobe After Effects с плагином Bodymovin который создает JSON файл. Любой анимированный стикер в Телеграме это заархивированный JSON файл в Lottie-формате.

https://core.telegram.org/animated_stickers

Для показа анимаций на сайте предоставляется JS библиотека
https://github.com/airbnb/lottie-web

Онлайн редактор для создания анимаций:
https://www.lottielab.com

**Обновление 2024-09-01**: Есть [плагин для Figma](https://lottiefiles.com/plugins/figma) для создания Lottie анимаций.

#telegram
