#!/bin/bash

git pull origin dev

sudo docker stop ticket-server

sudo docker cp . ticket-server:/app

sudo docker start ticket-server

echo "Starting service..."
sudo docker exec ticket-server /bin/bash -c \
    "cd src && go build main.go && ./main >out.log 2>&1 &"
echo "Deploy done!"
