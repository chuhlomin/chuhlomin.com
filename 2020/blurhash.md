---
created: 2020-04-20
---

# Blurhash

Today I want to talk about Blurhash.  
The idea is quite simple: an image can be "shrunk" to a particular hash represented as a short string.
From this string, you can generate a blurred image that remotely resembles the original.
Such a blurred image can be used as a muffler while the main image is loading or to hide the main image for other reasons (paid and/or adult content).

![Blurhash demo 1](blurhash.jpeg "Blurhash demo 1")

The promo page can explain it all more clearly: https://blurha.sh  
Repository: https://github.com/woltapp/blurhash

I experimented a bit with a Go library and made a CLI for it: https://github.com/chuhlomin/go-blurhash-experiment

[![Blurhash demo 2](blurhash.png "Blurhash demo 2")](blurhash.png)

Photo by Melissa Keizer on Unsplash https://unsplash.com/@keizgoesboom

#go #image
