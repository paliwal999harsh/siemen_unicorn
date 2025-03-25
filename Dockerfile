FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o unicorn_app ./cmd

FROM alpine:3.21.3

RUN addgroup --system unicorn && adduser --system --ingroup unicorn unicorn

WORKDIR /home/unicorn

COPY --from=builder /app/unicorn_app .
COPY --from=builder /app/res ./res

RUN chown -R unicorn:unicorn /home/unicorn/

USER unicorn

EXPOSE 8888

CMD ["./unicorn_app"]