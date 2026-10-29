package osType

import (
	"pkg/maps"
	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

type OSType string

const (
	Android OSType = "android"
	IOS     OSType = "iOS"
	IPadOS  OSType = "iPadOS"
	OSX     OSType = "OSX"
	WatchOS OSType = "WatchOS"
	Server  OSType = "server" // Устройство на стороне сервера, например MCP
)

// mappingProtoToModel содержит соответствие между значениями proto.OSType и OSType
var mappingProtoToModel = map[proto.OSType]OSType{
	proto.OSType_Unspecified: Server,
	proto.OSType_Android:     Android,
	proto.OSType_IOS:         IOS,
	proto.OSType_IPadOS:      IPadOS,
	proto.OSType_OSX:         OSX,
	proto.OSType_WatchOS:     WatchOS,
}

// ConvertToProto преобразует OSType в proto.OSType
func (b OSType) ConvertToProto() (osType proto.OSType, err error) {

	// Разворачиваем мапу
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	// Получаем значение
	protoOSType, ok := mappingModelToProto[b]
	if !ok {
		return protoOSType, errors.BadRequest.New("OSType undefined")
	}

	return protoOSType, nil
}

type ProtoOSType struct {
	proto.OSType
}

// ConvertToModel преобразует ProtoOSType в OSType
func (p ProtoOSType) ConvertToModel() (osType OSType, err error) {

	// Проверяем наличие значения
	osType, ok := mappingProtoToModel[p.OSType]
	if !ok {
		return osType, errors.BadRequest.New("OSType undefined")
	}

	return osType, nil
}
