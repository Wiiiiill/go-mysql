package orm

import "database/sql"

type BaseDB struct {
	WriteDatabase *sql.DB
	ReadDatabases []*sql.DB
}

type Orm struct {
	BaseDB
}

func (obj *BaseDB) Close() (err error) {
	if err = obj.WriteDatabase.Close(); err != nil {
		return err
	}
	for _, v := range obj.ReadDatabases {
		if err = v.Close(); err != nil {
			return err
		}
	}
	return err
}

func Open(db BaseDB) {

}
