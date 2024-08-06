# Build aşaması
FROM golang:1.22-alpine AS build

WORKDIR /app
COPY . .

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/sasa-elterminali-service ./cmd/api

# Çalışma aşaması
FROM scratch

EXPOSE 45013
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/cmd/api/config/config.json /config/config.json
COPY --from=build /bin/sasa-elterminali-service /bin/sasa-elterminali-service
CMD ["/bin/sasa-elterminali-service"]
