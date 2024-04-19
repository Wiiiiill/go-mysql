package orm

import (
	"database/sql"

	"github.com/Wiiiiill/go-mysql/orm"
)

func Open(url string) (*orm.BaseDB, error) {
	base := &orm.BaseDB{}
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	base.WriteDatabase = db
	base.ReadDatabases = make([]*sql.DB, 2)
	for i := 0; i < 2; i++ {
		db, err := sql.Open("mysql", url)
		if err != nil {
			return nil, err
		}
		base.ReadDatabases = append(base.ReadDatabases, db)
	}

	return base, nil
}
