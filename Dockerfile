FROM --platform=linux/amd64 golang:1.22.0-alpine3.19 as builder
# Installing libvips-dev
# RUN apk update && apk add --no-cache build-base gcc autoconf automake zlib-dev libpng-dev nasm bash vips-dev
RUN apk update && apk add --no-cache gcc libc-dev make pkgconfig vips-dev

WORKDIR /go/src/app

COPY . .

RUN go mod download
RUN GOOS=linux go build -o ./build/sizewise ./cmd/server

FROM alpine:3.19
RUN apk update && apk add --no-cache vips-dev jpeg-dev libpng-dev libwebp-dev tiff-dev giflib-dev

COPY --from=builder /go/src/app/build/sizewise /sizewise

ENTRYPOINT ["/sizewise"]
CMD ["--serve","--config=/sizewiserc.json"]

