package database

import "strings"

type ColumnInfo struct {
	Field   string `json:"field"`
	Type    string `json:"type"`
	IsNull  string `json:"is_null"`
	Key     string `json:"key"` //是否主键
	Default string `json:"default"`
	Comment string `json:"comment"`
}

func (c ColumnInfo) GetGoType() string {

	if strings.Index(c.Type, "int") >= 0 || strings.Index(c.Type, "tinyint") >= 0 {
		return "int"
	}

	if strings.Index(c.Type, "text") >= 0 || strings.Index(c.Type, "varchar") >= 0 || strings.Index(c.Type, "nvarchar") >= 0 {
		return "string"
	}

	if strings.Index(c.Type, "decimal") >= 0 || strings.Index(c.Type, "double") >= 0 {
		return "float64"
	}

	return "interface{}"
}
