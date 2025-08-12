# Dockerfile
FROM tip-alpine3.22

WORKDIR /tempbuild
# Copy App Over App
COPY . /tempbuild
RUN mkdir build
RUN go build -o build/tictactoe main.go

RUN mkdir /srv
RUN mv build/tictactoe /srv/tictactoe


# Launch
CMD ["/srv/tictactoe"]
