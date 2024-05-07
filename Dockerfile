# ビルドステージ
FROM golang:1.20-alpine as builder

WORKDIR /app

# 必要なパッケージのインストール
RUN apk update && \
    apk add --no-cache alpine-sdk bash tzdata

# タイムゾーン設定
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata

# 依存関係の管理
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピーとビルド
COPY . .
RUN go build -o /main cmd/app/main.go

# 実行ステージ
FROM alpine:latest

WORKDIR /app

# 必要なパッケージのみインストール
RUN apk add --no-cache ca-certificates bash

# ビルドステージから実行ファイルをコピー
COPY --from=builder /main ./

# ポート番号の公開
EXPOSE 8080

# 実行コマンド
CMD ["./main"]