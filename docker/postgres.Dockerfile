FROM ubuntu as intermediate

LABEL stage=intermediate

WORKDIR /root/temp

RUN apt-get update
RUN apt-get install -y git
RUN echo "meep"
RUN echo "moop"
RUN git clone https://github.com/alpstable/gpostgres.git

FROM golang:1.19

# Create the working directory.
WORKDIR /app

COPY --from=intermediate /root/temp/gpostgres .

RUN go mod tidy

# Run the tests.
CMD ["make", "tests"]
