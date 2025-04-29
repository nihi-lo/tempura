FROM golang:1.24

ENV TZ=Asia/Tokyo

WORKDIR /app

# Goの依存パッケージをキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# Taskをインストール(https://taskfile.dev/installation/#install-script)
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
