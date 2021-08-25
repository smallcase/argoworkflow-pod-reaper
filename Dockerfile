# builder image
FROM golang:1.16.7-alpine3.14 as builder

RUN mkdir /build
ADD *.go /build/
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o reaper .


# generate clean, final image for end users
FROM alpine:3.14
COPY --from=builder /build/reaper .

# executable
ENTRYPOINT [ "./reaper" ]