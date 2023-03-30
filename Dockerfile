FROM golang:1.19

RUN     mkdir /app
WORKDIR /app
ADD go.mod main.go /app/
EXPOSE 55555/udp
EXPOSE 55555/tcp
# ADD     go.mod main.go /app/
RUN     go build main

CMD     ./main

