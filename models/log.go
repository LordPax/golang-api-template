package models

import (
	"fmt"
	"golang-api/services"
	"time"
)

const (
	LOG_INFO  = "INFO"
	LOG_ERROR = "ERROR"
)

type Log struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type" gorm:"type:varchar(100)"`
	Tags      []string  `json:"tags" gorm:"json"`
	Date      time.Time `json:"date"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewLog(t string, tags []string, txt string) *Log {
	return &Log{
		Type: t,
		Tags: tags,
		Text: txt,
		Date: time.Now(),
	}
}

func FindAllLogs(query services.QueryFilter) ([]Log, error) {
	var logs []Log

	err := DB.Model(&Log{}).
		Offset(query.GetSkip()).
		Where(query.GetWhere()).
		Order(query.GetSort()).
		Find(&logs).Error

	return logs, err
}

func PrintLogf(tags []string, format string, v ...any) {
	ws := services.GetWebsocket()
	text := fmt.Sprintf(format, v...)

	_ = NewLog(LOG_INFO, tags, text).Save()
	_ = ws.Room("admin").Emit("log:new", text)
}

func ErrorLogf(tags []string, format string, v ...any) {
	ws := services.GetWebsocket()
	text := fmt.Sprintf(format, v...)

	_ = NewLog(LOG_ERROR, tags, text).Save()
	_ = ws.Room("admin").Emit("log:new", text)
}

func (l *Log) FindOneById(id int) error {
	return DB.Model(&Log{}).First(l, id).Error
}

func (l *Log) FindOne(key string, value any) error {
	return DB.Model(&Log{}).Where(key, value).First(l).Error
}

func (l *Log) Save() error {
	return DB.Model(&Log{}).Save(l).Error
}

func (l *Log) Delete() error {
	return DB.Model(&Log{}).Delete(l).Error
}
