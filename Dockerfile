FROM golang:1.13 as builder

LABEL maintainer="Niklas Engberg <per.niklas.engberg@gmail.com>"

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

ENV SPOTIFY_CLIENT_ID="abc"
ENV SPOTIFY_CLIENT_SECRET="def"
ENV SPOTIFY_AGGREGATION_LIST_ID="ghi"

EXPOSE 8080

# Command to run the executable
CMD ["./main"]