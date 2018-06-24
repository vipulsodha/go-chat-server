FROM golang:latest
RUN mkdir app/
ADD . /app/
WORKDIR /app
RUN go get golang.org/x/net/websocket
RUN go build chat-server.go
CMD ["/app/chat-server"]