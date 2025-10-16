package models

type Log struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Timestamp   string `gorm:"not null" json:"timestamp"`
	Method      string `gorm:"size:10;not null" json:"method"`
	Path        string `gorm:"size:255;not null" json:"path"`
	ClientIP    string `gorm:"size:100" json:"client_ip"`
	StatusCode  int    `gorm:"not null" json:"status_code"`
	User        string `gorm:"size:100" json:"user,omitempty"`
	QueryParams string `gorm:"type:text" json:"query_params,omitempty"`
	Error       string `gorm:"type:text" json:"error"`
}
