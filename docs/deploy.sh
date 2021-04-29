#!/bin/bash
 
git fetch --all
git checkout --force "origin/master"
go build /root/wishez-backend/main.go