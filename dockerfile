FROM golang:latest
RUN go version
env GOPATH=/
COPY ./ ./
RUN go mod download
RUN go build -o app .
CMD ["./app"]