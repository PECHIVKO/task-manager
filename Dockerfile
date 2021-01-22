FROM golang:1.15.7-alpine3.13 as builder
COPY go.mod go.sum /go/src/github.com/PECHIVKO/task-manager/
WORKDIR /go/src/github.com/PECHIVKO/task-manager
RUN go mod download
COPY . /go/src/github.com/PECHIVKO/task-manager
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/task-manager github.com/PECHIVKO/task-manager/cmd/api

FROM alpine:latest
RUN apk add --no-cache ca-certificates && update-ca-certificates

RUN mkdir /app/
COPY --from=builder /go/src/github.com/PECHIVKO/task-manager/build/task-manager /app/task-manager
COPY --from=builder /go/src/github.com/PECHIVKO/task-manager/db/migrations /app/migrations
COPY ./cmd/api/config_docker.yaml /app/config.yaml

WORKDIR /app/
RUN chmod +x /app/task-manager
EXPOSE 8181 8181
ENTRYPOINT ["/app/task-manager"]