FROM golang:1.14 AS build

WORKDIR /go/src/build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./app ./application.go

FROM alpine:latest

WORKDIR /go/src/deploy
COPY --from=build /go/src/build .
RUN chmod +x ./app

# You might need to change this settings according to your configuration
EXPOSE 3000

CMD ["./app"]
