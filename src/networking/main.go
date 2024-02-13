package networking

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/sirupsen/logrus"
)

type RPCMain struct {
	Logger *logrus.Logger
	Name string
	Address string
	BackendAddress string
	WorkerCount int

}
func NewRPCMain(logger logrus.Logger, name string, address string) *RPCMain{
	return &RPCMain{
		Logger: &logger,
		Name: name,
		Address: address,
		BackendAddress: "ipc://backend",
		WorkerCount: 3,
	}
}


func (rpc *RPCMain) Run(){
	context, err := zmq.NewContext()
	if err != nil{
		rpc.Logger.Error(err.Error())
		return
	}
	defer context.Term()
	// open ROUTER socket for outside comms
	socket, err := context.NewSocket(zmq.ROUTER)
	if err != nil{
		rpc.Logger.Error(err.Error())
		return
	}

	socket.SetIdentity(rpc.Name)
	err  = socket.Bind(rpc.Address)
	if err != nil{
		rpc.Logger.Error("Error Binding frontend socket: ", err)
		return
	}
	defer socket.Close()
	rpc.Logger.Info("Network Interface ready for connections")
	// Open DEALER socket for worker comms
	backend, err := context.NewSocket(zmq.DEALER)
	if err != nil {
		rpc.Logger.Error("error Creating backend socket", err)
		return
	}
	err = backend.Bind(rpc.BackendAddress)
	if err != nil{
		rpc.Logger.Error("Error Binding backend socket", err)
		return
	}
	defer backend.Close()

	// start background workers
	for i := 0; i< rpc.WorkerCount; i++ {
		wkr := NewWorker(*rpc.Logger, rpc.BackendAddress, *context)
		go wkr.Run()
	}
	rpc.Logger.Info("Backend Workers Up and running")

	// manually proxy connections
	err = zmq.Proxy(socket, backend, nil)
	if err != nil{
		rpc.Logger.Error(err)
		return
	}

}