package main

import (
	"github.com/brayomumo/raft-go/src/networking"
	"github.com/sirupsen/logrus"
)


func main(){
	logger := logrus.New()
	address := "tcp://*:5001"
	
	net := networking.NewRPCMain(*logger, "Node-name", address)
	net.Run()
}