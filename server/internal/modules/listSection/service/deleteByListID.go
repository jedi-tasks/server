package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"server/internal/modules/listSection/model"
	"server/internal/utils/necessary"
)

// DeleteByListID мягко удаляет все разделы списка внутри открытой транзакции удаления списка
func (s *ListSectionService) DeleteByListID(ctxTx context.Context, necessary necessary.NecessaryUserInformation, listID uuid.UUID, updatedAt time.Time) (lastAuditLogID uint32, err error) {
	ctxTx, span := tracer.Start(ctxTx, "DeleteByListID")
	defer span.End()

	sections, err := s.listSectionRepository.GetListSections(ctxTx, model.GetListSectionsReq{ //nolint:exhaustruct
		Necessary: necessary,
		ListIDs:   []uuid.UUID{listID},
	})
	if err != nil {
		return 0, err
	}

	for _, section := range sections {
		auditLogID, err := s.deleteListSection(ctxTx, model.DeleteListSectionReq{
			Necessary: necessary,
			ID:        section.ID,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return 0, err
		}
		lastAuditLogID = max(lastAuditLogID, auditLogID)
	}

	return lastAuditLogID, nil
}
