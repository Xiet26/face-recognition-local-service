package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"xiet26/Smart_Attendance_System/face-management/database"
)

type DeleteAttendTemp struct {
	BatchID string `json:"batchID"`
}

func (c *DeleteAttendTemp) Valid() error {
	if c.BatchID == "" {
		return fmt.Errorf("invalid struct")
	}

	_, err := govalidator.ValidateStruct(c)
	return err
}

type DeleteAttendTempHandler struct {
	AttendTempRepository database.AttendTempMongoRepository
}

func (h *DeleteAttendTempHandler) Handle(data *DeleteAttendTemp) error {
	if err := data.Valid(); err != nil {
		return err
	}

	return h.AttendTempRepository.DeleteByBatchID(data.BatchID)
}
