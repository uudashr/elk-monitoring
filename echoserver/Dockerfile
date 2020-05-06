# == Builder ==
FROM golang:1.14.2-alpine3.11 as builder

RUN apk add --no-cache bash=5.0.11-r1 git=2.24.3-r0

WORKDIR /app

# Copy dependencies definition
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy remaining source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /app/echoserver *.go

# == Runner ==
FROM alpine:3.11.6

# Copy binary from builder
COPY --from=builder /app/echoserver /app/echoserver

CMD ["/app/echoserver"] 