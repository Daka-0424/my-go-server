package entity

import (
	"reflect"
	"strings"
	"unicode"
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
	modelType := reflect.TypeOf(User{})
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("gorm")
		if field.Name != "Model" && !strings.Contains(tag, "foreignkey") {
			fields = append(fields, field.Name)
		}
	}
	return fields
}

func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
