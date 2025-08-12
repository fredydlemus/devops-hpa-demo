#Build
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY app/go.mod ./app/go.mod
WORKDIR /src/app
RUN go mod download
COPY app/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server .

#Run (distroless)
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/server /app/server
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/server"]