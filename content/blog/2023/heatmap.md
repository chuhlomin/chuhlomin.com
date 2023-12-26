---
date: 2023-12-25
image: heatmap.png
---

# Heatmap

A heatmap is one of the oldest ways to visualize data.
The data is arranged in a table, and its value determines the color of the cell: the greater the value, the brighter the color.

A similar map is used on GitHub to visualize user activity.

By the way, if the data is visualized on a map following the same principle, it is called a [choropleth](https://en.wikipedia.org/wiki/Choropleth_map).

Out of interest, I made such a map on D3 for posts in this blog:

https://chuhlomin.com/experiments/heatmap/

<iframe src="/experiments/heatmap/"
    width="100%"
    height="auto"
    style="border: none;"
    onload="const viewBox = this.contentWindow.document.querySelector('svg').getAttribute('viewBox');
    const [x, y, width, height] = viewBox.split(',');
    this.style.aspectRatio = width / height;"
></iframe>

#visualization #blog
