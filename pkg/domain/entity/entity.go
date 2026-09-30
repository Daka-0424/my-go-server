package entity

import (
	"time"

	"golang.org/x/exp/slices"
)

type EntityBase struct {
	ID        uint      `yaml:"id" gorm:"primaryKey;autoIncrement:true"`
	CreatedAt time.Time `yaml:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `yaml:"updated_at" gorm:"autoUpdateTime"`
}

// gromのUpdate用補機、更新したColumnを保持する
type gormAuxiliary struct {
	updateColumns []string
}

func (g *gormAuxiliary) updateColumn(columns ...string) {
	if g.updateColumns == nil {
		g.updateColumns = []string{}
	}
	for _, col := range columns {
		if !slices.Contains(g.updateColumns, col) {
			g.updateColumns = append(g.updateColumns, col)
		}
	}
}

func (g *gormAuxiliary) GetUpdateColumns() []string {
	if len(g.updateColumns) == 0 {
		return []string{}
	}
	return g.updateColumns
}
