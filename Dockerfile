FROM golang:1.14-alpine AS build

RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go get ./...
RUN GOOS=linux go build -o ./bin/activitygenerator

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
ENV BROKER_ADDR amqp://guest:guest@localhost:5672/
ENTRYPOINT /go/bin/activitygenerator