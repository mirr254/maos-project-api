FROM golang:1.21.5 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl bash && \
    addgroup -S app && \
    adduser -S app -G app

USER app

RUN curl -fsSL https://get.pulumi.com | sh -s -- --version 3.115.2
ENV PATH=$PATH:/home/app/.pulumi/bin

WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app .

EXPOSE 8080

CMD ["./app"]
