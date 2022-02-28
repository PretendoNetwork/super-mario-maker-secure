# syntax=docker/dockerfile:1
FROM golang:1.17.0-bullseye as builder
WORKDIR /app
# Set build architecture
ENV GOOS linux
ENV CGO_ENABLED 0
# Install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code
COPY . .
# build app
RUN go build -o build.out

FROM alpine:3.14 as runner
WORKDIR /app
# Add certificates
RUN apk add --no-cache ca-certificates
# Copy the app from the builder stage
COPY --from=builder /app/build.out /app/build.out
# Expose port
EXPOSE 60003/udp
CMD /app/build.out
