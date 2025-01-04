FROM golang:1.23-alpine AS base
WORKDIR /app 
COPY ./go.mod ./go.sum ./
RUN go mod download

FROM base AS builder 
ENV CGO_ENABLED=0\
    GOOS=linux\
    GOARCH=amd64
WORKDIR /app 
COPY . .
RUN go build -o server . 

FROM gcr.io/distroless/cc AS runtime 
WORKDIR /root
COPY --from=builder /app/server server
ENTRYPOINT [ "/root/server" ]
