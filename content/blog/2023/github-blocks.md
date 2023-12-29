---
date: 2023-03-30
refs:
  - blog/2020/graphviz.md
---

# GitHub Blocks

![GitHub Blocks](github-blocks.png)

GitHub launched a new Blocks platform in experimental techical preview.
Blocks tells GitHub how to render selected file or directory.
You may choose from list of available blocks or create your own.
Technically, they work as a React JSX (or TSX) components, rendered in iFrame. Note that such blocks can be embedded in any Markdown file.

https://blocks.githubnext.com

I made a simple block for rendering Graphviz DOT files.
It uses [graphviz-wasm](https://github.com/fabiospampinato/graphviz-wasm/)
package to run WASM port of Graphviz in browser 🤯.

https://github.com/chuhlomin/graphviz-block

![graphviz-block-demo](graphviz-block-demo.png)

**Update 2023-12-15**: The Blocks technical preview [ended](https://gist.github.com/idan/325676d192b32f169b032fde2d866c2c#github-next--technical-preview-sunsets).

#github #project
