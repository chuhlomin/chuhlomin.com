---
date: 2020-05-14
---

# Dive

If the word Docker is not entirely foreign to you,
I would like to advise a small utility that will simplify the Docker image analysis.

It shows the contents of the image layer by layer and highlights changes.
For example, you can use it to find unnecessary files to reduce the size of the image.

https://github.com/wagoodman/dive

<video width="720" height="450" controls>
  <source src="dive.mp4" type="video/mp4">
</video>

It's pretty easy run:

```bash
dive <your-image-tag>
```

But you can also run it as a Docker container if you don't want to install the application:

```bash
docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    wagoodman/dive:latest <dive arguments...>
```

Can you imagine, you can run a Dive container through Docker to analyze a Dive container.

#docker
