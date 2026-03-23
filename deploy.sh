#!/bin/bash

echo "Pullin latets code ..."
git pull origin main

echo "Restarting app..."
sudo systemctl restart belajar-server

echo "Done ..."
