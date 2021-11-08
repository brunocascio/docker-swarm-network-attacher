FROM golang:1.16-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/dsna

FROM alpine
COPY --from=build /app/dsna /dsna
ENTRYPOINT ["/dsna"]