FROM node:15.7 as build-vue
WORKDIR /usr/src/app
COPY web/vue.js/ ./
RUN yarn install && yarn build

FROM golang:1.15 as build-app
WORKDIR /go/src/vueapp
COPY ./ .
COPY --from=build-vue /usr/src/app/dist web/vue.js/dist
RUN go get -u github.com/mjibson/esc
RUN make build

FROM alpine
WORKDIR /
COPY --from=build-app /go/src/vueapp/bin/starter ./
COPY config.yml ./
EXPOSE 8080
CMD ["./starter"]
