---
date: 2020-06-02
---

# gthash

As a continuation of the Geohash topic, two years ago, guys wrote a library in TypeScript that makes a hash for coordinates and time.

https://github.com/ChrisChares/gthash

Based on the Z-order curve (like the geohash), the "geo-temporal" hashing algorithm can encode from kilometers and months to nanometers and pico-seconds (10-¹² s). Time can be coded within ±100,000 years of the beginning of the Unix Era.

https://en.wikipedia.org/wiki/Z-order_curve

I haven't figured out how it can be used yet. Even though it's short therefore seems ideal for QR code on some posters (time and place), it still seems more reasonable to use the [vEvents](https://icalendar.org/iCalendar-RFC-5545/3-6-1-event-component.html) format, which is easily read by phones as a calendar event.

#geo #qr
