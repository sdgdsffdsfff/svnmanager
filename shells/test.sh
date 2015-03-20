#!/bin/bash

a=1;
while [ true ]
do
    if [ $a -gt 5 ]
    then
        cd /home/aaa || a=10 && break
    else
        a=6
    fi
done

echo $a