FROM golang:1.20

RUN go version
ENV GOPATH=/
LABEL authors="adiletkemelkhan"


COPY ./ ./

RUN go mod download
RUN go build -o testproject ./main.go

ENTRYPOINT ["./testproject"]