FROM frolvlad/alpine-glibc:latest

WORKDIR /app
ADD ${source:-app} /app
RUN apk update && apk add ca-certificates
ENTRYPOINT ./main
