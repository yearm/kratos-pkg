package xgorm

import (
	"context"
	"github.com/jinzhu/gorm"
)

type (
	Scope func(db *gorm.DB) *gorm.DB
	Map   map[string]interface{}

	IGetter[T any] interface {
		GetByScopes(ctx context.Context, scopes []Scope, tx ...*gorm.DB) (*T, error)
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

// ToGormScopes ...
func ToGormScopes(scopes []Scope) []func(db *gorm.DB) *gorm.DB {
	_scopes := make([]func(db *gorm.DB) *gorm.DB, 0, len(scopes))
	for _, scope := range scopes {
		scope := scope
		_scopes = append(_scopes, func(db *gorm.DB) *gorm.DB {
			return scope(db)
		})
	}
	return _scopes
}

// ToMap ...
func ToMap(maps []Map) map[string]interface{} {
	m := make(map[string]interface{})
	for _, _map := range maps {
		for k, v := range _map {
			m[k] = v
		}
	}
	return m
}
