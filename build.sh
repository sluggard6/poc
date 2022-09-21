#!/bin/bash

cd ./assets
rm -rf ./dist
pnpm install
pnpm build
cd ..
statik -src=./assets/dist
