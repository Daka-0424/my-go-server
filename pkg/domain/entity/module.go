package entity

import (
	"reflect"
)

const DEFAULT_DB_ID = 0

func Entity() []any {
	return []any{
		// User
		&User{},
		&UserLoginState{},
		&UserPointSummary{},
		&UserSummaryRelation{},

		// Admin
		&Admin{},

		// VC
		&EarnedPoint{},
		&ImitationPoint{},
		&SpendPointHistory{},
		&SpendPointRelation{},

		// Seed
		&PlatformProduct{},
	}
}

func Seed() []any {
	return []any{
		&PlatformProduct{},
	}
}

type ISeedType interface {
	PlatformProduct
}

func GetEntityFields(entity interface{}) []string {
	var fields []string
	modelType := reflect.TypeOf(entity)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" {
			fields = append(fields, tag)
		} else if field.Name == "Term" {
			termFields := GetEntityFields(Term{})
			fields = append(fields, termFields...)
		}
	}
	return fields
}
