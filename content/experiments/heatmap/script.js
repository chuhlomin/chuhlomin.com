/* jshint esversion: 6, browser: true */
/* globals d3 */

const font = "-apple-system, 'Segoe UI', 'Open Sans', Helvetica, Arial, sans-serif";

const dateValues = sample.map(dv => ({
  date: d3.timeDay(new Date(dv.date + "T05:00:00Z")),
  value: Number(dv.value),
  title: dv.title,
  url: dv.url
}));

const years = Array.from(
  d3.rollup(
    dateValues,
    v => v.length,
    d => d.date.getFullYear()
  ),
  ([year, count]) => ({ year, count })
).sort((a, b) => b.year - a.year);

const daysInWeek = 7;
const weeksInYear = 54;
const cellSize = 16;
const padding = cellSize;
const yearVerticalMargin = cellSize;
const yearHeight = cellSize * daysInWeek;
const yearWidth = cellSize * weeksInYear;
const leftMargin = cellSize * 4;

const totalWidth = leftMargin + yearWidth + padding;
const totalHeight = yearHeight * years.length + yearVerticalMargin * (years.length - 1) + 2*padding;

const formatDate = d3.utcFormat("%Y-%m-%d");
const colorFn = d3
  .scaleSequential(d3.interpolateBuGn)
  .domain([Math.floor(0), Math.ceil(1)]);
const format = d3.format("+.2%");

const svg = d3
  .create('svg')
  .attr("viewBox", [0, 0, totalWidth, totalHeight])
  .attr("preserveAspectRatio", "xMidYMid meet");

const group = svg.append("g");

// Year groups
const year = group
  .selectAll("g")
  .data(years)
  .join("g")
  .attr(
    "transform",
    (d, i) => `translate(${leftMargin}, ${(yearHeight + yearVerticalMargin) * i + padding})`
  )

// Year labels
year
  .append("text")
  .attr("x", -5)
  .attr("y", -2*cellSize)
  .attr("text-anchor", "end")
  .attr("class", "year")
  .attr("font-size", 16)
  .attr("font-weight", 550)
  .attr("transform", "rotate(270)")
  .attr("fill", "#222")
  .text(d => d.year);

// Mon, Wed, Fri labels
year
  .append("g")
  .attr("class", "weekdays")
  .attr("text-anchor", "middle")
  .selectAll("text")
  .data(weekDays)
  .join("text")
  .attr("x", -8)
  .attr("y", (d, i) => (i + 0.5) * cellSize)
  .filter(d => d) // if d == "", skip the label
  .attr("dy", "0.31em")
  .attr("font-size", 12)
  .attr("fill", "#222")
  .text(d => d);

// Day cells, empty
year
  .append("g")
  .selectAll("rect")
  .data(d => d3.timeDays(new Date(d.year, 0, 1), new Date(d.year + 1, 0, 1)))
  .join("rect")
  .attr("width", cellSize - 1)
  .attr("height", cellSize - 1)
  .attr("x", d => dayX(d) * cellSize + 0.5)
  .attr("y", d => dayY(d) * cellSize + 0.5)
  .attr("data-date", d => formatDate(d))
  .attr("data-x", d => dayX(d))
  .attr("data-y", d => dayY(d))
  .attr("fill", "#eee");

// Day cells, using sample data
year
  .append("g")
  .selectAll("rect")
  .data(d => dateValues.filter(v => v.date.getFullYear() === d.year))
  .join("a") // Wrap each rectangle with an <a> tag
  .on("mouseover", showTooltip)
  .on("mouseout", hideTooltip)
  .on("mousemove", moveTooltip)
  .attr("href", d => d.url)
  .attr("target", "_top")
  .append("rect") // Append a rectangle element inside the <a> tag
  .attr("data-date", d => formatDate(d.date))
  .attr("data-title", d => d.title)
  .attr("width", cellSize - 1)
  .attr("height", cellSize - 1)
  .attr("x", d => dayX(d.date) * cellSize + 0.5)
  .attr("y", d => dayY(d.date) * cellSize + 0.5)
  .attr("data-date", d => formatDate(d.date))
  .attr("data-x", d => dayX(d.date))
  .attr("data-y", d => dayY(d.date))
  .attr("fill", "#69C16D");

