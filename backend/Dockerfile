FROM golang:1.21-alpine AS builder

WORKDIR /app

ARG GIT_SHA
ENV GIT_SHA=${GIT_SHA}

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
	-ldflags "-X main.Version=${GIT_SHA}" \
	-o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .

RUN mkdir -p /root/uploads

EXPOSE 8080

CMD ["./main"]
