# syntax=docker/dockerfile:1
FROM golang:1.25 AS builder

ARG VERSION
ARG BUILD
ARG CI_JOB_TOKEN
ARG BUILD_DATE
ENV GOPRIVATE=github.com/jedi-tasks/go-server-grpc
RUN apt install ca-certificates
RUN git config --global url.https://${CI_JOB_TOKEN}@github.com/.insteadOf https://github.com/
ADD . /go-server
WORKDIR /go-server/server
RUN sed -i -e "s|@{build}|${BUILD}|g" -e "s|@{version}|${VERSION}|g" -e "s|@{buildDate}|${BUILD_DATE}|g" internal/main.go
RUN go build -o ../main internal/main.go

FROM ubuntu:latest
WORKDIR /go-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go-server/main .
EXPOSE 8080
EXPOSE 8090
USER root:root
CMD ["./main"]
