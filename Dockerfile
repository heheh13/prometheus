FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prometheus-demo .


FROM alpine:latest

WORKDIR /app
## defining flags from to copy form the previous os 
RUN apk add curl
COPY --from=builder /app/prometheus-demo .
COPY .env .env

ENTRYPOINT ["./prometheus-demo"]