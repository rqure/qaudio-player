# Build the application from source
FROM golang:1.21.6 AS build-stage

RUN apt-get update && apt-get install -y libasound2-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

COPY audio/*.mp3 ./

RUN CGO_ENABLED=1 GOOS=linux go build -o /qapp

# Deploy the application binary into a lean image
FROM debian:bookworm-slim AS build-release-stage

RUN apt-get update && apt-get install -y libasound2-dev pulseaudio alsa-utils libasound2-plugins

WORKDIR /

COPY --from=build-stage /qapp /qapp
COPY --from=build-stage /app/*.mp3 /

ENV PULSE_SERVER=172.17.0.1
ENV PULSE_COOKIE=/run/pulse/cookie

ENTRYPOINT ["/qapp"]
