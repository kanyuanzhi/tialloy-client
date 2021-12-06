package ticnet

import (
	"github.com/kanyuanzhi/tialloy-client/ticface"
)

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request ticface.IRequest) {}

func (br *BaseRouter) Handle(request ticface.IRequest) {}

func (br *BaseRouter) PostHandle(request ticface.IRequest) {}
