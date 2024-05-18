FROM golang:latest

COPY ./ ./
RUN go mod download
RUN go build -o /userprofile-server cmd/userprofile-server/main.go

CMD ["/userprofile-server", "--host=0.0.0.0", "--port=8080"]
