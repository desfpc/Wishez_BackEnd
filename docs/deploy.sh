#!/bin/bash
 
git fetch --all
git checkout --force "origin/master"
systemctl stop wishez-backend.service
go build /root/wishez-backend/main.go
go install /root/wishez-backend/main.go
systemctl start wishez-backend.service