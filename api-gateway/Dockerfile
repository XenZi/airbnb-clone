FROM nginx:latest

COPY ./api-gateway.conf /etc/nginx/conf.d/api-gateway.conf
COPY ./localhost.crt /etc/ssl/certs/localhost.crt
COPY ./localhost.csr /etc/ssl/csr/localhost.csr
COPY ./localhost.key /etc/ssl/private/localhost.key
COPY ./localhost.key.org /etc/ssl/keyorg/localhost.key.org
COPY ./localhost.ext /etc/ssl/v3/localhost.ext
RUN rm /etc/nginx/conf.d/default.conf

RUN apt-get update ; apt-get install unzip
ADD https://github.com/open-telemetry/opentelemetry-cpp-contrib/releases/download/webserver%2Fv1.0.0/opentelemetry-webserver-sdk-x64-linux.tgz.zip /opt
RUN cd /opt ; unzip opentelemetry-webserver-sdk-x64-linux.tgz.zip; tar xvfz opentelemetry-webserver-sdk-x64-linux.tgz
RUN cd /opt/opentelemetry-webserver-sdk; ./install.sh
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/opentelemetry-webserver-sdk/sdk_lib/lib
RUN echo "load_module /opt/opentelemetry-webserver-sdk/WebServerModule/Nginx/ngx_http_opentelemetry_module.so;\n$(cat /etc/nginx/nginx.conf)" > /etc/nginx/nginx.conf
COPY opentelemetry_module.conf /etc/nginx/conf.d