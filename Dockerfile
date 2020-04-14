# --- Build stage ---
FROM golang:1.14 AS builder
WORKDIR /app
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
# Copy the dependency definition
COPY go.mod .
COPY go.sum .
# Download dependencies
RUN go mod download
# Copy the source code
COPY . .
# Build for release
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# --- Final stage ---
FROM gcr.azk8s.cn/distroless/static-debian10:latest
COPY --from=builder /app/message-hub-mock /
ENTRYPOINT ["/message-hub-mock"]
EXPOSE 3000
LABEL Name=message-hub-mock \
	Version=0.1
