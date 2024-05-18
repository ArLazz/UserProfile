FROM golang:1.21.1
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /userprofile-server cmd/userprofile-server/main.go

CMD ["/userprofile-server", "--port=8080"]