#!/bin/bash

source ../common/env.sh

vtbench \
  --config-file-not-found-handling ignore \
	--protocol mysql \
	--host 127.0.0.1 \
	--port 15306 \
	--db commerce@primary \
	--sql "select * from customer where customer_id = :random_int" \
	--threads 1 \
	--count 50000
