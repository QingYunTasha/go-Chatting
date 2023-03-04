FROM golang:1.20.1-bullseye
WORKDIR /builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main.exe cmd/web/main.go


FROM alpine:latest
WORKDIR /app
COPY --from=0 /builder/main.exe ./
CMD ["./main.exe"]