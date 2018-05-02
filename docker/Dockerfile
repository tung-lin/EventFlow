FROM frolvlad/alpine-glibc:latest

ARG APP
WORKDIR /app
ADD ${APP} /app
RUN apk update && apk add ca-certificates
ENTRYPOINT ./main
