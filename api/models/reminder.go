package models

import (
	"errors"

	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	TaskID      uint   `json:"task_id"`
	Description string `gorm:"size:500;not null" json:"description"`
}

func (r *Reminder) GetReminders(db *gorm.DB, id uint64) (*[]Reminder, error) {
	var err error
	reminders := []Reminder{}
	task := Task{}
	db.First(&task, id)
	if err = db.Model(Reminder{}).Where("task_id = ?", task.ID).Find(&reminders).Error; err != nil {
		return &[]Reminder{}, err
	}
	return &reminders, err
}

func (r *Reminder) GetReminderByID(db *gorm.DB, reminderId, taskId uint32) (*Reminder, error) {
	var err error
	task := Task{}
	db.First(&task, taskId)
	if err = db.Model(Reminder{}).Where("id = ? AND task_id = ?", reminderId, taskId).First(&r).Error; err != nil {
		return &Reminder{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Reminder{}, errors.New("Reminder Not Found")
	}
	return r, err
}

func (r *Reminder) CreateReminder(db *gorm.DB, id uint64) (*Reminder, error) {
	var err error
	task := Task{}
	db.First(&task, id)
	r.TaskID = task.ID
	if err = db.Debug().Create(&r).Error; err != nil {
		return &Reminder{}, err
	}
	return r, nil
}

func (r *Reminder) DeleteReminder(db *gorm.DB, reminderId, taskId uint32) (int64, error) {
	var err error
	task := Task{}
	db.First(&task, taskId)
	if err = db.Model(Reminder{}).Where("id = ? AND task_id = ?", reminderId, taskId).Delete(&Reminder{}).Error; err != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
