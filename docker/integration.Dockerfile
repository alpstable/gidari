FROM golang:1.19

# Create the working directory.
WORKDIR /app

# Copy the source code.
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Run the tests.
ENTRYPOINT ["go", "test", "-v", "./..."]
