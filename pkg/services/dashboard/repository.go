package dashboard

import (
	"github.com/labbs/alfred/pkg/database"
	"github.com/labbs/alfred/pkg/exception"
)

type DashboardRepositoryDB struct {
	client database.DbConnection
}

func NewDashboardRepository() DashboardRepositoryDB {
	client := database.GetDbConnection()
	return DashboardRepositoryDB{client: client}
}

func (d DashboardRepositoryDB) GetAllDashboards(userId string) ([]Dashboard, *exception.AppError) {
	var b []Dashboard
	r := d.client.DB.
		Where("user_id = ?", userId).Find(&b)
	if r.Error != nil {
		return []Dashboard{}, exception.NewUnexpectedError("unable to find dashboard(s)", r.Error)
	}
	return b, nil
}

func (d DashboardRepositoryDB) GetDashboardById(userId string, id string) (Dashboard, *exception.AppError) {
	b := Dashboard{}
	r := d.client.DB.
		Preload("Widgets").
		Where("id = ? and user_id = ?", id, userId).First(&b)
	if r.Error != nil {
		return Dashboard{}, exception.NewUnexpectedError("unable to find dashboard", r.Error)
	}
	return b, nil
}

func (d DashboardRepositoryDB) GetDefaultDashboard(userId string) (Dashboard, *exception.AppError) {
	b := Dashboard{}
	r := d.client.DB.
		Preload("Widgets").
		Where("user_id = ? and `default` = 1", userId).First(&b)
	if r.Error != nil {
		return Dashboard{}, exception.NewUnexpectedError("unable to find dashboard", r.Error)
	}
	return b, nil
}

func (d DashboardRepositoryDB) CreateDashboard(b Dashboard) *exception.AppError {
	r := d.client.DB.Create(&b)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to create dashboard", r.Error)
	}
	return nil
}

func (d DashboardRepositoryDB) UpdateDashboard(b Dashboard) *exception.AppError {
	r := d.client.DB.Save(&b)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to update dashboard", r.Error)
	}
	return nil
}

func (d DashboardRepositoryDB) DeleteDashboard(id string, userId string) *exception.AppError {
	r := d.client.DB.Where("id = ? and user_id = ?", id, userId).Delete(&Dashboard{})
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to delete dashboard", r.Error)
	}
	return nil
}

func (d DashboardRepositoryDB) SetDefaultDashboard(id string, userId string) *exception.AppError {
	r := d.client.DB.Model(&Dashboard{}).Where("id = ? and user_id = ?", id, userId).Update("default", true)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to set default dashboard", r.Error)
	}
	return nil
}

func (d DashboardRepositoryDB) UpdateWidget(widget Widget) *exception.AppError {
	r := d.client.DB.Save(&widget)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to update widget", r.Error)
	}
	return nil
}

func (d DashboardRepositoryDB) DeleteWidget(id string, userId string) *exception.AppError {
	r := d.client.DB.Where("id = ? and user_id = ?", id, userId).Delete(&Widget{})
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to delete widget", r.Error)
	}
	return nil
}

func (d DashboardRepositoryDB) CreateWidget(widget Widget) *exception.AppError {
	r := d.client.DB.Create(&widget)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to create widget", r.Error)
	}
	return nil
}

func (d DashboardRepositoryDB) GetWidgetsByDashboardId(dashboardId, userId string) ([]Widget, *exception.AppError) {
	var w []Widget
	r := d.client.DB.
		Where("dashboard_id = ? and user_id = ?", dashboardId, userId).Find(&w)
	if r.Error != nil {
		return []Widget{}, exception.NewUnexpectedError("unable to find widget(s)", r.Error)
	}
	return w, nil
}

func (d DashboardRepositoryDB) GetWidgetById(id, userId string) (Widget, *exception.AppError) {
	w := Widget{}
	r := d.client.DB.
		Where("id = ? and user_id = ?", id, userId).First(&w)
	if r.Error != nil {
		return Widget{}, exception.NewUnexpectedError("unable to find widget", r.Error)
	}
	return w, nil
}
