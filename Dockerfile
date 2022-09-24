FROM golang:1.19-alpine as builder

LABEL maintainer="xyedo | Hafid Mahdi"

RUN apk update && apk add --no-caache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o webapi -race ./cmd/web

FROM golang:1.19-alpine


RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/webapi .

EXPOSE 4000

CMD ["./webapi"]