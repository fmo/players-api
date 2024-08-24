FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make players_api_amd

CMD ["./playersApiApp"]
