#!/bin/sh

git add ./
git commit -m "update"
git push
git tag $1
git push origin $1
