package config

import (
	"github.com/gophero/goal/stringx"
	"strconv"
)

const (
	WDB = "wdb"
	RDB = "rdb"
)

type DBs struct {
	DBs []DB `mapstructure:"db" yaml:"db"`
}

func (d DBs) Wdb() DB {
	for _, v := range d.DBs {
		if v.Name == WDB {
			return v
		}
	}
	return DB{}
}

func (d DBs) Rdb() DB {
	for _, v := range d.DBs {
		if v.Name == RDB {
			return v
		}
	}
	return DB{}
}

type DB struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Database string `mapstructure:"database" yaml:"database"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	ShowSql  bool   `mapstructure:"show-sql" yaml:"show-sql"`
	Params   string `mapstructure:"params" yaml:"params"`
}

func (db DB) Dsn() string {
	var p = db.Params
	if p != "" && !stringx.StartsWith(p, "?") {
		p = "?" + p
	}
	return db.Username + ":" + db.Password + "@tcp(" + db.Host + ":" + strconv.Itoa(db.Port) + ")/" + db.Database + p
}
