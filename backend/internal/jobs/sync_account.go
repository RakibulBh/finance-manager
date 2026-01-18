package jobs

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeSyncAccount = "sync:account"
)

type SyncAccountPayload struct {
	FamilyID uuid.UUID `json:"family_id"`
	ItemID   string    `json:"item_id"`
}

func NewSyncAccountTask(familyID uuid.UUID, itemID string) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncAccountPayload{
		FamilyID: familyID,
		ItemID:   itemID,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSyncAccount, payload), nil
}
