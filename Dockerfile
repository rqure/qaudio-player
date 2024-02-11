# Build the application from source
FROM golang:alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

COPY audio/*.mp3 ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /qapp

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

RUN apk update && apk add --no-cache alsa-lib pulseaudio alsa-utils alsa-plugins-pulse mpg123

WORKDIR /

COPY --from=build-stage /qapp /qapp
COPY --from=build-stage /app/*.mp3 /

ENV PULSE_SERVER=172.17.0.1
ENV PULSE_COOKIE=/run/pulse/cookie

ENTRYPOINT ["/qapp"]
