#!/bin/bash
docker run -v $PWD:/src --env-file config.env --network=host gobuffalo/buffalo:v0.16.21 buffalo "$@"
