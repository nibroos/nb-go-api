package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nibroos/nb-go-api/service/internal/dtos"
	"github.com/nibroos/nb-go-api/service/internal/models"
	"github.com/nibroos/nb-go-api/service/internal/utils"
	"gorm.io/gorm"
)

type IdentifierRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewIdentifierRepository(db *gorm.DB, sqlDB *sqlx.DB) *IdentifierRepository {
	return &IdentifierRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *IdentifierRepository) ListIdentifiers(ctx context.Context, filters map[string]string) ([]dtos.IdentifierListDTO, int, error) {
	identifiers := []dtos.IdentifierListDTO{}
	var total int

	from := `FROM (
        SELECT i.id, i.user_id, i.type_identifier_id, i.ref_num, i.status, i.created_at, i.updated_at,
        u.name as user_name,
        ti.name as type_identifier_name

        FROM identifiers i
        JOIN users u ON i.user_id = u.id
        JOIN mix_values ti ON i.type_identifier_id = ti.id
        WHERE i.deleted_at IS NULL
    ) AS alias WHERE 1=1`

	query := `SELECT * ` + from
	countQuery := `SELECT COUNT(*) ` + from

	var args []interface{}
	i := 1
	for key, value := range filters {
		switch key {
		case "ref_num":
			if value != "" {
				query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				countQuery += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				args = append(args, "%"+value+"%")
				i++
			}
		}
	}

	if value, ok := filters["user_id"]; ok && value != "" {
		query += fmt.Sprintf(" AND user_id = $%d", i)
		countQuery += fmt.Sprintf(" AND user_id = $%d", i)
		args = append(args, value)
		i++
	}

	if value, ok := filters["global"]; ok && value != "" {
		query += fmt.Sprintf(" AND (ref_num ILIKE $%d OR user_name ILIKE $%d OR type_identifier_name ILIKE $%d)", i, i+1, i+2)
		countQuery += fmt.Sprintf(" AND (ref_num ILIKE $%d OR user_name ILIKE $%d OR type_identifier_name ILIKE $%d)", i, i+1, i+2)
		args = append(args, "%"+value+"%", "%"+value+"%", "%"+value+"%")
		i += 3
	}

	countArgs := append([]interface{}{}, args...)

	allowedOrderColumns := []string{"id", "ref_num", "user_name", "type_identifier_name"}
	orderColumn := utils.GetStringOrDefaultFromArray(filters["order_column"], allowedOrderColumns, "id")
	orderDirection := utils.GetStringOrDefault(filters["order_direction"], "asc")
	query += fmt.Sprintf(" ORDER BY %s %s", orderColumn, orderDirection)

	perPage := utils.GetIntOrDefault(filters["per_page"], 10)
	currentPage := utils.GetIntOrDefault(filters["page"], 1)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, perPage, (currentPage-1)*perPage)

	// Channels for concurrent execution
	countChan := make(chan error)
	selectChan := make(chan error)

	// Goroutine for count query
	go func() {
		err := r.sqlDB.GetContext(ctx, &total, countQuery, countArgs...)
		countChan <- err
	}()

	// Goroutine for select query
	go func() {
		err := r.sqlDB.SelectContext(ctx, &identifiers, query, args...)
		selectChan <- err
	}()

	// Wait for both goroutines to finish
	countErr := <-countChan
	selectErr := <-selectChan

	if countErr != nil {
		return nil, 0, countErr
	}

	if selectErr != nil {
		return nil, 0, selectErr
	}

	return identifiers, total, nil
}

func (r *IdentifierRepository) GetIdentifierByID(ctx context.Context, params *dtos.GetIdentifierParams) (*dtos.IdentifierDetailDTO, error) {
	var identifier dtos.IdentifierDetailDTO
	// deletedAt := params.IsDeleted

	query := `SELECT i.id, i.user_id, i.type_identifier_id, i.ref_num, i.status, i.created_at, i.updated_at, i.deleted_at,
	u.name as user_name,
	ti.name as type_identifier_name

	FROM identifiers i
	JOIN users u ON i.user_id = u.id
	JOIN mix_values ti ON i.type_identifier_id = ti.id
	WHERE 1=1`

	var args []interface{}

	i := 1
	query += " AND i.id = $1"
	args = append(args, params.ID)
	i++

	isDeletedQuery := ` AND i.deleted_at IS NULL`
	if params.IsDeleted != nil && *params.IsDeleted == 1 {
		isDeletedQuery = " AND i.deleted_at IS NOT NULL"
	}

	query += isDeletedQuery

	if err := r.sqlDB.Get(&identifier, query, args...); err != nil {
		return nil, err
	}

	return &identifier, nil
}

// BeginTransaction starts a new transaction
func (r *IdentifierRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *IdentifierRepository) CreateIdentifier(tx *gorm.DB, identifier *models.Identifier) error {
	if err := tx.Create(identifier).Error; err != nil {
		return err
	}
	return nil
}

func (r *IdentifierRepository) UpdateIdentifier(tx *gorm.DB, identifier *models.Identifier) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(identifier).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *IdentifierRepository) DeleteIdentifier(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Identifier{}, id).Error; err != nil {
		if err := tx.Delete(&models.Identifier{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *IdentifierRepository) RestoreIdentifier(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE identifiers SET deleted_at = NULL WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
