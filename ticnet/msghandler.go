package ticnet

import (
	"github.com/kanyuanzhi/tialloy-client/ticface"
	"github.com/kanyuanzhi/tialloy-client/utils"
)

type MsgHandler struct {
	Apis             map[uint32]ticface.IRouter // Apis[msgID] = handler
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	TaskQueue        []chan ticface.IRequest
}

func NewMsgHandler() ticface.IMsgHandler {
	return &MsgHandler{
		Apis:             make(map[uint32]ticface.IRouter),
		WorkerPoolSize:   utils.GlobalObject.TcpWorkerPoolSize,
		MaxWorkerTaskLen: utils.GlobalObject.TcpMaxWorkerTaskLen,
		TaskQueue: make([]chan ticface.IRequest, utils.GlobalObject.TcpWorkerPoolSize),
	}
}

func (mh *MsgHandler) DoMsgHandler(request ticface.IRequest) {
	msgID := request.GetMsgID()
	handler, ok := mh.Apis[msgID]
	if !ok {
		utils.GlobalLog.Warnf("api msgID=%d is not found", msgID)
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ticface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		utils.GlobalLog.Warnf("api msgID=%d repeated", msgID)
		return
	}
	mh.Apis[msgID] = router
	utils.GlobalLog.Tracef("api msgID=%d added", msgID)
}

func (mh *MsgHandler) StartOneWorkerPool(workerID int, taskQueue chan ticface.IRequest) {
	utils.GlobalLog.Tracef("worker id=%d started", workerID)
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ticface.IRequest, mh.MaxWorkerTaskLen)
		go mh.StartOneWorkerPool(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ticface.IRequest) {
	workerID := request.GetMsgID() % mh.WorkerPoolSize
	utils.GlobalLog.Tracef("add msgID=%d to workerID=%d", request.GetMsgID(), workerID)
	mh.TaskQueue[workerID] <- request
}
