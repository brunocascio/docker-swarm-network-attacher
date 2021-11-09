# syntax = docker/dockerfile:1.3
FROM golang:1.16-alpine AS base
RUN apk --no-cache add curl
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=target=/root/.cache,type=cache go mod download
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s \
  && mv ./bin/air /usr/bin/air
ENTRYPOINT ["air", "-c", ".air.toml"]

FROM base as build
COPY . .
RUN go build -o /app/dsna

FROM alpine
COPY --from=build /app/dsna /dsna
ENTRYPOINT ["/dsna"]