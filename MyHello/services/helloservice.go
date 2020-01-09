package services

import (
	"MyHello/protos"
	"context"
)

type ConfigService struct {
	BaseService
}

func (c ConfigService) Hello(context context.Context, request *proto.Request) (*proto.Response, error) {
	var model proto.Response
	model.Msg = "ok"
	return &model, nil
}
