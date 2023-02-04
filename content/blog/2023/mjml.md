---
date: 2023-02-04
image: mjml_og.png
---

# MJML

![MJML logo](mjml.png)

Email clients are known for their poor (and inconsistent) support of HTML and CSS.
MJML is a markup language designed to reduce the pain of coding a responsive email.

https://mjml.io

https://documentation.mjml.io

https://github.com/mjmlio/mjml

Example MJML code:

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

Then running `mjml input.mjml -o output.html` will generate ugly HTML code,
that can be rendered correctly by email clients.
