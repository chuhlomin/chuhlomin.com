---
date: 2020-06-28
---

# .well-known

In 2010, the RFC "Defining Well-Known Uniform Resource Identifiers (URIs)" came out and suggested the prefix "/.well-known/" for so-called "known places" on the resource (site).  
https://tools.ietf.org/html/rfc5785

For example, if you follow the link https://github.com/.well-known/change-password, you will be redirected to a password change page. Safari in particular uses this: if it detects that your password is too weak, it will offer to change it.  
https://wicg.github.io/change-password-url/

Wikipedia lists the most common ".well-known" that can be provided by services:  
https://en.wikipedia.org/wiki/List_of_%2F.well-known%2F_services_offered_by_webservers

To avoid conflicts, new proposed services must be registered with the IANA:  
https://www.iana.org/assignments/well-known-uris/well-known-uris.xhtml

I especially recommend looking at security.txt  
https://securitytxt.org/

#rfc #ops
