package helpers

import (
	"fmt"

	"gorm.io/gorm"
)

func TruncateDb(db *gorm.DB, schema any) error {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(schema)
	tx := db.Exec(fmt.Sprintf("TRUNCATE TABLE  %s;", stmt.Schema.Table))
	return tx.Error
}
