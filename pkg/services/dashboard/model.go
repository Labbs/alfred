package dashboard

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/labbs/alfred/pkg/exception"
)

type Dashboard struct {
	Id      string `json:"id" gorm:"primaryKey"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
	UserId  string `gorm:"index" json:"-"`

	Widgets []Widget `json:"widgets" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Widget struct {
	Id   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"w"`
	H    int    `json:"h"`
	HTML string `gorm:"type:longtext" json:"html"`
	CSS  string `gorm:"type:longtext" json:"css"`
	JS   string `gorm:"type:longtext" json:"js"`

	DashboardId string `gorm:"index" json:"-"`
	UserId      string `gorm:"index" json:"-"`
}

type ConfigurationMap map[string]interface{}

func (sla *ConfigurationMap) Scan(value interface{}) error {
	var skills map[string]interface{}
	err := json.Unmarshal([]byte(value.(string)), &skills)
	if err != nil {
		return err
	}
	*sla = skills
	return nil
}

func (sla ConfigurationMap) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type DashboardRepository interface {
	GetAllDashboards(userId string) ([]Dashboard, *exception.AppError)
	GetDashboardById(userId string, id string) (Dashboard, *exception.AppError)
	GetDefaultDashboard(userId string) (Dashboard, *exception.AppError)
	CreateDashboard(d Dashboard) *exception.AppError
	UpdateDashboard(d Dashboard) *exception.AppError
	DeleteDashboard(id string, userId string) *exception.AppError
	SetDefaultDashboard(id string, userId string) *exception.AppError
	UpdateWidget(widget Widget) *exception.AppError
	DeleteWidget(id string, userId string) *exception.AppError
	CreateWidget(widget Widget) *exception.AppError
	GetWidgetsByDashboardId(dashboardId, userId string) ([]Widget, *exception.AppError)
	GetWidgetById(id, userId string) (Widget, *exception.AppError)
}
