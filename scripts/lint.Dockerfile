FROM golangci/golangci-lint:v1.50.0
# Create the working directory.
WORKDIR /app

COPY . .

# Run the tests.
CMD ["golangci-lint", "run", "--config", ".golangci.yml"]