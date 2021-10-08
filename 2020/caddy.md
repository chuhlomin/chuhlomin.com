---
created: 2020-06-23
---

# Caddy (v2)

![Caddy logo](caddy.png "Caddy logo")

Many people must have heard about nginx - a high-performance web server used worldwide, written by Igor Sysoev (2004).
It works great, and it would seem that there is nothing more to be desired. On the other hand, it's been years, and now I think many things would be done differently.

And then the increasingly popular Caddy, written in Go, comes on the scene:

- also an open-source project
- with simple, minimalist configurations,
- Automatically gets and updates certificates from Let's Encrypt, auto-redirect HTTP→HTTPS
- and a lot more tasty stuff.

https://caddyserver.com

For example, this is what the config for one project looks like right now:

```
ts.chuhlomin.com {
  reverse_proxy timestamp:80
}
```

That's it, I don't have to think about certificates anymore.

The caddy-docker-proxy project deserves a special mention, as it automatically updates the Caddy config (the server itself "picks up" the changes) based on the labels of the containers:
https://github.com/lucaslorentz/caddy-docker-proxy

#app #ops #go
