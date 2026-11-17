package repos

import (
	"gorm.io/gorm"

	pgmodels "aman/internal/adapters/outbound/postgres/models"
	"aman/internal/domains/app"
	"aman/internal/ports/outbound"
)

var _ outbound.IAppRepository = (*AppRepository)(nil)

type AppRepository struct {
	db *gorm.DB
}

func NewAppRepository(db *gorm.DB) *AppRepository {
	return &AppRepository{db: db}
}

func (r *AppRepository) Create(a app.App) (app.App, error) {
	m := pgmodels.ToModelApp(a)
	if err := r.db.Create(&m).Error; err != nil {
		return app.App{}, err
	}
	return pgmodels.ToDomainApp(m), nil
}

func (r *AppRepository) FindByID(id string) (*app.App, error) {
	var m pgmodels.App
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, app.ErrNotFound
		}
		return nil, err
	}
	a := pgmodels.ToDomainApp(m)
	return &a, nil
}

func (r *AppRepository) FindAll() ([]app.App, error) {
	var result []pgmodels.App
	if err := r.db.Find(&result).Error; err != nil {
		return nil, err
	}
	apps := make([]app.App, len(result))
	for i, m := range result {
		apps[i] = pgmodels.ToDomainApp(m)
	}
	return apps, nil
}

func (r *AppRepository) Update(a app.App) (app.App, error) {
	m := pgmodels.ToModelApp(a)
	if err := r.db.Save(&m).Error; err != nil {
		return app.App{}, err
	}
	return pgmodels.ToDomainApp(m), nil
}

func (r *AppRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&pgmodels.App{}).Error
}
