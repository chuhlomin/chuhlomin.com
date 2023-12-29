// Source: https://github.com/ebradyjobory/eclipse-explorer/blob/master/JSEX/SE2901.js 
//
// Besselian Elements for Solar Eclipses from:
//   "Five Millennium Canon of Solar Eclipses: -1999 to +3000",
//      Fred Espenak and Jean Meeus, NASA/TP-2006-214141, October 2006
//
const eclipseBesselianElements = new Array(
  2460054.679120, // 2023  4 20
  2460232.250470, // 2023 10 14
  2460409.262840, // 2024  4  8
  2460586.282098, // 2024 10  2
  2460763.950417, // 2025  3 29
  2460940.321576, // 2025  9 21
);

function calcSolarEclipses(year) {
  // calculate solar eclipse based on Besselian elements from SE2001 constant
  var eclipseList = [];

  for (var i = 0; i < eclipseBesselianElements.length; i++) {
    var date = toDate(eclipseBesselianElements[i]);
    if (date.getUTCFullYear() != year) {
      continue;
    }

    eclipseList.push(date);
  }

  return eclipseList;
}
