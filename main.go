package main

import (
	"time"

	"./register"
	"./common"
)

func main() {
	ip := common.GetIPv4ByIFName("eth1")
	servTypes := common.GetServFromConf("config.ini")

	var w *register.Worker = nil

	// keep trying to connect to etcd and docker daemon
	// until success
	for w == nil {
		w = register.NewWorker(
			ip,
			servTypes,
			[]string{"http://192.168.0.2:4001"},
			"tcp://127.0.0.1:4243")

		if w == nil {
			time.Sleep(1)
		}
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for _ = range ticker.C {
			w.HeartBeat()
		}
	}()
}
