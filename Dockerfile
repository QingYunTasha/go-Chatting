FROM golang:latest
WORKDIR /builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/main.exe cmd/main.go


FROM alpine:latest
WORKDIR /app
COPY --from=0 /builder/cmd/main.exe ./
CMD ["./main.exe"]
