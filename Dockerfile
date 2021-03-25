FROM golang:1.15 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app .

FROM alpine
COPY --from=builder /bin/app /app

ENV ADDRESS=default
ENV Port=8000
ENV Schema=http

ENTRYPOINT ["/app"]
