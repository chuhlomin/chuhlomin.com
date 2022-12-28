/* jshint const: true, esversion: 6, undef: true */
const sectionHeight = 20;
const padding = 20;
const radius = 500;

const svg = d3
  .create('svg')
  .attr("viewBox", [0, 0, radius*2, radius*2])
  .attr("preserveAspectRatio", "xMidYMid meet")
  .style("font", "10px sans-serif");

const year = 2023;

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
    .startAngle(d => {
      let daysBeforeMonth = 0;
      for (let i = 0; i < d.getMonth(); i++) {
        daysBeforeMonth += monthDays.get(i);
      }
      return daysBeforeMonth * 2 * Math.PI / days.length;
    })
    .endAngle(d => {
      let daysBeforeEndOfMonth = 0;
      for (let i = 0; i <= d.getMonth(); i++) {
        daysBeforeEndOfMonth += monthDays.get(i);
      }
      return daysBeforeEndOfMonth * 2 * Math.PI / days.length;
    })
    .padRadius(radius - sectionHeight)
  );

// add SVG path for each month
d3.timeMonths(new Date(year, 0, 1), new Date(year + 1, 0, 1)).forEach((d, i) => {
  let daysBeforeMonth = 0;
  for (let j = 0; j < d.getMonth(); j++) {
    daysBeforeMonth += monthDays.get(j);
  }
  let daysBeforeEndOfMonth = 0;
  for (let j = 0; j <= d.getMonth(); j++) {
    daysBeforeEndOfMonth += monthDays.get(j);
  }
  const startAngle = daysBeforeMonth * 2 * Math.PI / days.length;
  const endAngle = daysBeforeEndOfMonth * 2 * Math.PI / days.length;
  console.log(d.toLocaleString('default', { month: 'long' }), startAngle, endAngle);
  const path = d3.path();
  path.moveTo(radius, radius);
  // draw arc where 0 is at 12 o'clock
  path.arc(radius, radius, radius - 10, startAngle - Math.PI / 2, endAngle - Math.PI / 2);
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
  .attr('font-family', 'sans-serif')
  .attr('font-size', 20)
  .attr('text-anchor', 'middle')
  .selectAll('text')
  .data(d3.timeMonths(new Date(year, 0, 1), new Date(year + 1, 0, 1)))
  .join('text')
  .append('textPath')
  .attr('xlink:href', (d, i) => "#month" + i)
  .attr('startOffset', '50%')
  .text(d => d.toLocaleString('default', { month: 'long' }));

document.body.append(svg.node());
