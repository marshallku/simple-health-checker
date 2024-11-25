FROM golang:1.23.3-alpine3.20 AS build

WORKDIR /app
COPY . .

RUN go build -ldflags "-s -w" -o statusy .

FROM scratch

COPY --from=build /app/statusy /statusy

EXPOSE 8080
CMD ["/statusy", "--mode", "server"]
