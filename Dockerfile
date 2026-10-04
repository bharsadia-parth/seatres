FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server
RUN CGO_ENABLED=0 go build -o /out/burst ./cmd/burst

FROM alpine:3.20
RUN adduser -D app
USER app
COPY --from=build /out/server /server
COPY --from=build /out/burst /burst
EXPOSE 8080
CMD ["/server"]