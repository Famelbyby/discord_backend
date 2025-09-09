#!/bin/bash

LOKI_PID=
ALLOY_PID=
GRAFANA_PID=
graceful_shutdown() {
 kill -SIGTERM ${LOKI_PID};
 wait ${LOKI_PID};
 kill -SIGTERM ${ALLOY_PID}
 wait ${ALLOY_PID};
 kill -SIGTERM ${GRAFANA_PID}
 wait ${GRAFANA_PID};
 exit 0;
}
trap graceful_shutdown SIGHUP SIGINT SIGQUIT SIGTERM

./run_loki.sh &
LOKI_PID=$!

./run_alloy.sh &
ALLOY_PID=$!

chmod 777 run_grafana.sh
./run_grafana.sh &
GRAFANA_PID=$!

while [ 1 ]; do sleep 15; done;


