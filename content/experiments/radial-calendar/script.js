/* jshint esversion: 6, browser: true */
/* globals d3 */

const sectionHeight = 20;
const padding = 20;
const radius = 500;
const innerRadius = radius - padding - sectionHeight;
const outerRadius = radius - padding;
const moonRadius = innerRadius - padding/2;
const solsticesRadius = innerRadius - padding;
const font = '-apple-system, "Segoe UI", "Open Sans", Helvetica, Arial, sans-serif';

const svg = d3
  .create('svg')
  .attr("viewBox", [0, 0, radius*2, radius*2])
  .attr("preserveAspectRatio", "xMidYMid meet");

// now in local time
const now = new Date();
const year = now.getFullYear();

// create array of days in year to display in D3
const days = d3.timeDays(new Date(year, 0, 1), new Date(year + 1, 0, 1));
// count days in each month
const monthDays = d3.rollup(days, v => v.length, d => d.getMonth());

// draw all days in a year as a circle
svg
  .append('g')
  .attr('fill', 'none')
  .attr('stroke', 'black')
  .attr('stroke-width', 0.25)
  .attr('stroke-linejoin', 'round')
  .attr('transform', `translate(${radius},${radius})`)
  .selectAll('path')
  .data(days)
  .join('path')
  // if day is weekend, use gray background color
  .attr('fill', d => d.getDay() === 0 || d.getDay() === 6 ? '#ddd' : 'none')
  .attr('d', d3.arc()
    .innerRadius(radius - padding - sectionHeight)
    .outerRadius(radius - padding)
    .startAngle(d => d3.timeDay.count(d3.timeYear(d), d) * 2 * Math.PI / days.length)
    .endAngle(d => (d3.timeDay.count(d3.timeYear(d), d) + 1) * 2 * Math.PI / days.length)
    .padRadius(radius - sectionHeight)
  );

// fill daysBeforeMonth dictionary with days before each month
const daysBeforeMonth = new Map();
daysBeforeMonth.set(-1, 0);
for (let i = 0; i < 12; i++) {
  let count = monthDays.get(i);
  if (i > 0) {
    count += daysBeforeMonth.get(i - 1);
  }

  daysBeforeMonth.set(i, count);
}

// draw all months in a year as a circle
svg
  .append('g')
  .attr('fill', 'none')
  .attr('stroke', 'black')
  .attr('stroke-width', 1.5)
  .attr('stroke-linejoin', 'round')
  .attr('transform', `translate(${radius},${radius})`)
  .selectAll('path')
  .data(d3.timeMonths(new Date(year, 0, 1), new Date(year + 1, 0, 1)))
  .join('path')
  .attr('d', d3.arc()
    .innerRadius(radius - padding - sectionHeight)
    .outerRadius(radius - padding)
    .startAngle(d => daysBeforeMonth.get(d.getMonth() - 1) * 2 * Math.PI / days.length)
    .endAngle(d => daysBeforeMonth.get(d.getMonth()) * 2 * Math.PI / days.length)
    .padRadius(radius - sectionHeight)
  );

// add SVG path for each month
d3.timeMonths(new Date(year, 0, 1), new Date(year + 1, 0, 1)).forEach((d, i) => {
  const startAngle = daysBeforeMonth.get(d.getMonth() - 1) * 2 * Math.PI / days.length;
  const endAngle = daysBeforeMonth.get(d.getMonth()) * 2 * Math.PI / days.length;
  const path = d3.path();
  path.moveTo(radius, radius);

  // draw arc where 0 is at 12 o'clock
  // change text direction for months from 3 to 9
  if (d.getMonth() >= 3 && d.getMonth() <= 8) {
    path.arc(radius, radius, radius + 4, endAngle - Math.PI / 2, startAngle - Math.PI / 2, true);
  } else {
    path.arc(radius, radius, radius - 10, startAngle - Math.PI / 2, endAngle - Math.PI / 2);
  }
  path.lineTo(radius, radius);
  path.closePath();
  svg.append("path")
    .attr("id", "month" + i)
    .attr("d", path)
    .attr("fill", "none")
    .attr("stroke", "none");
});

// add month labels, curved text along the circle
svg
  .append('g')
  .attr('font-family', font)
  .attr('font-size', 20)
  .attr('text-anchor', 'middle')
  .selectAll('text')
  .data(d3.timeMonths(new Date(year, 0, 1), new Date(year + 1, 0, 1)))
  .join('text')
  .append('textPath')
  .attr('xlink:href', (d, i) => "#month" + i)
  .attr('startOffset', '50%')
  .text(d => months[d.getMonth()]);

