#!/bin/bash

prefix=~/storage_data

for i in $(seq 1 6); do
  mkdir -p ${prefix}/$i/objects
  mkdir -p ${prefix}/$i/temp
  mkdir -p ${prefix}/$i/garbage
done

mkdir -p ../log

sudo ifconfig wlan0:1 10.29.1.1/16
sudo ifconfig wlan0:2 10.29.1.2/16
sudo ifconfig wlan0:3 10.29.1.3/16
sudo ifconfig wlan0:4 10.29.1.4/16
sudo ifconfig wlan0:5 10.29.1.5/16
sudo ifconfig wlan0:6 10.29.1.6/16
