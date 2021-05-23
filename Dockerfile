FROM node:16-alpine3.13 as build-web

ARG NODE_ENV=production

WORKDIR /src
COPY . .
RUN cd web && npm install
RUN cd web && npm run build

FROM golang:1.16-alpine3.13 as build

ARG CGO_ENABLED=0

WORKDIR /src
COPY --from=build-web . .
RUN make TAGS=production

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build upl /upl

RUN ["/upl"]
