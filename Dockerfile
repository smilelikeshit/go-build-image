FROM alpine:latest

WORKDIR /app

RUN echo "hello" > bandung.txt

RUN echo "bandung" > jakarta.txt

RUN echo "semarang" > semarang.txt