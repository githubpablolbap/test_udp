FROM golang:1.19

RUN     mkdir /app
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o main
EXPOSE 8080
# EXPOSE 22222/udp
# EXPOSE 22222/tcp
CMD ["/app/main"]

