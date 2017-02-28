package app

import (
	"chatoperations/operations"
)

type Notifier interface {
	NotifyRequestSubmitted(o operations.Request) error
	NotifyRequestApproved(o operations.Request) error
	NotifyRequestRejected(o operations.Request) error
}

type OperationsInteractor struct {
	actionStore:  operations.ActionRepository,
	serverStore:  operations.ServerRepository,
	requestStore: operations.RequestRepository,
	notifier:     Notifier,
}

func (ops *OperationsInteractor) SetRequestAction(o operations.Request, actionId int) (operations.Request, error) {
	action, err := ops.actionStore.FindById(actionId)

	o.Action = action

	err := ops.requestStore.Store(o)
	return o, err
}

func (ops *OperationsInteractor) SetRequestServer(requestId int, serverId int) (operations.Request, error) {
	o, err := ops.requestStore.FindById(requestId)
	server, err := ops.serverStore.FindById(serverId)

	o.Server = server

	err := ops.requestStore.Store(o)
	return o, err
}

func (ops *OperationsInteractor) SubmitRequest(requestId int) error {
	o, err := ops.requestStore.FindById(requestId)

	return ops.notifier.NotifyRequestSubmitted(o)
}

func (ops *OperationsInteractor) ApproveRequest(requestId int, responder string) error {
	o, err := ops.requestStore.FindById(requestId)
	o.Approved = true
	o.Responder = responder

	go ops.notifier.NotifyRequestSubmitted(o)
	// go exec request

	err := ops.requestStore.Store(o)
	return o, err
}

func (ops *OperationsInteractor) RejectRequest(requestId int, responder string) error {
	o, err := ops.requestStore.FindById(requestId)
	o.Approved = false
	o.Responder = responder

	go ops.notifier.NotifyRequestRejected(o)

	err := ops.requestStore.Store(o)
	return o, err
}
