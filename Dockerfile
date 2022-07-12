FROM openapitools/openapi-generator-cli as openapi-generator

WORKDIR /local
COPY oapi-spec.yaml .

RUN /usr/local/bin/docker-entrypoint.sh generate -g go-server -i /local/oapi-spec.yaml \
    -o /local/pkg/openapi \
    --additional-properties=outputAsLibrary=true,onlyInterfaces=true,sourceFolder=openapi,addResponseHeaders=true && \
    sed -i -e 's/"github.com\/gorilla\/mux"//g' pkg/openapi/openapi/api_auth.go && \
    sed -i -e 's/"encoding\/json"//g' pkg/openapi/openapi/api_import.go && \
    sed -i -e 's/"encoding\/json"//g' -e 's/"github.com\/gorilla\/mux"//g' pkg/openapi/openapi/api_playlists.go

RUN /usr/local/bin/docker-entrypoint.sh generate -g javascript -i /local/oapi-spec.yaml -o \
    /local/assets/openapi \
    --additional-properties=usePromises=true,moduleName=@/openapi \
    --global-property models,modelTests=false \
    --global-property apis,apiTests=false \
    --global-property supportingFiles

# api build
FROM golang:1.18 as builder

ARG BUILD_NUMBER=dev

WORKDIR /build

COPY --from=openapi-generator /local/pkg/openapi /build/pkg/openapi

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN sed -i "s/dev-build/${BUILD_NUMBER}/" assets/App.vue && \
    make build-linux

# frontend build
FROM node:14-alpine as frontend_builder

WORKDIR /build

COPY --from=openapi-generator /local/assets/openapi/src /build/assets/openapi

COPY package.json .
COPY package-lock.json .
RUN npm remove cypress
RUN npm ci

COPY . .
RUN npm run build

# runtime
FROM alpine:3

LABEL maintainer="lukas.gruber1@gmail.com"

RUN apk --no-cache add ca-certificates

RUN addgroup -S spolyr && adduser --system spolyr
USER spolyr

WORKDIR /app
COPY --from=builder --chown=spolyr:spolyr /build/spolyr .
COPY --from=frontend_builder --chown=spolyr:spolyr /build/public public

ENTRYPOINT ["/app/spolyr"]
CMD ["web"]