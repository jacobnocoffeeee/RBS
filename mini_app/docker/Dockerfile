FROM golang:1.24.4-alpine3.22 AS builder
RUN addgroup --g 1000 nonroot \
&& adduser -u 1000 -G nonroot -s /bin/sh -D nonroot
USER nonroot
WORKDIR /app
COPY ./src .
RUN go build -o main

FROM alpine:3.22
RUN addgroup --g 1000 nonroot \
&& adduser -u 1000 -G nonroot -s /bin/sh -D nonroot
USER nonroot
COPY --from=builder app/main .
ENTRYPOINT ["./main"]