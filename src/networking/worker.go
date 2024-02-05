package networking

import (
	"encoding/json"

	"github.com/brayomumo/raft-go/src/rpc"
	zmq "github.com/pebbe/zmq4"
	"github.com/sirupsen/logrus"
)


type RPCWorker struct {
	Logger *logrus.Logger
	Context *zmq.Context
	BackendAddress string
	Name string
	socket *zmq.Socket
	raftHandler rpc.RPCHandler // Used to process commands received by worker, initalized by network module
}

func NewWorker(logger logrus.Logger, address string, ctxt zmq.Context) *RPCWorker{
	return &RPCWorker{
		Context: &ctxt,
		BackendAddress: address,
		Logger: &logger,
		socket: nil,
	}
}

func (w *RPCWorker)Run(){
	// create socket
	socket, err :=  w.Context.NewSocket(zmq.DEALER)
	if err != nil{
		w.Logger.Error(err.Error())
		return
	}
	defer socket.Close()
	w.socket = socket
	// connect to backend
	w.socket.Connect(w.BackendAddress)
	w.socket.SetIdentity(w.Name)
	
	// wait for data, pass to rpc handler, send back response
	for {
		mesg, err := w.socket.RecvMessage(0)
		w.Logger.Info(w.Name, "Worker Got a request: ", mesg)
		if err != nil{
			w.Logger.Error(err)
			continue // We dont want the loop to terminate since this terminates the worker
		}
		sender := mesg[0]
		payload := mesg[1]
		command := payload[1]
		var argsFormat map[string]interface{}
		if err := json.Unmarshal([]byte(payload), &argsFormat); err != nil{
			w.Logger.Error(err)
			continue
		}
		// This expects command: str  and args: map[string]string
		resp := w.raftHandler.Handle(string(command), argsFormat)
		
		// send response
		// SendMessage is send_multipart of []string
		w.socket.SendMessage([]string{sender, resp})

	}
}