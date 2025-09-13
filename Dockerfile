# ビルドステージ
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Goモジュールのダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピーとビルド
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./main.go

# 実行ステージ
FROM alpine:3.17

WORKDIR /app

# タイムゾーンの設定
RUN apk add --no-cache tzdata

# ビルドステージから実行ファイルをコピー
COPY --from=builder /app/bin/server .

# 環境変数の設定
ENV MYSQL_HOST=mysql
ENV MYSQL_PORT=3306
ENV MYSQL_USER=root
ENV MYSQL_PASSWORD=password
ENV MYSQL_DATABASE=vital_chi_db

# ポートの公開
EXPOSE 8080

# コンテナ起動時のコマンド
CMD ["./server"]