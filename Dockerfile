FROM golang:1.20-alpine

WORKDIR /app

RUN apk update && apk add alpine-sdk && apk add bash
RUN apk add --no-cache tzdata && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && apk del tzdata
ENV TZ=Asia/Tokyo

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o build/main cmd/app/main.go

EXPOSE 8080
CMD ["./build/main"]