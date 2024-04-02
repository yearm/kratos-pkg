package xgormv2

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync/atomic"
)

type (
	IGetter[T any] interface {
		GetByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (*T, error)
		GetListByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) ([]*T, error)
		GetCountByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (int64, error)
	}
	ICreator[T any] interface {
		Create(ctx context.Context, t *T, tx ...*gorm.DB) error
		CreateInBatches(ctx context.Context, ts []*T, batchSize int, tx ...*gorm.DB) error
	}
	IUpdater[T any] interface {
		SaveByScopes(ctx context.Context, t *T, scopes []Scope, tx ...*gorm.DB) error
		UpdateByScopesAndMaps(ctx context.Context, scopes []Scope, maps []Map, tx ...*gorm.DB) (int64, error)
	}
	ITx interface {
		Transaction(fc func(tx *gorm.DB) error) error
	}
)

type BaseRepo[T schema.Tabler] struct {
	tb   T
	rwDB *gorm.DB
	rDBs []*gorm.DB
	next uint32
}

// NewBaseRepo ...
func NewBaseRepo[T schema.Tabler](rwDB *gorm.DB, rDBs []*gorm.DB) *BaseRepo[T] {
	var tb T
	return &BaseRepo[T]{
		tb:   tb,
		rwDB: rwDB,
		rDBs: rDBs,
		next: 0,
	}
}

func (b *BaseRepo[T]) GetRwDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 {
		return tx[0]
	}
	return b.rwDB
}

// GetRDB 默认采用轮询负载
func (b *BaseRepo[T]) GetRDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 {
		return tx[0]
	}
	if len(b.rDBs) == 1 {
		return b.rDBs[0]
	}
	nextIndex := atomic.AddUint32(&b.next, 1)
	return b.rDBs[nextIndex%uint32(len(b.rDBs))]
}

// GetByScopes ...
func (b *BaseRepo[T]) GetByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (*T, error) {
	result := new(T)
	if err := b.GetRDB(tx...).
		Scopes(ToGormScopes(scopes)...).
		WithContext(ctx).
		First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// GetCountByScopes ...
func (b *BaseRepo[T]) GetCountByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (int64, error) {
	var count int64
	if err := b.GetRDB(tx...).
		Model(new(T)).
		Scopes(ToGormScopes(scopes)...).
		WithContext(ctx).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetListByScopes ...
func (b *BaseRepo[T]) GetListByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) ([]*T, error) {
	var results = make([]*T, 0)
	if err := b.GetRDB(tx...).
		Scopes(ToGormScopes(scopes)...).
		WithContext(ctx).
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// Create ...
func (b *BaseRepo[T]) Create(ctx context.Context, t *T, tx ...*gorm.DB) error {
	r := b.GetRwDB(tx...).
		WithContext(ctx).
		Create(t)
	return r.Error
}

// CreateInBatches ...
func (b *BaseRepo[T]) CreateInBatches(ctx context.Context, ts []*T, batchSize int, tx ...*gorm.DB) error {
	r := b.GetRwDB(tx...).
		WithContext(ctx).
		CreateInBatches(ts, batchSize)
	return r.Error
}

// SaveByScopes ...
func (b *BaseRepo[T]) SaveByScopes(ctx context.Context, t *T, scopes []Scope, tx ...*gorm.DB) error {
	r := b.GetRwDB(tx...).
		Scopes(ToGormScopes(scopes)...).
		WithContext(ctx).
		Save(t)
	return r.Error
}

// UpdateByScopesAndMaps ...
func (b *BaseRepo[T]) UpdateByScopesAndMaps(ctx context.Context, scopes []Scope, maps []Map, tx ...*gorm.DB) (int64, error) {
	r := b.GetRwDB(tx...).
		Model(new(T)).
		Scopes(ToGormScopes(scopes)...).
		WithContext(ctx).
		Updates(ToMap(maps))
	return r.RowsAffected, r.Error
}

// Transaction ...
func (b *BaseRepo[T]) Transaction(f func(tx *gorm.DB) error) error {
	return b.rwDB.Transaction(f)
}

// ScopeEqual ...
func (b *BaseRepo[T]) ScopeEqual(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", b.fieldName(k)), v)
	}
}

// ScopeNotEqual ...
func (b *BaseRepo[T]) ScopeNotEqual(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s != ?", b.fieldName(k)), v)
	}
}

