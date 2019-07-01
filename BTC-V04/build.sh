#!/bin/bash
rm blockchain
rm *.db
rm -f *.dat

go build -o blockchain *.go
