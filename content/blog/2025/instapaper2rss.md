---
date: 2025-02-15
image: instapaper2rss_og.png
refs:
  - blog/2021/atom.md
---

# Instapaper2RSS

There are many ways to build your "read later" list. Safari Reading List, Pocket, Instapaper, etc.
But I always struggled with accessing them later, from another device, perhaps even offline, and ideally, in a RSS reader.

So I built a small tool that bridges this gap – Instapaper2RSS.
It's a Go application that runs in GitHub Action, fetching articles from my Instapaper
account and generating an Atom feed that can be subscribed to in any RSS reader.

https://github.com/chuhlomin/instapaper2rss

Key points:

- The app authenticates with Instapaper API using OAuth
- Fetches new bookmarks incrementally (only ones added since last run)
- Stores state in BoltDB database
- Generates an Atom feed with full article content
- Uploads everything to Cloudflare R2 storage

My feed to far: https://instapaper.chuhlomin.com/atom.xml

#project #go
