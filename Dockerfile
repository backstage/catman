FROM golang:buster as app
RUN mkdir /app
WORKDIR /app
COPY cmd ./cmd
COPY go.mod .
COPY go.sum .
RUN go build ./cmd/catman

FROM gcr.io/distroless/base
COPY --from=app /app/catman /
ENTRYPOINT ["/catman"]