// add text path for day labels inside arc
days.forEach((d, i) => {
  // if day is 10th of month, 20th or last day of month, add label
  if (!(d.getDate() === 10 || d.getDate() === 20 || d.getDate() === monthDays.get(d.getMonth()))) {
    return;
  }
  
  const halfDayAngle = Math.PI / days.length / 2;
  const startAngle = d3.timeDay.count(d3.timeYear(d), d) * 2 * Math.PI / days.length;
  const endAngle = (d3.timeDay.count(d3.timeYear(d), d) + 1) * 2 * Math.PI / days.length;
  const middleAngle = startAngle + (endAngle - startAngle) / 2;
  const path = d3.path();
  
  // change direction of text for days from 6 to 12 month
  if (d.getMonth() >= 6 && d.getMonth() <= 11) {
    const middleX = radius + outerRadius * Math.cos(middleAngle - halfDayAngle - Math.PI / 2);
    const middleY = radius + outerRadius * Math.sin(middleAngle - halfDayAngle - Math.PI / 2);
    const middleX2 = radius + innerRadius * Math.cos(middleAngle - halfDayAngle - Math.PI / 2);
    const middleY2 = radius + innerRadius * Math.sin(middleAngle - halfDayAngle - Math.PI / 2);
    path.moveTo(middleX, middleY);
    path.lineTo(middleX2, middleY2);
  } else {
    const middleX = radius + innerRadius * Math.cos(middleAngle + halfDayAngle - Math.PI / 2);
    const middleY = radius + innerRadius * Math.sin(middleAngle + halfDayAngle - Math.PI / 2);
    const middleX2 = radius + outerRadius * Math.cos(middleAngle + halfDayAngle - Math.PI / 2);
    const middleY2 = radius + outerRadius * Math.sin(middleAngle + halfDayAngle - Math.PI / 2);
    path.moveTo(middleX, middleY);
    path.lineTo(middleX2, middleY2);
  }
  svg.append("path")
    .attr("id", "day" + d.toLocaleString('default', { month: 'long' }) + d.getDate())
    .attr("d", path)
    .attr("fill", "none")
    .attr("stroke", "none");

  // add text path for day labels in the middle of path
  svg
    .append('g')
    .attr('font-family', font)
    .attr('font-size', 6)
    .attr('text-anchor', 'middle')
    .append('text')
    .append('textPath')
    .attr('xlink:href', "#day" + d.toLocaleString('default', { month: 'long' }) + d.getDate())
    .attr('startOffset', '50%')
    .text(d.getDate());
});

// draw clock hand pointing to current time in the year
const day = d3.timeDay.count(d3.timeYear(now), now);
// get number of hours in the day in local time
const hours = now.getHours() + now.getMinutes() / 60 + now.getSeconds() / 3600;
const dayAngle = day * 2 * Math.PI / days.length + hours * 2 * Math.PI / days.length / 24;

// add to svg
svg
  .append('g')
  .attr('transform', `translate(${radius}, ${radius})`)
  .append('line')
  .attr('x1', 0)
  .attr('y1', 0)
  .attr('x2', innerRadius * Math.cos(dayAngle - Math.PI / 2))
  .attr('y2', innerRadius * Math.sin(dayAngle - Math.PI / 2))
  .attr('stroke', 'black')
  .attr('stroke-width', 1);

// add a small circle at the end of the hand
svg
  .append('g')
  .attr('transform', `translate(${radius}, ${radius})`)
  .append('circle')
  .attr('cx', innerRadius * Math.cos(dayAngle - Math.PI / 2))
  .attr('cy', innerRadius * Math.sin(dayAngle - Math.PI / 2))
  .attr('r', 2)
  .attr('fill', 'black');

// add a big circle in the middle
svg
  .append('g')
  .attr('transform', `translate(${radius}, ${radius})`)
  .append('circle')
  .attr('cx', 0)
  .attr('cy', 0)
  .attr('r', 4)
  .attr('fill', 'black');

// add dotted lines from center to beginning of each month
svg
  .append('g')
  .attr('transform', `translate(${radius}, ${radius})`)
  .selectAll('line')
  .data(d3.timeMonths(new Date(year, 0, 1), new Date(year + 1, 0, 1)))
  .join('line')
  .attr('x1', 0)
  .attr('y1', 0)
  .attr('x2', d => innerRadius * Math.cos(daysBeforeMonth.get(d.getMonth()) * 2 * Math.PI / days.length - Math.PI / 2))
  .attr('y2', d => innerRadius * Math.sin(daysBeforeMonth.get(d.getMonth()) * 2 * Math.PI / days.length - Math.PI / 2))
  .attr('stroke', 'black')
  .attr('stroke-width', 0.25)
  .attr('stroke-dasharray', '1 1');

// add moon phases
var fullMoons = phaseRange(
  new Date(year, 0, 1),
  new Date(year + 1, 0, 1),
  FULL
);

// add full moon circles
svg
  .append('g')
  .attr('transform', `translate(${radius}, ${radius})`)
  .selectAll('circle')
  .data(fullMoons)
  .join('circle')
  .attr('cx', d => moonRadius * Math.cos(d3.timeDay.count(d3.timeYear(d), d) * 2 * Math.PI / days.length - Math.PI / 2))
  .attr('cy', d => moonRadius * Math.sin(d3.timeDay.count(d3.timeYear(d), d) * 2 * Math.PI / days.length - Math.PI / 2))
  .attr('r', 2)
  .attr('fill', 'white')
  .attr('stroke', 'black')
  .attr('stroke-width', 0.5);

// add summer and winter Solstice
var solstices = [
  calcMeanInstant(0, year),
  calcMeanInstant(1, year),
  calcMeanInstant(2, year),
  calcMeanInstant(3, year)
].map(toDate);

// add solstice circles
svg
  .append('g')
  .attr('transform', `translate(${radius}, ${radius})`)
  .selectAll('circle')
  .data(solstices)
  .join('circle')
  .attr('cx', d => solsticesRadius * Math.cos(d3.timeDay.count(d3.timeYear(d), d) * 2 * Math.PI / days.length - Math.PI / 2))
  .attr('cy', d => solsticesRadius * Math.sin(d3.timeDay.count(d3.timeYear(d), d) * 2 * Math.PI / days.length - Math.PI / 2))
  .attr('r', 2)
  .attr('fill', 'black');

document.body.append(svg.node());
