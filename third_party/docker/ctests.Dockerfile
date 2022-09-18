FROM golang:1.19

# Create the working directory.
WORKDIR /app

# Copy the source code.
COPY . .

RUN go mod download

# Run the tests.
CMD ["go", "test", "-v", "./..."]
