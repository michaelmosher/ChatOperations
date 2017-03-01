package application

import (
	"chatoperations/operations"
)

type Notifier interface {
	NotifyRequestSubmitted(o operations.Request) error
	NotifyRequestApproved(o operations.Request) error
	NotifyRequestRejected(o operations.Request) error
}

type OperationsInteractor struct {
	ActionStore  operations.ActionRepository
	ServerStore  operations.ServerRepository
	RequestStore operations.RequestRepository
	Notifier     Notifier
}

func (ops *OperationsInteractor) loadRequest(requestId int) (operations.Request, error) {
	o, err := ops.RequestStore.FindById(requestId)

	if err != nil {
		return o, err
	}

	o.Action, err = ops.ActionStore.FindById(int(o.Action.Id))
	o.Server, err = ops.ServerStore.FindById(int(o.Server.Id))

	return o, err
}

func (ops *OperationsInteractor) ActionOptions() []operations.Action {
	actions, _ := ops.ActionStore.FindAll()
	return actions
}

func (ops *OperationsInteractor) ServerOptions() []operations.Server {
	 servers, _ := ops.ServerStore.FindAll()
	 return servers
}

func (ops *OperationsInteractor) SetRequestAction(o operations.Request, actionId int) (operations.Request, error) {
	action, err := ops.ActionStore.FindById(actionId)
	o.Action = action

	requestId, err := ops.RequestStore.Store(o)
	o.Id = requestId

	return o, err
}

func (ops *OperationsInteractor) SetRequestServer(requestId int, serverId int) (operations.Request, error) {
	o, err := ops.loadRequest(requestId)
	server, err := ops.ServerStore.FindById(serverId)

	o.Server = server

	_, err = ops.RequestStore.Store(o)
	return o, err
}

func (ops *OperationsInteractor) SubmitRequest(requestId int) error {
	o, err := ops.loadRequest(requestId)

	if err != nil {
		return err
	}

	return ops.Notifier.NotifyRequestSubmitted(o)
}

func (ops *OperationsInteractor) ApproveRequest(requestId int, responder string) (operations.Request, error) {
	o, err := ops.loadRequest(requestId)
	o.Approved = true
	o.Responder = responder

	go ops.Notifier.NotifyRequestSubmitted(o)
	// go exec request

	_, err = ops.RequestStore.Store(o)
	return o, err
}

func (ops *OperationsInteractor) RejectRequest(requestId int, responder string) (operations.Request, error) {
	o, err := ops.loadRequest(requestId)
	o.Approved = false
	o.Responder = responder

	go ops.Notifier.NotifyRequestRejected(o)

	_, err = ops.RequestStore.Store(o)
	return o, err
}