FROM golang:1.19-alpine

# Create the working directory.
WORKDIR /app

# Install gofumpt
RUN go install mvdan.cc/gofumpt@latest

# Run gofumpt
CMD ["gofumpt", "-l", "-w", "."]
