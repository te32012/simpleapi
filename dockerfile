FROM golang:1.22 as build
COPY . .
RUN go build -o /server ./cmd/main.go



FROM ubuntu:22.04 as deploy
RUN apt-get update && apt-get -y upgrade

WORKDIR /app
COPY --from=build /server .

CMD ["/app/server"]