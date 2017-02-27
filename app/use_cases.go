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
	action := ops.actionStore.FindById(actionId)
	o.Action = action

	err := ops.requestStore.Store(o)
	return o, err
}

func (ops *OperationsInteractor) SetRequestServer(o operations.Request, serverId int) (operations.Request, error) {
	server := ops.serverStore.FindById(serverId)
	o.Server = server

	err := ops.requestStore.Store(o)
	return o, err
}

func (ops *OperationsInteractor) SubmitRequest(o operations.Request) error {
	return ops.notifier.NotifyRequestSubmitted(o)
}

func (ops *OperationsInteractor) ApproveRequest(o operations.Request) error {
	o.Approved = true

	go ops.notifier.NotifyRequestSubmitted(o)
	// go exec request

	err := ops.requestStore.Store(o)
	return o, err
}

func (ops *OperationsInteractor) RejectRequest(o operations.Request) error {
	o.Approved = false

	go ops.notifier.NotifyRequestRejected(o)

	err := ops.requestStore.Store(o)
	return o, err
}
