package models

import (
	"errors"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	AssignerID  int        `json:"assigner_id"`
	AssigneeID  int        `json:"assignee_id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Description string     `json:"description"`
	Reminders   []Reminder `gorm:"constraint:foreignKey:TaskID,OnUpdate:CASCADE,OnDelete:SET NULL;" json:"reminders"`
	Year        int        `gorm:"not null" json:"year"`
	Month       int        `gorm:"not null" json:"month"`
	Day         int        `gorm:"not null" json:"day"`
	Duration    int        `gorm:"not null" json:"duration"`
	DateTask    string     `json:"date_task (YYYY-MM-DD)"`
}

func (t *Task) GetTasks(db *gorm.DB) (*[]Task, error) {
	var err error
	tasks := []Task{}
	if err = db.Model(&Task{}).Find(&tasks).Error; err != nil {
		return &[]Task{}, err
	}
	return &tasks, err
}

func (t *Task) GetTaskByID(db *gorm.DB, uid uint32) (*Task, error) {
	var err error
	if err = db.Model(Task{}).Where("id = ?", uid).First(&t).Error; err != nil {
		return &Task{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Task{}, errors.New("Task Not Found")
	}
	return t, err
}

func (t *Task) CreateTask(db *gorm.DB) (*Task, error) {
	var err error
	if err = db.Debug().Create(&t).Error; err != nil {
		return &Task{}, err
	}
	return t, nil
}

func (t *Task) DeleteTask(db *gorm.DB, id uint32) (int64, error) {
	if error := db.Model(&Task{}).Where("id = ?", id).Delete(&Task{}).Error; error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
