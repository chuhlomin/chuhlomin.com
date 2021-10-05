---
created: 2020-06-28
---

# .well-known

В 2010м вышел RFC «Defining Well-Known Uniform Resource Identifiers (URIs)» который предложил префикс «/.well-known/» для так называемых «известных мест» на ресурсе (сайте).<br>
https://tools.ietf.org/html/rfc5785

Например, если пройти по ссылке https://github.com/.well-known/change-password, то вы будете переадресованы на страницу изменения пароля. Это использует, в частности, Сафари: если обнаружит, что ваш пароль слишком слабый, то предложит его поменять.<br>
https://wicg.github.io/change-password-url/

Wikipedia приводит список наиболее распространенных «.well-known», которые могут предоставлять сервисы:<br>
https://en.wikipedia.org/wiki/List_of_%2F.well-known%2F_services_offered_by_webservers

Чтобы избежать коллизий, новые предлагаемые сервисы должны быть зарегистрированы в IANA:<br>
https://www.iana.org/assignments/well-known-uris/well-known-uris.xhtml

Особенно рекомендую посмотреть на security.txt<br>
https://securitytxt.org/ 

#rfc #ops
