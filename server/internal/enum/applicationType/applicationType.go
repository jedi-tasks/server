package applicationType

import (
	"context"
	"pkg/maps"

	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

type ApplicationType string

// enums:"ios,android,web,server"
const (
	IOS     = ApplicationType("ios")
	Android = ApplicationType("android")
	Web     = ApplicationType("web")
	Server  = ApplicationType("server")
)

func (t *ApplicationType) Validate(ctx context.Context) error {
	if t == nil {
		return nil
	}
	switch *t {
	case IOS, Android, Web, Server:
	default:
		return errors.BadRequest.New("Unknown application type").
			WithContextParams(ctx).
			SkipThisCall().
			WithParams("type", *t).
			WithCustomHumanText("Неизвестный тип приложения")
	}
	return nil
}

// mappingProtoToModel содержит соответствие между значениями proto.ApplicationType и ApplicationType
var mappingProtoToModel = map[proto.ApplicationType]ApplicationType{
	proto.ApplicationType_IOS:     IOS,
	proto.ApplicationType_Android: Android,
	proto.ApplicationType_Web:     Web,
	proto.ApplicationType_Server:  Server,
}

// ConvertToProto преобразует ApplicationType в proto.ApplicationType
func (b ApplicationType) ConvertToProto() (applicationType proto.ApplicationType, err error) {

	// Разворачиваем мапу
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	// Получаем значение
	protoApplicationType, ok := mappingModelToProto[b]
	if !ok {
		return protoApplicationType, errors.BadRequest.New("ApplicationType undefined")
	}

	return protoApplicationType, nil
}

type ProtoApplicationType struct {
	proto.ApplicationType
}

// ConvertToModel преобразует ProtoApplicationType в ApplicationType
func (p ProtoApplicationType) ConvertToModel() (applicationType ApplicationType, err error) {

	// Проверяем наличие значения
	applicationType, ok := mappingProtoToModel[p.ApplicationType]
	if !ok {
		return applicationType, errors.BadRequest.New("ApplicationType undefined")
	}

	return applicationType, nil
}
