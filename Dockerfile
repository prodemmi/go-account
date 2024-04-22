FROM golang:1.22.2

ENV HTTPS_PROXY http://fodev.org:8118

# Create and change to the server directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server

CMD ["./server"]
