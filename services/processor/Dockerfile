FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
COPY . .            

RUN go mod tidy
RUN go build -o /processor

EXPOSE 2112

CMD ["/processor"]