// ScopeIN ...
func (b *BaseRepo[T]) ScopeIN(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN(?)", b.fieldName(k)), v)
	}
}

// ScopeNotIN ...
func (b *BaseRepo[T]) ScopeNotIN(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s NOT IN(?)", b.fieldName(k)), v)
	}
}

// ScopeBetween ...
func (b *BaseRepo[T]) ScopeBetween(k string, v1, v2 interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", b.fieldName(k)), v1, v2)
	}
}

// ScopeGT ...
func (b *BaseRepo[T]) ScopeGT(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s > ?", b.fieldName(k)), v)
	}
}

// ScopeGTE ...
func (b *BaseRepo[T]) ScopeGTE(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s >= ?", b.fieldName(k)), v)
	}
}

// ScopeLT ...
func (b *BaseRepo[T]) ScopeLT(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s < ?", b.fieldName(k)), v)
	}
}

// ScopeLTE ...
func (b *BaseRepo[T]) ScopeLTE(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s <= ?", b.fieldName(k)), v)
	}
}

// ScopeLike ...
func (b *BaseRepo[T]) ScopeLike(k string, v interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s LIKE ?", b.fieldName(k)), v)
	}
}

// ScopeOrderBy expr: ASC、DESC、>=0、<= 0 ...
func (b *BaseRepo[T]) ScopeOrderBy(k, expr string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s %s", b.fieldName(k), expr))
	}
}

// ScopeLimit ...
func (b *BaseRepo[T]) ScopeLimit(v int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(v)
	}
}

// ScopeOffset ...
func (b *BaseRepo[T]) ScopeOffset(v int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(v)
	}
}

// ScopeLeftJoin ...
func (b *BaseRepo[T]) ScopeLeftJoin(joinTable string, onTableFiled, onJoinTableField string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(fmt.Sprintf("LEFT JOIN `%s` ON %s = %s", joinTable, b.fieldName(onTableFiled), b.fieldName(onJoinTableField, joinTable)))
	}
}

// ScopeInnerJoin ...
func (b *BaseRepo[T]) ScopeInnerJoin(joinTable string, onTableFiled, onJoinTableField string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.InnerJoins(fmt.Sprintf("INNER JOIN `%s` ON %s = %s", joinTable, b.fieldName(onTableFiled), b.fieldName(onJoinTableField, joinTable)))
	}
}

func (b *BaseRepo[T]) MapKV(k string, v interface{}) Map {
	return Map{fmt.Sprintf("%s", b.fieldName(k)): v}
}

// MapExprAdd ...
func (b *BaseRepo[T]) MapExprAdd(k string, v interface{}) Map {
	return Map{
		fmt.Sprintf("%s", b.fieldName(k)): gorm.Expr(fmt.Sprintf("%s + ?", b.fieldName(k)), v),
	}
}

// MapExprSub ...
func (b *BaseRepo[T]) MapExprSub(k string, v interface{}) Map {
	return Map{
		fmt.Sprintf("%s", b.fieldName(k)): gorm.Expr(fmt.Sprintf("%s - ?", b.fieldName(k)), v),
	}
}

func (b *BaseRepo[T]) fieldName(k string, tableName ...string) string {
	if len(tableName) > 0 {
		return fmt.Sprintf("`%s`.`%s`", tableName[0], k)
	}
	return fmt.Sprintf("`%s`.`%s`", b.tb.TableName(), k)
}
