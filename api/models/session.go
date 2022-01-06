package models

import (
	"errors"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Title        string `json:"title"`
	Description  string `json:"description"`
	CoachID      int    `json:"coach_id"`
	ClientID     int    `json:"client_id"`
	Year         int    `gorm:"not null" json:"year"`
	Month        int    `gorm:"not null" json:"month"`
	Day          int    `gorm:"not null" json:"day"`
	StartingTime string `gorm:"not null" json:"starting_time"`
	Duration     int    `gorm:"not null" json:"duration"`
	DateSession  string `json:"date_session (YYYY-MM-DD HH:MM)"`
}

func (s *Session) GetSessions(db *gorm.DB) (*[]Session, error) {
	var err error
	sessions := []Session{}
	if err = db.Model(&Session{}).Find(&sessions).Error; err != nil {
		return &[]Session{}, err
	}
	return &sessions, err
}

func (s *Session) GetSessionByID(db *gorm.DB, uid uint32) (*Session, error) {
	var err error
	if err = db.Model(Session{}).Where("id = ?", uid).First(&s).Error; err != nil {
		return &Session{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Session{}, errors.New("Session Not Found")
	}
	return s, err
}

func (s *Session) CreateSession(db *gorm.DB) (*Session, error) {
	var err error
	if err = db.Debug().Create(&s).Error; err != nil {
		return &Session{}, err
	}
	return s, nil
}

func (s *Session) DeleteSession(db *gorm.DB, id uint32) (int64, error) {
	if error := db.Model(&Session{}).Where("id = ?", id).Delete(&Session{}).Error; error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
