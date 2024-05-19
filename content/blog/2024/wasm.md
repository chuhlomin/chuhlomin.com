---
date: 2024-05-19
image: wasm.png
---

# WebAssembly

![WebAssembly logo](wasm.png "WebAssembly logo")

WebAssembly (WASM) is a binary code format for a stack-based virtual machine.
WASM is designed as a portable target for compiling high-level languages like C/C++/Rust/Go, enabling deployment on the web for both client and server applications.

https://webassembly.org/

All modern browsers support it, and it can also run on the server side, for example, on [Fastly Compute@Edge](https://docs.fastly.com/products/compute) or [Cloudflare Workers](https://developers.cloudflare.com/workers/runtime-apis/webassembly/#webassembly-wasm) platforms.

I rewrote a small project for generating PDF files with grid patterns from Go running on the server to Go-compiled-into-WASM running in the browser:
https://grid.chuhlomin.com/

#web #project #go
