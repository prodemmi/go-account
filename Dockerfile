FROM golang:1.22.2

ENV HTTPS_PROXY http://fodev.org:8118

# Create and change to the app directory.
WORKDIR /app

# Copy local code to the container image.
COPY ./server ./

EXPOSE 9090

CMD ["./server"]