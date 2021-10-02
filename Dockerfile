FROM nginx:alpine

ADD nginx.conf /etc/nginx/nginx.conf

ADD output/ /usr/share/nginx/html/

ADD static/ /usr/share/nginx/html/
