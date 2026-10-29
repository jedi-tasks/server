package model

import "github.com/jedi-tasks/go-server-grpc/proto"

type ApplicationInformation struct {
	BundleID string `json:"bundleID" validate:"required" db:"application_bundle_id"` // Бандл приложения
	Version  string `json:"version" validate:"required" db:"application_version"`    // Версия приложения
	Build    string `json:"build" validate:"required" db:"application_build"`        // Билд приложения
}

type ProtoApplicationInformation struct {
	*proto.ApplicationInformation
}

func (p ProtoApplicationInformation) ConvertToModel() (ApplicationInformation, error) {
	if p.ApplicationInformation == nil {
		return ApplicationInformation{}, nil
	}

	return ApplicationInformation{
		BundleID: p.BundleID,
		Version:  p.Version,
		Build:    p.Build,
	}, nil
}
