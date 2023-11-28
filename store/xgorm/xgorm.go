package xgorm

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/samber/lo"
)

type (
	IGetter[T any] interface {
		GetByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (*T, error)
		GetCountByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (int64, error)
		GetListByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) ([]*T, error)
	}
	ICreator[T any] interface {
		Create(ctx context.Context, t *T, tx ...*gorm.DB) error
	}
	IUpdater[T any] interface {
		SaveByScopes(ctx context.Context, t *T, scopes []Scope, tx ...*gorm.DB) error
		UpdateByScopesAndMaps(ctx context.Context, scopes []Scope, maps []Map, tx ...*gorm.DB) error
	}
	ITx interface {
		Transaction(fc func(tx *gorm.DB) error) error
	}
)

type BaseRepo[T any] struct {
	rDB  *gorm.DB
	rwDB *gorm.DB
}

// NewBaseRepo ...
func NewBaseRepo[T any](rDB *gorm.DB, rwDB *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{rDB: rDB, rwDB: rwDB}
}

// GetByScopes ...
func (b *BaseRepo[T]) GetByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (*T, error) {
	result := new(T)
	db := lo.If(len(tx) <= 0, b.rDB).ElseF(func() *gorm.DB { return tx[0] })
	if err := db.Scopes(ToGormScopes(scopes)...).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// GetCountByScopes ...
func (b *BaseRepo[T]) GetCountByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (int64, error) {
	var count int64
	db := lo.If(len(tx) <= 0, b.rDB).ElseF(func() *gorm.DB { return tx[0] })
	if err := db.Model(new(T)).Scopes(ToGormScopes(scopes)...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetListByScopes ...
func (b *BaseRepo[T]) GetListByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) ([]*T, error) {
	var results = make([]*T, 0)
	db := lo.If(len(tx) <= 0, b.rDB).ElseF(func() *gorm.DB { return tx[0] })
	if err := db.Scopes(ToGormScopes(scopes)...).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Create ...
func (b *BaseRepo[T]) Create(ctx context.Context, t *T, tx ...*gorm.DB) error {
	db := lo.If(len(tx) <= 0, b.rwDB).ElseF(func() *gorm.DB { return tx[0] })
	return db.Create(t).Error
}

// SaveByScopes ...
func (b *BaseRepo[T]) SaveByScopes(ctx context.Context, t *T, scopes []Scope, tx ...*gorm.DB) error {
	db := lo.If(len(tx) <= 0, b.rwDB).ElseF(func() *gorm.DB { return tx[0] })
	return db.Scopes(ToGormScopes(scopes)...).Save(t).Error
}

// UpdateByScopesAndMaps ...
func (b *BaseRepo[T]) UpdateByScopesAndMaps(ctx context.Context, scopes []Scope, maps []Map, tx ...*gorm.DB) error {
	db := lo.If(len(tx) <= 0, b.rwDB).ElseF(func() *gorm.DB { return tx[0] })
	return db.Model(new(T)).Scopes(ToGormScopes(scopes)...).Updates(ToMap(maps)).Error
}

// Transaction ...
func (b *BaseRepo[T]) Transaction(f func(tx *gorm.DB) error) error {
	return b.rwDB.Transaction(f)
}

// ScopeEqual ...
func (b *BaseRepo[T]) ScopeEqual(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` = ?", k), v)
	}
}

// ScopeNotEqual ...
func (b *BaseRepo[T]) ScopeNotEqual(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` != ?", k), v)
	}
}

// ScopeIN ...
func (b *BaseRepo[T]) ScopeIN(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` IN(?)", k), v)
	}
}

// ScopeBetween ...
func (b *BaseRepo[T]) ScopeBetween(k string, v1, v2 interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` BETWEEN ? AND ?", k), v1, v2)
	}
}

// ScopeGT ...
func (b *BaseRepo[T]) ScopeGT(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` > ?", k), v)
	}
}

// ScopeGTE ...
func (b *BaseRepo[T]) ScopeGTE(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` >= ?", k), v)
	}
}

// ScopeLT ...
func (b *BaseRepo[T]) ScopeLT(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` < ?", k), v)
	}
}

// ScopeLTE ...
func (b *BaseRepo[T]) ScopeLTE(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` <= ?", k), v)
	}
}

// ScopeLike ...
func (b *BaseRepo[T]) ScopeLike(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("`%s` LIKE ?", k), v)
	}
}

// ScopeOrderBy ...
func (b *BaseRepo[T]) ScopeOrderBy(v string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(v)
	}
}

// ScopeLimit ...
func (b *BaseRepo[T]) ScopeLimit(v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(v)
	}
}

// ScopeOffset ...
func (b *BaseRepo[T]) ScopeOffset(v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(v)
	}
}

func (b *BaseRepo[T]) MapKV(k string, v interface{}) Map {
	return Map{k: v}
}

// MapExprAdd ...
func (b *BaseRepo[T]) MapExprAdd(k string, v interface{}) Map {
	return Map{
		k: gorm.Expr(fmt.Sprintf("`%s` + ?", k), v),
	}
}

// MapExprSub ...
func (b *BaseRepo[T]) MapExprSub(k string, v interface{}) Map {
	return Map{
		k: gorm.Expr(fmt.Sprintf("`%s` - ?", k), v),
	}
}
