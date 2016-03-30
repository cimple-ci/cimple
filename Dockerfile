FROM alpine

COPY ./output/downloads/snapshot/cimple_linux_amd64.tar.gz /opt/cimple_linux_amd64.tar.gz
COPY ./frontend/templates /opt/frontend/templates
COPY ./frontend/assets /opt/frontend/assets

WORKDIR /opt

RUN tar xvzf /opt/cimple_linux_amd64.tar.gz -C /opt
RUN ln -s /opt/cimple_linux_amd64/cimple /usr/local/bin/cimple

ENTRYPOINT ["cimple"]
