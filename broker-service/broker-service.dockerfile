# base image for golang
FROM golang:1.18-alpine as builder 

RUN mkdir /app

COPY . /app

WORKDIR /app 

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api 

RUN chmod +x /app/brokerApp

# copy from the golang image to a base image without golang init only the comlied is copied
FROM alpine:latest 

RUN mkdir /app 

COPY --from=builder /app/brokerApp /app 

CMD ["/app/brokerApp"]

