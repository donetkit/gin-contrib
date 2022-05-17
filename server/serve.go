package server

import "os"

type IService interface {
	Run()
	Shutdown()
	SetRunMode(mode string)
	StopNotify(sig os.Signal)
}
