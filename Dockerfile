FROM golang:1.21

WORKDIR /

COPY ./ ./

RUN go mod download

RUN cd cmd && go build -o /source
RUN cd ../

CMD ["/source"]

