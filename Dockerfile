FROM alpine:latest

RUN mkdir injector
COPY envsecrets injector/envsecrets
RUN chmod +x injector/envsecrets