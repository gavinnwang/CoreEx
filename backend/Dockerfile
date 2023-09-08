FROM golang:latest

RUN go install github.com/cosmtrek/air@latest

RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Intall go migrate tool
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app/src

COPY go.* ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/main ./cmd/main.go

EXPOSE 8080

CMD [ "./start.sh" ] 