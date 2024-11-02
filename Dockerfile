FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go get github.com/gin-gonic/gin
RUN go install github.com/jessevdk/go-assets-builder@latest
RUN make build

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 3030
ENTRYPOINT /go/bin/sm
