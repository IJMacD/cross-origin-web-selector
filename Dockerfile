FROM golang:1.23.1 AS build
WORKDIR /usr/src/app
COPY src/go.mod src/go.sum ./
RUN go mod download -x
COPY src ./
RUN CGO_ENABLED=0 go build -v -o /app main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /app
ENTRYPOINT [ "/app" ]