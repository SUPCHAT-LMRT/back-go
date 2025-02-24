FROM golang:1.24-alpine3.21 AS build

ARG TOKEN

WORKDIR /app-data

RUN apk add git gcc build-base

COPY go.mod go.sum ./

# download dependencies and cache them using buildkit
RUN --mount=type=cache,target=/root/go/pkg/mod \
    go mod download -x

COPY . .
RUN mkdir /dist

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/root/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags='-s -w' -o /dist/app cmd/main.go

RUN chmod +x /dist/app

FROM alpine:3.21
WORKDIR /app-data
ENV TZ=Europe/Paris

COPY --from=build /dist /

ENTRYPOINT ["/app"]
