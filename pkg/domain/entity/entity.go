package entity

import (
	"golang.org/x/exp/slices"
)

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
