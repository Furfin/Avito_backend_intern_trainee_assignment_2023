FROM golang:1.21

WORKDIR /

COPY ./ ./

RUN go mod download

RUN cd cmd && go build -o /source
RUN cd ../

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose

# Run
CMD ["/source"]

