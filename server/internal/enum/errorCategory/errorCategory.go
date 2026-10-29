package errorCategory

import (
	"pkg/maps"

	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

// mappingProtoToModel содержит соответствие между значениями proto.ErrorCategory и errors.Category
var mappingProtoToModel = map[proto.ErrorCategory]errors.Category{
	proto.ErrorCategory_Internal:           errors.CategoryInternal,
	proto.ErrorCategory_NeedToLogout:       errors.CategoryNeedToLogout,
	proto.ErrorCategory_NeedToSync:         errors.CategoryNeedToSync,
	proto.ErrorCategory_Other:              errors.CategoryOther,
	proto.ErrorCategory_NeedToRefreshToken: errors.CategoryNeedToRefreshToken,
}

// ConvertToProto преобразует errors.Category в proto.ErrorCategory
func ConvertToProto(category errors.Category) (protoCategory proto.ErrorCategory, err error) {

	// Разворачиваем мапу
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	// Получаем значение
	protoCategory, ok := mappingModelToProto[category]
	if !ok {
		return protoCategory, errors.BadRequest.New("ErrorCategory undefined")
	}

	return protoCategory, nil
}

// ProtoErrorCategory - обертка над proto.ErrorCategory для конвертации в модель
type ProtoErrorCategory struct {
	proto.ErrorCategory
}

// ConvertToModel преобразует ProtoErrorCategory в errors.Category
func (p ProtoErrorCategory) ConvertToModel() (category errors.Category, err error) {

	// Проверяем наличие значения
	category, ok := mappingProtoToModel[p.ErrorCategory]
	if !ok {
		return category, errors.BadRequest.New("ErrorCategory undefined")
	}

	return category, nil
}
