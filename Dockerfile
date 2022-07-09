FROM nginx:alpine

ADD nginx.conf /etc/nginx/nginx.conf

ADD output/ /usr/share/nginx/html/

ADD _static/ /usr/share/nginx/html/
