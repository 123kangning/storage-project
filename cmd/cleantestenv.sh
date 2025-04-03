#!/bin/bash

prefix=~/storage_data

for i in `seq 1 6`
do
    rm -rf ${prefix}/$i/objects/*
    rm -rf ${prefix}/$i/temp/*
    rm -rf ${prefix}/$i/garbage/*
done
