FROM --platform=$BUILDPLATFORM docker.io/golang:1.21 as server-builder
ARG TARGETPLATFORM
WORKDIR /usr/src/app

COPY . .
RUN GOOS=linux GOARCH=$(echo $TARGETPLATFORM | sed 's/linux\///') \
  go build -o dist/smg src/main.go

FROM docker.io/debian:stable-slim as runner
WORKDIR /app
COPY --from=server-builder /usr/src/app/dist/smg /app
COPY templates /app/templates
COPY static /app/static

EXPOSE 3333
CMD ["./smg"]