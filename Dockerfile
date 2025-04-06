FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
ADD /pkg /app/pkg
RUN go mod download && go mod verify

WORKDIR /app/pkg
RUN go build -o /bin/tokenring

FROM golang:1.23-alpine
COPY --from=builder /bin/tokenring /bin/tokenring

ENTRYPOINT ["/bin/tokenring"]
