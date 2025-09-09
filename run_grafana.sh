#!/bin/sh


export GF_PATHS_DATA="/var/lib/grafana"
export GF_PATHS_LOGS="/var/log/grafana"
export GF_PATHS_PLUGINS="/var/lib/grafana/plugins"
export GF_PATHS_PROVISIONING="/etc/grafana/provisioning"
export GF_PATHS_CONFIG="/etc/grafana/grafana.ini"
export GF_PATHS_HOME="/usr/share/grafana"

chown -R grafana:grafana \
  "${GF_PATHS_DATA}" \
  "${GF_PATHS_LOGS}" \
  "${GF_PATHS_PLUGINS}" \
  "${GF_PATHS_PROVISIONING}" \
  /etc/grafana


cp grafana.ini /etc/grafana/grafana.ini
cp grafana.ini /usr/share/grafana/conf/defaults.ini
mkdir -p /etc/grafana/provisioning/plugins
mkdir -p /etc/grafana/provisioning/alerting

/usr/sbin/grafana-server --homepath=${GF_PATHS_HOME} \
     --config=${GF_PATHS_CONFIG} \
     cfg:default.log.mode=console \
     cfg:default.paths.data=${GF_PATHS_DATA}, \
     cfg:default.paths.logs=${GF_PATHS_LOGS}, \
     cfg:default.paths.plugins=${GF_PATHS_PLUGINS}, \
     cfg:default.paths.provisioning=${GF_PATHS_PROVISIONING}
