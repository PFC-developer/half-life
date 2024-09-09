#!/bin/sh

echo "-- using the following"
echo 'CONFIGFILE='${CONFIGFILE}
if [ "$CONFIGFILE" != "${CONFIGFILE#http://}" ] ; then
    curl -s -L $CONFIGFILE -o /status/config.yaml
elif [ "$CONFIGFILE" != "${CONFIGFILE#https://}" ]  ;then
    curl  -L $CONFIGFILE -o /status/config.yaml
elif [ -d "$CONFIGFILE" ] ;then
    yq eval-all '. as $item ireduce ({}; . *+ $item)' /config/*.yaml > /status/config.yaml
else
    cp $CONFIGFILE /status/config.yaml
fi

/usr/bin/halflife monitor -f /status/config.yaml -s /status/status.yaml
