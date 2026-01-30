FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM alpine:latest

ENV ENVIRONMENT=prod

WORKDIR /app

COPY --from=builder /app/main .
COPY config/ /app/config/

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 2000

CMD /app/main -config /app/config/$ENVIRONMENT.yaml
