FROM alpine:latest

ARG GF_UID=472
ARG GF_GID=472
ARG EDGE_COMMUNITY_REPO="https://dl-cdn.alpinelinux.org/alpine/edge/community"

ENV GF_PATHS_DATA="/var/lib/grafana"
ENV GF_PATHS_LOGS="/var/log/grafana"
ENV GF_PATHS_PLUGINS="/var/lib/grafana/plugins"
ENV GF_PATHS_PROVISIONING="/etc/grafana/provisioning"
ENV GF_PATHS_CONFIG="/etc/grafana/grafana.ini"
ENV GF_PATHS_HOME="/usr/share/grafana"

ENV DEPLOYENV=DEV
ENV MSNAME=logger-host-1.0

ENTRYPOINT sh -c /entry_point.sh


# 1. Install base dependencies, OpenRC, and a logger (syslogd)
# This is necessary for the 'grafana-openrc' post-install scripts to succeed.
# Also include other common Grafana dependencies.
RUN apk update && \
    apk add --no-cache ca-certificates curl libstdc++ openrc su-exec tzdata fontconfig  ttf-dejavu && \
    # Create minimal OpenRC run directory and softlevel file
    # This allows 'rc-update' (used by some package scripts) to work
    mkdir -p /run/openrc && \
    touch /run/openrc/softlevel && \
    rc-update add devfs sysinit && \
    rc-update add dmesg sysinit && \
    rc-update add hwdrivers sysinit && \
    apk add grafana

# 4. Create necessary directories and set permissions
RUN mkdir -p ${GF_PATHS_DATA} ${GF_PATHS_LOGS} ${GF_PATHS_PLUGINS} ${GF_PATHS_PROVISIONING}/datasources ${GF_PATHS_PROVISIONING}/dashboards ${GF_PATHS_PROVISIONING}/notifiers && \
    chmod -R 777 ${GF_PATHS_DATA} ${GF_PATHS_LOGS} /etc/grafana ${GF_PATHS_PROVISIONING}


#ALLOY isntallation
RUN echo "${EDGE_COMMUNITY_REPO}" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache alloy && \
    sed -i '$d' /etc/apk/repositories
    

# Create directory for Alloy configuration and data
RUN mkdir -p /etc/alloy /var/lib/alloy && \
    chmod -R 777 /etc/alloy /var/lib/alloy



# Volume for persistent data
VOLUME ["${GF_PATHS_DATA}"]


WORKDIR /logger_host
COPY config/* .
COPY ./entry_point.sh /
COPY run_scripts/* /logger_host/


RUN apk update && apk add libc6-compat && \
wget https://github.com/grafana/loki/releases/download/v3.4.4/loki-linux-amd64.zip && \
unzip loki-linux-amd64.zip -d /logger_host/
	
EXPOSE 3100 9096 12345 3000
RUN chmod 777 /entry_point.sh
