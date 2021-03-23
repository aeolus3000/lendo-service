#!/bin/bash
docker run -v $PWD:/src --network=host gobuffalo/buffalo:v0.16.21 buffalo "$@"
