package model

import "github.com/jedi-tasks/go-server-grpc/proto"

type Version struct {
	Version string `json:"version"` // Версия приложения
	Build   string `json:"build"`   // Номер сборки
}

func (v Version) ConvertToProto() (res *proto.Version, err error) {
	return &proto.Version{
		Version: v.Version,
		Build:   v.Build,
	}, nil
}