// Function to generate the path for each month
function pathMonth(t) {
  const x = (dayX(t) + 1) * cellSize;
  const y = dayY(t) * cellSize;
  
  // In the ideal case, the line will be straight.
  if (y == 0) {
    return `M${x-cellSize},0.5 V${yearHeight - 0.5}`;
  }
  
  // But if weeks are not aligned, the line will be slanted.
  return `M${x},0.5 V${y} H${x-cellSize} V${yearHeight - 0.5}`;
}

// Draw deviders between months
year
  .append("g")
  .attr("class", "months")
  .selectAll("path")
  .data(d => d3.timeMonths(new Date(d.year, 1, 1), new Date(d.year + 1, 0, 1)))
  .join("path")
  .attr("d", d => pathMonth(d))
  .attr("stroke", "#222")
  .attr("stroke-width", 1)
  .attr("fill", "none");

// ------------------------------

// Add tooltip element
const tooltip = svg
  .append("g")
  .attr("class", "tooltip");

// Tooltip collout
tooltip.append("path")
  .attr("fill", "#FEFECA")
  .attr("stroke", d3.color("#FEFECA").darker());

const tooltipText = tooltip.append("text")
  .attr("font-size", 16)
  .attr("font-family", font)
  .attr("fill", "#222");

tooltipText.append("tspan").attr("class", "date");
tooltipText.append("tspan").attr("class", "title");

const tooltipPaddingX = 8;
const tooltipPaddingY = 4;
const cursorOffset = 24;

function showTooltip(evt) {
  moveTooltip(evt);
  tooltipEl.classList.add("tooltip--visible");
}

function moveTooltip(evt) {
  var CTM = svgEl.getScreenCTM(); // coordinate transformation matrix
  // e – translation along x axis
  // f – translation along y axis
  // a – scaling along x axis
  // d – scaling along y axis

  tooltipElText.querySelector(".date").textContent = evt.target.dataset.date;
  tooltipElText.querySelector(".title").textContent = ": " + evt.target.dataset.title;
  var bbox = tooltipElText.getBBox();

  const mouse = {
    x: (evt.clientX - CTM.e) / CTM.a,
    y: (evt.clientY - CTM.f) / CTM.d
  };

  const pos = getTooltipPosition(CTM, bbox, mouse);
  const textRectCoords = getTextRectCoords(CTM, bbox, pos, mouse);
  const tooltipRectSize = {
    width: bbox.width + 2*tooltipPaddingX,
    height: bbox.height + 2*tooltipPaddingY
  };
  
  tooltipElText.setAttribute("x", textRectCoords.textX);
  tooltipElText.setAttribute("y", textRectCoords.textY);

  drawTooltipCallout(CTM, pos, textRectCoords, tooltipRectSize, evt.target);
}

function hideTooltip() {
  tooltipEl.classList.remove("tooltip--visible");
}

function getTooltipPosition(CTM, bbox, mouse) {
  var x = "right"
  var y = "bottom"

  // handle corner cases: too close to the right or bottom edge
  if (mouse.x + bbox.width + 2*tooltipPaddingX + cursorOffset / CTM.a > (totalWidth - padding / CTM.a)) {
    x = "left";
  }

  // handle corner cases: too close to the bottom edge
  if (mouse.y + bbox.height + 2*tooltipPaddingY + cursorOffset / CTM.d > (totalHeight - padding / CTM.d)) {
    y = "top";
  }

  return {x, y};
}

function getTextRectCoords(CTM, bbox, pos, mouse) {
  var textX;
  var textY;
  var rectX;
  var rectY;

  var cursorOffsetScaled = {
    x: cursorOffset,// / CTM.a,
    y: cursorOffset,// / CTM.d
  };

  switch (pos.x) {
    case "left":
      textX = mouse.x - bbox.width - cursorOffsetScaled.x;
      rectX = mouse.x - bbox.width - tooltipPaddingX - cursorOffsetScaled.x;
      break;
    default:
      textX = mouse.x + tooltipPaddingX + cursorOffsetScaled.x;
      rectX = mouse.x + cursorOffsetScaled.x;
  }

  switch (pos.y) {
    case "top":
      textY = mouse.y - cursorOffsetScaled.y;
      rectY = mouse.y - bbox.height - cursorOffsetScaled.y;
      break;
    default:
      textY = mouse.y + bbox.height + cursorOffsetScaled.y;
      rectY = mouse.y + cursorOffsetScaled.y;
  }

  return {textX, textY, rectX, rectY};
}

