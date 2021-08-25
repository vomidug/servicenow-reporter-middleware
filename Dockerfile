FROM golang:1.16-alpine AS build
RUN apk update; apk add upx
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags "-s -w" -o /main
RUN upx --best --lzma /main

FROM scratch
WORKDIR /
COPY --from=build /main /main
ENTRYPOINT ["/main"]
