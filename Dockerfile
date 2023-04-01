FROM golang:1.19

RUN     mkdir /app
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o main
# EXPOSE 8080
EXPOSE 55555/udp
EXPOSE 55555/tcp
CMD ["/app/main"]

