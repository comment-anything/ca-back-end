FROM golang:1.20.3-alpine3.16 AS build
ARG port

WORKDIR /CANY
COPY . .
RUN go get .
RUN go build -o ca-back .
CMD ["bin/sh"]

FROM alpine:3.17.3

EXPOSE $port
COPY --from=build /CANY/ca-back /CANY/serve
COPY --from=build /CANY/.env /CANY/.env


CMD ["/CANY/serve", "--env=/CANY/.env", "--docker=true"]