FROM node:16-alpine3.14 as build-web

WORKDIR /src
COPY . .
RUN cd web && npm install

ARG NODE_ENV=production
RUN cd web && npm run build

FROM golang:1.16-alpine3.14 as build

RUN apk add \
	make
WORKDIR /src
COPY --from=build-web /src .

ARG CGO_ENABLED=0
RUN make TAGS=production

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /src/upl /upl

CMD ["/upl"]
