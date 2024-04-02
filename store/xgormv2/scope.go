package xgormv2

import (
	"gorm.io/gorm"
)

type (
	Scope func(db *gorm.DB) *gorm.DB
	Map   map[string]interface{}
)

// ToGormScopes ...
func ToGormScopes(scopes []Scope) []func(db *gorm.DB) *gorm.DB {
	ss := make([]func(db *gorm.DB) *gorm.DB, 0, len(scopes))
	for _, scope := range scopes {
		scope := scope
		ss = append(ss, func(db *gorm.DB) *gorm.DB {
			return scope(db)
		})
	}
	return ss
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
