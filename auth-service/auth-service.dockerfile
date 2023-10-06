# base image for golang
FROM golang:1.18-alpine as builder 

RUN mkdir /app

COPY . /app

WORKDIR /app 

RUN CGO_ENABLED=0 go build -o authApp ./cmd/api 

RUN chmod +x /app/authApp

# copy from the golang image to a base image without golang init only the comlied is copied
FROM alpine:latest 

RUN mkdir /app 

COPY --from=builder /app/authApp /app 

CMD ["/app/authApp"]

