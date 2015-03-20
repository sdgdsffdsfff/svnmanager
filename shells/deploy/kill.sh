#!/bin/bash

isError=false;
times=0;
while [ `ps aux|grep java | wc -l` -gt 1 ]
do
    sh /usr/local/tomcat6/bin/shutdown.sh || isError=true && break
    sleep 6
    times=$[$times+1]
    if [[ $times -ge 5 && `ps aux|grep java | wc -l` -gt 1 ]]
    then
        ps aux | grep java | awk '{f++;if(NF>12){id=$2}} END {print "kill -9 " id}' |sh
    fi
    isError=false
done

if [ isError ]
then
    echo "error"
else
    echo "complete"
fi