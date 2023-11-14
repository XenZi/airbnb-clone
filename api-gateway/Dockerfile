FROM nginx:latest

COPY ./api-gateway.conf /etc/nginx/conf.d/api-gateway.conf
COPY ./localhost.crt /etc/ssl/certs/localhost.crt
COPY ./localhost.csr /etc/ssl/csr/localhost.csr
COPY ./localhost.key /etc/ssl/private/localhost.key
COPY ./localhost.key.org /etc/ssl/keyorg/localhost.key.org
COPY ./localhost.ext /etc/ssl/v3/localhost.ext
RUN rm /etc/nginx/conf.d/default.conf