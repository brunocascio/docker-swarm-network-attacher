FROM golang:1.16-buster AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/dsna

FROM gcr.io/distroless/base-debian10
COPY --from=build /app/dsna /dsna
USER nonroot:nonroot
ENTRYPOINT ["/dsna"]