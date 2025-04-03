#!/bin/bash

prefix=~/storage_data

for i in $(seq 1 6); do
  mkdir -p ${prefix}/$i/objects
  mkdir -p ${prefix}/$i/temp
  mkdir -p ${prefix}/$i/garbage
done

sudo ifconfig wlan0:1 10.29.1.1/16
sudo ifconfig wlan0:2 10.29.1.2/16
sudo ifconfig wlan0:3 10.29.1.3/16
sudo ifconfig wlan0:4 10.29.1.4/16
sudo ifconfig wlan0:5 10.29.1.5/16
sudo ifconfig wlan0:6 10.29.1.6/16

#以下适用于mac

## 为 en0 接口创建别名并分配 IP 地址
#sudo ifconfig en0 alias 10.29.1.1/16
#sudo ifconfig en0 alias 10.29.1.2/16
#sudo ifconfig en0 alias 10.29.1.3/16
#sudo ifconfig en0 alias 10.29.1.4/16
#sudo ifconfig en0 alias 10.29.1.5/16
#sudo ifconfig en0 alias 10.29.1.6/16
#sudo ifconfig en0 alias 10.29.2.1/16
#sudo ifconfig en0 alias 10.29.2.2/16
#sudo ifconfig en0 alias 10.29.3.1/16
