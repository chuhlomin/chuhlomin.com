---
created: 2020-04-22
---

# QR

I recently discovered that the search engine DuckDuckGo generates a QR code for a query that begins with "qr" (it has been doing this since 2013).

![DuckDuckGo QR demo](qr.jpeg "DuckDuckGo QR demo")

By the way, in addition to the plain text in the QR code can be written:

- Wi-Fi network data, including password (look for QR code on your modem): `WIFI:T:WPA;S:mynetwork;P:mypass;,`
- geolocation: `geo:40.71872,-73.98905,100,`
- events (name, start and end dates and times)
- malicious link,
- as well as many other things.

More information is collected here: https://github.com/zxing/zxing/wiki/Barcode-Contents

There are several QR code libraries for the Go programming language, and I used this one: https://github.com/liyue201/goqr

#qr #golang
