---
date: 2023-02-04
image: mjml_og.png
---

# MJML

![MJML logo](mjml.png)

Email клиенты известны своей плохой (и несогласованной) поддержкой HTML и CSS.
MJML – это язык разметки, разработанный для уменьшения боли при создании адаптивных электронных писем.

https://mjml.io

https://documentation.mjml.io

https://github.com/mjmlio/mjml

Пример кода MJML:

```html
<mjml lang="en">
  <mj-head>
    <mj-include path="./styles.css" type="css" />
    <mj-font name="MyFont" href="https://cdn.mysite.com/fonts/MyFont.css" />
    <mj-all font-family="MyFont,-apple-system,blinkmacsystemfont,roboto,helvetica neue,segoe ui,arial,sans-serif"></mj-all>
  </mj-head>
  <mj-body>
    <mj-section>
      <mj-column>
        <mj-include path="_header.mjml" />
        <mj-text>
          <p>Hello World!</p>
        </mj-text>
        <mj-include path="_footer.mjml" />
      </mj-column>
    </mj-section>
  </mj-body>
</mjml>
```

Затем запуск `mjml input.mjml -o output.html` сгенерирует ужасный HTML код,
который будет корректно отображен в почтовых клиентах.
