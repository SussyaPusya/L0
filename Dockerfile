FROM golang:1.25.0-alpine3.21 as builder



WORKDIR /appl

COPY . .





RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app_service cmd/app/main.go


FROM alpine:3.21


WORKDIR /appl

COPY --from=builder /appl/app_service .

# Ensure /app/migrations exists in your build context before copying, or remove this line if not needed
# COPY --from=builder /app/migrations ./migrations

COPY --from=builder /appl/db  ./db 
COPY --from=builder /appl/static  ./static


CMD [ "./app_service" ]