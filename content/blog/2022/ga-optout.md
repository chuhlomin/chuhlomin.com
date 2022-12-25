---
date: 2022-11-25
image: ga.png
---

# Google Analytics Opt-out

<figure>
<img src="./ga.jpg" atl="Google Analytics Illustration">
<figcaption>Image generated by <a href="https://www.midjourney.com">Midjourney</a></figcaption>
</figure>

If you are concerned about Google Analytics collecting data about you, here is the official way to opt out:

https://tools.google.com/dlpage/gaoptout

It is a browser extension that adds a tag to the page:

```html
<script type="text/javascript">window["_gaUserPrefs"] = { ioo : function() { return true; } }</script>
```

#privacy #advice