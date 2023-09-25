---
date: 2023-09-24
---

# HTMX

![HTMX](htmx.png)

The web needs to be simplified. All these JavaScript frameworks, build tools, and more make it difficult to create a simple website. HTMX is a new approach to web development.

https://htmx.org/

HTMX is a ~14k JavaScript library that allows you to add AJAX to HTML pages without writing JavaScript code. HTML attributes are used to define the behavior of the page. For example, a link's `hx-get` attribute changes its behavior: when clicked, it loads content from the server and replaces the current page. And the `hx-swap` attribute configures how this replacement will happen.

Example:

```html
<div hx-get="/hello" hx-swap="outerHTML">
  <p>Hello, World!</p>
</div>
```

A free Hypermedia Systems book on building hypermedia-based applications using HTMX and others has recently been released.

https://hypermedia.systems/

Example of using HTMX with WebSockets with a Go server:

https://github.com/chuhlomin/htmx-websockets

The interactivity of [search](https://chuhlomin.com/blog/search) in this blog is implemented on HTMX.

#web #project #go
