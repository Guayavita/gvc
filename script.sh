#!/bin/bash

./bin/guayavita 

./bin/guayavita compile --help

./bin/guayavita compile -o test-bin/ ./test-data/hello-simple.gvt --benchmark

ls -alh ./test-bin/hello-simple

./test-bin/hello-simple
