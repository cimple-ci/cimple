FROM alpine:3.4

ARG CIMPLE_VERSION

WORKDIR /opt/workspace
ENTRYPOINT ["cimple"]
VOLUME ["/opt/workspace"]

ENV PATH /opt/cimple/bin:$PATH

RUN apk add --no-cache git

COPY ./output/cimple-alpine /opt/cimple/bin/cimple
COPY ./server/frontend/templates /opt/frontend/templates
COPY ./server/frontend/assets /opt/frontend/assets

RUN chmod +x /opt/cimple/bin/cimple

