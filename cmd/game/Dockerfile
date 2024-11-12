FROM golang:1.22.7 as builder

ENV GOOS=linux
ENV GOARCH=amd64
COPY . .
RUN go build -o /game ./cmd/game

FROM ubuntu:focal

RUN apt update
RUN apt install -y ca-certificates
COPY --from=builder /game /

CMD ["/game"]
