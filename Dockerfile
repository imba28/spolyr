# api build
FROM golang:1.17 as builder

ARG BUILD_NUMBER=dev

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN sed -i "s/dev-build/${BUILD_NUMBER}/" internal/template/files/includes/footer.html && \
    make build-linux

# frontend build
FROM node:16-alpine as frontend_builder

WORKDIR /build

COPY package.json .
COPY package-lock.json .
RUN npm ci

COPY . .
RUN npm run lint && npm run build

# runtime
FROM alpine:3

LABEL maintainer="lukas.gruber1@gmail.com"

RUN apk --no-cache add ca-certificates

RUN addgroup -S spolyr && adduser --system spolyr
USER spolyr

WORKDIR /app
COPY --from=builder --chown=spolyr:spolyr /build/spolyr .
COPY --from=builder --chown=spolyr:spolyr /build/public public
COPY --from=frontend_builder --chown=spolyr:spolyr /build/public/dist public/dist

CMD ["/app/spolyr"]