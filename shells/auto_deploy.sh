#!/bin/bash
cd /opt/wings
echo "cd /opt/wings"
echo "--------------------svn up---------------------------"
svn up
sleep 3
echo "--------------------mvn clean start------------------"
mvn clean:clean compile 
echo "--------------------mvn clean end--------------------"
sleep 5
echo "--------------------kill tomcat----------------------"
gao
times=0;
while [ `ps aux|grep java | wc -l` -gt 1 ]
do
        sh /usr/local/tomcat6/bin/shutdown.sh
        echo "shutdown ..."
        sleep 6
        times=$[$times+1]
        echo $times
        if [[ $times -ge 5 && `ps aux|grep java | wc -l` -gt 1 ]]
        then
            ps aux | grep java | awk '{f++;if(NF>12){id=$2}} END {print "kill -9 " id}' |sh
        fi
done

sleep 1
echo "--------------------mv ROOT 2 /opt/bak---------------"
mv /usr/local/tomcat6/webapps/ROOT /opt/bak/ROOT_$(date +%Y-%m-%d-%H:%M:%S)
echo "--------------------mvn war start--------------------"
mvn war:exploded
# 60 = 1 minute
sleep 3
echo "--------------------tstart---------------------------"
ulimit -n 65534;/usr/local/tomcat6/bin/startup.sh
tail -fn 100 /usr/local/tomcat6/logs/catalina.out
