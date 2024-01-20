FROM golang:1.21-alpine AS build-stage
RUN apk --no-cache add g++


WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o syncbot

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=build-stage /app/syncbot .
COPY --from=build-stage /app/config.yaml .

ENTRYPOINT ["./syncbot"]