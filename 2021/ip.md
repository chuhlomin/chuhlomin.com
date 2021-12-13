---
date: 2021-12-06
---

# IP

Every device on the Internet has an IP address.
It can be an IPv4 version (e.g. 1.1.1.1) or IPv6 (e.g. 2001:db8:0:1234:0:567:8:1).

Sometimes you need to know your IP or see information about another IP.  
There are several services for this purpose:

* https://ifconfig.co
* https://httpbin.org/ip
* https://ipinfo.io
* https://whatismyipaddress.com

I wanted something more minimalistic, without styles, to use it from the console with `curl`.

Not long thinking, I wrote a little service: https://ip.chuhlomin.com/  
It uses:

* [GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) databases for city and provider (ASN) lookups
* whois.iana.org for Whois queries, for example, https://ip.chuhlomin.com/1.1.1.1/whois

Tested only on IPv4, so PRs are welcome: https://github.com/chuhlomin/ip

#project #ops #go
