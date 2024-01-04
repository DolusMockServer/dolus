FROM golang:1.21.5-alpine AS builder
RUN apk update && apk add --no-cache git
RUN apk add build-base
RUN apk --no-cache add ca-certificates && update-ca-certificates


ENV CGO_ENABLED=0

WORKDIR /app


COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

COPY . .
RUN cd cmd/dolus && go build -o /dolus

FROM scratch

EXPOSE 1080

WORKDIR /app 

COPY --from=builder /dolus /app/dolus
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


ENTRYPOINT ["/app/dolus"]
