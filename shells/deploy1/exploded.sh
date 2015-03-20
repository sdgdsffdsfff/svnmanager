#!/bin/bash
{
    cd /opt/wings &&
    echo "complete"
} || {
    echo "error"
}