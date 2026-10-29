package interceptor

import (
	"context"
	"fmt"
	pkgErrors "pkg/errors"
	"pkg/log"
	"server/internal/enum/errorCategory"
	"server/internal/utils/errors"

	"google.golang.org/grpc"
	_ "google.golang.org/protobuf/types/dynamicpb"

	pb "github.com/jedi-tasks/go-server-grpc/proto"
)

type ErrorHandlerInterceptor struct{}

func NewErrorHandlerInterceptor() *ErrorHandlerInterceptor {
	return &ErrorHandlerInterceptor{}
}

func (i *ErrorHandlerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		res, err := handler(ctx, req)
		if err != nil {
			res, err = handleError(res, err)
		}
		return res, err
	}
}

func (i *ErrorHandlerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := handler(srv, stream); err != nil {

			// Просто логгируем ошибку и возвращаем ошибку GRPC, потому что на этом этапе мы не знаем, какой тип ответа у нас будет
			res, err := handleError(nil, err)
			if err != nil {
				log.Error(err)
				return err
			}

			if err := stream.SendMsg(res); err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	}
}

func handleError(res any, err error) (any, error) {

	// Если ошибки нет, а мы сюда попали, значит какие-то проблемы в чейне вызовов gRPC и логгируем это
	if err == nil {
		log.Error("Мы попали в интерсептор обработки ошибок с пустой ошибкой")
		return res, nil
	}

	// Кастуем ошибку, если она не обернута, внутри оборачиваем ее в InternalServer,
	// на вызовы начнет алертить система оповещения и разработчик быстро пофиксит проблему
	customErr := pkgErrors.CastError(err)

	customErr.DeveloperText = customErr.Err.Error()

	customErr.SystemInfo = log.GetSystemInfo()

	// Логгируем ошибку
	log.LogError(customErr)

	params := make(map[string]string, len(customErr.Params))
	for k, v := range customErr.Params {
		params[k] = fmt.Sprintf("%+v", v) // Приводим к строке
	}

	// Конвертируем категорию ошибки в proto-представление
	protoCategory, err := errorCategory.ConvertToProto(errors.GetCategory(customErr.ErrorType))
	if err != nil {
		log.Error(err)
	}

	protoErr := &pb.Error{
		Category:      protoCategory,
		Message:       customErr.HumanText,
		SystemMessage: customErr.DeveloperText,
		Params:        params,
	}

	return &pb.OnlyError{Error: protoErr}, nil
}
