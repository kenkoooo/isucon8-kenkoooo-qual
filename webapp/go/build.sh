#!/bin/sh -x
docker run --rm -it -v $(pwd):/go golang make
