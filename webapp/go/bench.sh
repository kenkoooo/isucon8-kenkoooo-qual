#!/bin/bash -x

DIR=$(cd $(dirname $0); pwd)
cd ${DIR}/../../bench/
./bin_linux/bench
