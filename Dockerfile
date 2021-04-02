FROM golang:1.16 as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN make build-linux

# runtime
FROM alpine:3

LABEL maintainer="lukas.gruber1@gmail.com"

RUN apk --no-cache add ca-certificates

RUN addgroup -S spolyr && adduser --system spolyr
USER spolyr

WORKDIR /app
COPY --from=builder --chown=spolyr:spolyr /build/spolyr .
COPY --from=builder --chown=spolyr:spolyr /build/public public

CMD ["/app/spolyr"]