---
date: 2021-11-16
image: zoom.png
---

# Zoom-links

If you work remotely, you almost certainly use Zoom.
The experience of joining a meeting is something like this (at least for me):

1. you click on a link in your calendar/letter/slack/application,
2. opens a new tab in your browser (Safari),
3. you get an alert like "can I open Zoom?",
4. you press "yes," and Zoom opens.
5. The tab in the browser stays open.

It gets annoying after a while.

Roman Timushev wrote how you could get rid of extra steps with Hammerspoon and a Lua script:

https://timushev.com/posts/2020/12/05/conference-call-links/

Hammerspoon should be the default browser and works as a "proxy":
if a link looks like a Zoom-meeting link, then it will open Zoom right away,
and if it is not – it will be opened in the browser specified in the script.

I tried this method, and now I do not understand how I lived before.

PS: New blog tag! It will be for posts with recommended articles.

#link
