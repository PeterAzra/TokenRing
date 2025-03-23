FROM golang:1.23

WORKDIR /home/dev/Projects/TokenRing
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /home/dev/Projects/TokenRing/pkg ./...

CMD ["/home/dev/Projects/TokenRing/pkg/pkg"]
