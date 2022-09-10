FROM golang:1.19 as build
WORKDIR /go/src/app

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static

COPY --from=build /go/bin/app /
COPY --from=build /go/src/app/static/. /static/.
CMD ["/app"]