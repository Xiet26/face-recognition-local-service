package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"xiet26/goface/face-management/database"
	"xiet26/goface/face-management/model"
)

type GetAttendTemp struct {
	BatchID string `json:"batchID"`
}

func (c *GetAttendTemp) Valid() error {
	if c.BatchID == "" {
		return fmt.Errorf("invalid struct")
	}

	_, err := govalidator.ValidateStruct(c)
	return err
}

type GetAttendTempHandler struct {
	AttendTempRepository database.AttendTempMongoRepository
}

func (h *GetAttendTempHandler) Handle(data *GetAttendTemp) ([]model.AttendTemp, error) {
	if err := data.Valid(); err != nil {
		return nil, err
	}

	return h.AttendTempRepository.ReadByBatchID(data.BatchID)
}
