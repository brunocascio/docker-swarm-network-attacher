FROM golang:1.16-alpine AS base
RUN apk --no-cache add curl
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ADD https://github.com/cosmtrek/air/releases/download/v1.27.3/air_1.27.3_linux_arm64 /usr/bin/air
RUN chmod +x /usr/bin/air
ENTRYPOINT ["air"]

FROM base as build
COPY . .
RUN go build -o /app/dsna

FROM alpine
COPY --from=build /app/dsna /dsna
ENTRYPOINT ["/dsna"]