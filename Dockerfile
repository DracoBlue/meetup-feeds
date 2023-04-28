FROM golang:1.19.1-alpine3.16 AS build
WORKDIR /app/
ADD go.mod go.sum /app/
RUN go mod download
ADD *.go /app/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /meetup-feeds

FROM scratch
COPY --from=build /meetup-feeds /meetup-feeds
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/meetup-feeds"]
