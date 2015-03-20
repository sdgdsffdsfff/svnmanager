#!/bin/bash
set -e
{
    mv /usr/local/tomcat6/webapps/ROOT /opt/bak/ROOT_$(date +%Y-%m-%d-%H:%M:%S) &&
    echo "complete"
} || {
    echo "error"
}