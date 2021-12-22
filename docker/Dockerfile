FROM golang:1.17-alpine AS builder
RUN mkdir -p /build
COPY ./ /build/

WORKDIR "/build"
RUN go mod tidy
RUN go build -o bin/ ./...

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/bin/ ./

RUN mkdir -p certs
#RUN mkdir -p content
COPY --from=builder /build/content/ ./content
#RUN mkdir -p templates
COPY --from=builder /build/templates/*.html ./templates/


ENTRYPOINT ./easytls
