FROM golang:1.12

MAINTAINER xiliangMa "xiliangMa@outlook.com"

EXPOSE 8080
EXPOSE 8088
EXPOSE 10443


WORKDIR /usr/share/diss-backend
RUN mkdir -p /usr/share/diss-backend/conf && \
    mkdir -p /usr/share/diss-backend/swagger


COPY entrypoint.sh /usr/share/diss-backend/
COPY bin/diss-backend /usr/share/diss-backend/
COPY conf /usr/share/diss-backend/conf
COPY swagger /usr/share/diss-backend/swagger

RUN chmod +x /usr/share/diss-backend/entrypoint.sh

# RUN ln -s /usr/share/diss-backend/diss-backend /usr/local/bin/diss-backend
ENTRYPOINT ["sh", "/usr/share/diss-backend/entrypoint.sh"]
