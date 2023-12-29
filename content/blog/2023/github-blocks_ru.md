---
date: 2023-03-30
refs:
  - blog/2020/graphviz_ru.md
---

# GitHub Blocks

![GitHub Blocks](github-blocks.png)

GitHub запустил новую платформу Blocks в экспериментальном техническом превью.
Блоки влияют на то как GitHub должен отобразить файл или каталог. Вы можете выбрать из списка доступных блоков или создать свой собственный. Технически они работают как компоненты React JSX (или TSX), отображаемые в iFrame. Такие блоки можно встраивать в любой Markdown файл.

https://blocks.githubnext.com

Я сделал простой блок для отображения файлов Graphviz DOT.
Он использует пакет [graphviz-wasm](https://github.com/fabiospampinato/graphviz-wasm/) чтобы запустить Graphviz в браузере (WASM порт 🤯).

https://github.com/chuhlomin/graphviz-block

![graphviz-block-demo](graphviz-block-demo.png)

**Обновление 2023-12-15**: Техническое превью Blocks [закончилось](https://gist.github.com/idan/325676d192b32f169b032fde2d866c2c#github-next--technical-preview-sunsets).

#github #project
