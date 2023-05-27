FROM golang:1.20.3 AS builder
run mkdir /distributing
WORKDIR /distributing
COPY notify/ notify/
COPY pomo/ pomo/
WORKDIR /distributing/pomo
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -tags=containers

FROM alpine:latest
RUN mkdir /app && adduser -h /app -D pomo
WORKDIR /app
COPY --chown=pomo --from=builder /distributing/pomo/pomo .
CMD ["/app/pomo"]