FROM golang:1.22

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o cars-catalog ./cmd/cars-catalog/main.go

CMD ["./cars-catalog"]