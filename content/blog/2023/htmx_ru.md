---
date: 2023-09-24
---

# HTMX

![HTMX](htmx.png)

Bеб нуждается в упрощении. Все эти JavaScript-фреймворки, инструменты сборки и прочее затрудняют создание простого сайта. HTMX предлагает новый подход к веб-разработке.

https://htmx.org/

HTMX – это JavaScript-библиотека размером ≈14k, которая позволяет добавлять AJAX на HTML-страницы без написания JavaScript-кода. Для определения поведения страницы используются атрибуты HTML. Например, атрибут `hx-get` у ссылки меняет её поведение: при клике загружается содержимое с сервера и заменяет текущую страницу. А атрибут `hx-swap` настроит как произойдет эта замена.

Пример:

```html
<div hx-get="/hello" hx-swap="outerHTML">
  <p>Привет, мир!</p>
</div>
```

Недавно выпущена бесплатная книга Hypermedia Systems, посвященная созданию приложений на основе гипермедиа с использованием HTMX и др.

https://hypermedia.systems/

Пример использования HTMX c WebSockets с сервером на Go:

https://github.com/chuhlomin/htmx-websockets

Интерактивноть [поиска](https://chuhlomin.com/blog/search?lang=ru) в этом блоге реализована на HTMX.

#web #project #go