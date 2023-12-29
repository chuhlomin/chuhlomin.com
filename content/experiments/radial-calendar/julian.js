/**
 * Julian date conversion functions
 *
 * https://github.com/ryanseys/lune/
 * 
 * Copyright 2004
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

'use strict'

function fromDate (date) {
  return date.getTime() / 86400000 + 2440587.5
}

function toDate (julian) {
  return new Date((julian - 2440587.5) * 86400000)
}

function calcMeanInstant(k, year) {
  const Y = (year - 2000) / 1000;
  switch (k) {
    case 0: // March equinox
      return 2451623.80984 +
        365242.37404 * Y +
        0.05169 * Math.pow(Y, 2) -
        0.00411 * Math.pow(Y, 3) -
        0.00057 * Math.pow(Y, 4);
    case 1: // June solstice
      return 2451716.56767 +
        365241.62603 * Y +
        0.00325 * Math.pow(Y, 2) +
        0.00888 * Math.pow(Y, 3) -
        0.0003 * Math.pow(Y, 4);
    case 2: //  September equinox
      return 2451810.21715 +
        365242.01767 * Y -
        0.11575 * Math.pow(Y, 2) +
        0.00337 * Math.pow(Y, 3) +
        0.00078 * Math.pow(Y, 4);
    case 3: //  December solstice
      return 2451900.05952 +
        365242.74049 * Y -
        0.06223 * Math.pow(Y, 2) -
        0.00823 * Math.pow(Y, 3) +
        0.00032 * Math.pow(Y, 4);
  }
};
