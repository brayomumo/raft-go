package rpc

import (
	"sync"

	"github.com/sirupsen/logrus"
)


type RPCHandler struct{
	mu sync.Mutex
	Logger *logrus.Logger

}

func NewRPCHandler(logger *logrus.Logger) *RPCHandler{
	return &RPCHandler{
		Logger: logger,
	}
}

/* 
This is the core of the RPC handler. It takes the command name(string) and payload which is 
map for the arguments required for the given command in the form :
 	<argument_name> : <argument_value>
for simplicity.
*/

func (handler *RPCHandler) Handle(command string, payload map[string]interface{} ) string{
	// TODO: Move lock to funs with write/read logic
	handler.mu.Lock()
	// lock so only single goroutine can change data in this handler
	defer handler.mu.Unlock()
	switch command{
	case "APPENDENTRIES":
		handler.Logger.Info("Append Entries Called!")
		return "true"
	case "REQUESTVOTE":
		handler.Logger.Info("Request Vote called")
		return "true"
	default:
		handler.Logger.Warn("Unknown command ", command)
		return "false"
	}
}