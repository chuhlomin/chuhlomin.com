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
