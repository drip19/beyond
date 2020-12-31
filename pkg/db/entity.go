package db

import (
	"gorm.io/gorm"
)

type InchAction struct {
	gorm.Model
	Pair        string  `json:"pair"`
	InPut       string  `json:"input"`
	OutPut      string  `json:"output"`
	Gas         int64   `json:"gas"`
	Price       string  `json:price`
	Profit      float64 `json:"profit"`
	TxHash      string  `json:"tx_hash"`
	IncludeSake bool    `json:"includeSake"`
}
