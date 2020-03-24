FROM golang:1.12 as builder
MAINTAINER xiliangMa "xiliangMa@outlook.com"
EXPOSE 8080
EXPOSE 8088
EXPOSE 10443
WORKDIR /build
COPY entrypoint.sh .
COPY bin/diss-backend .
COPY conf ./conf
COPY upload ./upload
COPY swagger ./swagger

FROM alpine:3.11
WORKDIR /opt/diss-backend
COPY --from=builder /build/ .
RUN chmod +x ./entrypoint.sh
VOLUME /opt/diss-backend
ENTRYPOINT ["sh", "./entrypoint.sh"]
