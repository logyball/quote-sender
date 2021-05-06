FROM golang:alpine3.13 as builder
WORKDIR /app
ENV GOARCH=arm

RUN mkdir /dist

COPY . /app/.
RUN go mod download
RUN go test -v ./...
RUN go build -o /dist/quote

# Deployment container
FROM arm64v8/alpine:3.13
WORKDIR /dist
COPY --from=builder /dist/quote /dist/quote
