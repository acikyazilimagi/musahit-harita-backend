# stage 1
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -o api cmd/api/main.go

# stage 2
FROM gcr.io/distroless/static-debian11 as runner
COPY --from=builder --chown=nonroot:nonroot /app/api /
EXPOSE 80
ENTRYPOINT ["/api"]