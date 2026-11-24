package repos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	pgmodels "aman/internal/adapters/outbound/postgres/models"
	"aman/internal/domains/app"
	"aman/internal/ports/outbound"
)

var _ outbound.AppRepository = (*AppRepository)(nil)

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
	if _, err := uuid.Parse(id); err != nil {
		return nil, app.ErrNotFound
	}

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
	if err := r.db.Order("created_at DESC").Find(&result).Error; err != nil {
		return nil, err
	}
	apps := make([]app.App, len(result))
	for i, m := range result {
		apps[i] = pgmodels.ToDomainApp(m)
	}
	return apps, nil
}

func (r *AppRepository) FindPage(search string, createdFrom, createdUntil *time.Time, limit, offset int) ([]app.App, int64, error) {
	query := r.db.Model(&pgmodels.App{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("(name ILIKE ? OR url ILIKE ? OR callback_url ILIKE ?)", like, like, like)
	}
	if createdFrom != nil {
		query = query.Where("created_at >= ?", *createdFrom)
	}
	if createdUntil != nil {
		query = query.Where("created_at < ?", *createdUntil)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var result []pgmodels.App
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	apps := make([]app.App, len(result))
	for i, m := range result {
		apps[i] = pgmodels.ToDomainApp(m)
	}
	return apps, total, nil
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
