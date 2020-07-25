FROM golang:alpine
RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY resources .
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
RUN mkdir -p ./storage
RUN chown appuser ./storage
RUN mkdir -p ./uploads
RUN chown appuser ./uploads
USER appuser
EXPOSE 1421
CMD ["./main"]