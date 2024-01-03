#!/bin/bash


# TODO: introduce version so as to not overwrite previous folder contents
#
DESTINATION_DIR="$HOME/cue/github.com/DolusMockServer/dolus"

mkdir -p $DESTINATION_DIR

cp -r cue-expectations $DESTINATION_DIR 
