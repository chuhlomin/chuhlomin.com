---
date: 2023-03-30
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
If you have access to technical preview, you can
[see it in action](https://blocks.githubnext.com/chuhlomin/graphviz-block-demo/blob/main/README.md).

https://github.com/chuhlomin/graphviz-block

![graphviz-block-demo](graphviz-block-demo.png)

#github #project
