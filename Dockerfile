FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/glink ./cmd/app

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /bin/glink /app/glink

EXPOSE 8080

ENTRYPOINT ["/app/glink"]