function drawTooltipCallout(CTM, pos, textRectCoords, tooltipRectSize, target) {
  const bounding = target.getBoundingClientRect();

  const from = {
    x: (bounding.x - CTM.e) / CTM.a,
    y: (bounding.y - CTM.f) / CTM.d
  };

  const to = {
    x: textRectCoords.rectX,
    y: textRectCoords.rectY
  };

  const rect = {
    centerX: (bounding.x + bounding.width / 2 - CTM.e) / CTM.a,
    centerY: (bounding.y + bounding.height / 2 - CTM.f) / CTM.d,
    width: bounding.width / CTM.a,
    height: bounding.height / CTM.d
  };

  switch (pos.x) {
    case "left":
      to.x += tooltipRectSize.width;
      break;
    default:
      from.x += bounding.width / CTM.a;
  }

  switch (pos.y) {
    case "top":
      to.y += tooltipRectSize.height;
      break;
    default:
      from.y += bounding.height / CTM.d;
  }

  var calloutPath;
  // Draw a callout box with a triangle pointing to the target.
  // There are 4 possible shapes, depending on the position of the target.
  const post = pos.x + "_" + pos.y;

  // find point where line from target to tooltip intersects with the rect side
  const kl = rect.height / 2;
  const kn = to.y - rect.centerY;
  const kj = to.x - rect.centerX;

  const shift = kl * kj / kn;

  if (pos.y == "bottom") {
    from.x = rect.centerX + shift;
  } else {
    from.x = rect.centerX - shift;
  }
  if (from.x > rect.centerX + rect.width / 2) {
    from.x = rect.centerX + rect.width / 2;
    from.y = rect.centerY + kl * kn / kj;
  } else if (from.x < rect.centerX - rect.width / 2) {
    from.x = rect.centerX - rect.width / 2;
    from.y = rect.centerY - kl * kn / kj;
  }

  switch (post) {
    case "left_top":
      calloutPath = `M${from.x},${from.y} L${to.x},${to.y - 7} L${to.x},${to.y - tooltipRectSize.height} L${to.x - tooltipRectSize.width},${to.y - tooltipRectSize.height} L${to.x - tooltipRectSize.width},${to.y} L${to.x - 7},${to.y} Z`;
      break;
    case "left_bottom":
      calloutPath = `M${from.x},${from.y} L${to.x},${to.y + 7} L${to.x},${to.y + tooltipRectSize.height} L${to.x - tooltipRectSize.width},${to.y + tooltipRectSize.height} L${to.x - tooltipRectSize.width},${to.y} L${to.x - 7},${to.y} Z`;
      break;
    case "right_top":
      calloutPath = `M${from.x},${from.y} L${to.x},${to.y - 7} L${to.x},${to.y - tooltipRectSize.height} L${to.x + tooltipRectSize.width},${to.y - tooltipRectSize.height} L${to.x + tooltipRectSize.width},${to.y} L${to.x + 7},${to.y} Z`;
      break;
    case "right_bottom":
      calloutPath = `M${from.x},${from.y} L${to.x},${to.y + 7} L${to.x},${to.y + tooltipRectSize.height} L${to.x + tooltipRectSize.width},${to.y + tooltipRectSize.height} L${to.x + tooltipRectSize.width},${to.y} L${to.x + 7},${to.y} Z`;
      break;
  }

  tooltipElPath.setAttribute("d", calloutPath);
}

// ------------------------------

document.body.append(svg.node());

const tooltipEl = document.querySelector(".tooltip");
const tooltipElText = tooltipEl.querySelector("text");
const tooltipElPath = tooltipEl.querySelector("path");
const svgEl = document.querySelector("svg");
