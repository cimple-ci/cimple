FROM progrium/busybox

ARG CIMPLE_VERSION

WORKDIR /opt/workspace
ENTRYPOINT ["cimple"]
VOLUME ["/opt/workspace"]

ENV PATH /opt/cimple/bin:$PATH

RUN opkg-install git

COPY ./output/downloads/$CIMPLE_VERSION/cimple_${CIMPLE_VERSION}_linux_amd64.tar.gz /opt/cimple_linux_amd64.tar.gz
COPY ./server/frontend/templates /opt/frontend/templates
COPY ./server/frontend/assets /opt/frontend/assets

RUN cd /tmp \
    && zcat /opt/cimple_linux_amd64.tar.gz | tar -xvf - \
    && chmod +x /tmp/cimple_${CIMPLE_VERSION}_linux_amd64/cimple \
    && mkdir -p /opt/cimple/bin \
    && mv /tmp/cimple_${CIMPLE_VERSION}_linux_amd64/cimple /opt/cimple/bin \
    && rm -rf /tmp/cimple_linux_amd64 \
    && rm /opt/cimple_linux_amd64.tar.gz

