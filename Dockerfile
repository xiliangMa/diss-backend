FROM golang:1.15 as builder
MAINTAINER xiliangMa "xiliangMa@outlook.com"
EXPOSE 8080
EXPOSE 8088
EXPOSE 10443
WORKDIR /build
COPY entrypoint.sh .
COPY build/bin/diss-backend .
COPY conf ./conf
COPY swagger ./swagger
COPY upload ./upload

FROM alpine:3.14
WORKDIR /opt/diss-backend
COPY --from=builder /build/ .
RUN chmod +x ./entrypoint.sh
RUN apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
VOLUME /opt/diss-backend
ENTRYPOINT ["sh", "./entrypoint.sh"]
