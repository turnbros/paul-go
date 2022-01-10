#!/bin/sh

CMD_DIR="./cmd"
BUILD_DIR="./dist"

mkdir $BUILD_DIR
for file in ${CMD_DIR}/*
do
  go build -o $BUILD_DIR $file
done