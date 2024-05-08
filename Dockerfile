FROM golang:1.22.2-alpine3.19

ENV HTTPS_PROXY http://fodev.org:8118

# Create and change to the server directory.
WORKDIR /app

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server

# Expose port
EXPOSE 9090

CMD ["./server"]
