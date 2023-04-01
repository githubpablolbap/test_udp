FROM golang:1.19

RUN     mkdir /app
WORKDIR /app
COPY . /app
RUN go mod download
# EXPOSE 55555/udp
# EXPOSE 55555/tcp
RUN go build -o main
EXPOSE 8080
CMD ["/app/main"]

