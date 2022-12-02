FROM golang:alpine3.17 as builder
WORKDIR /app
ENV GOARCH=arm

RUN mkdir /dist

COPY . /app/.
RUN go mod download
RUN go build -o /dist/quoteCats

# Deployment container
FROM arm64v8/alpine:3.17
WORKDIR /dist
COPY --from=builder /dist/quoteCats /dist/quoteCats
