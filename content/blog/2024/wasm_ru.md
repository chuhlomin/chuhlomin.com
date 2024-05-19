---
date: 2024-05-19
image: wasm.png
---

# WebAssembly

![Логотип WebAssembly](wasm.png "Логотип WebAssembly")

WebAssembly (WASM) – это бинарный формат кода для стековой виртуальной машины.
Программирование идёт на обычных статически типизированных языках, таких как Си, C++, C#, Rust, Go. После компиляции кода на этих языках в WASM, его можно запустить в браузере или на сервере.

https://webassembly.org/

Его поддерживают все современные браузеры. Также он может работать на стороне сервера, например, на платформах [Fastly Compute@Edge](https://docs.fastly.com/products/compute) или [Cloudflare Workers](https://developers.cloudflare.com/workers/runtime-apis/webassembly/#webassembly-wasm).

Я переписал небольшой проект по генерации PDF-файлов с шаблонами сеток с Go, работающего на сервере, на Go, скомпилированный в WASM, работающий в браузере:
https://grid.chuhlomin.com/

#web #project #go

