FROM golang:1.24.4

WORKDIR /back

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd

CMD [ "./main" ]