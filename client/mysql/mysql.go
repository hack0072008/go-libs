package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	Options struct {
		Driver  string `yaml:"driver" mapstructure:"driver"`
		ConnStr string `yaml:"dsn" mapstructure:"dsn"`
		// 定时保活
		KeepAlive int `yaml:"keep_alive" mapstructure:"keep_alive"`
		// 最大可空闲连接数量
		MaxIdles int `yaml:"max_idles" mapstructure:"max_idles"`
		// 最大连接数量
		MaxOpens    int  `yaml:"max_opens" mapstructure:"max_opens"`
		MaxLifeTime int  `yaml:"max_life_time" mapstructure:"max_life_time"`
		LogMode     bool `yaml:"log_mode" mapstructure:"log_mode"`
	}
	DB struct {
		*gorm.DB
		ticker *time.Ticker
	}
)

func NewDB(ops Options) (*DB, error) {
	db, err := gorm.Open(ops.Driver, ops.ConnStr)
	if err != nil {
		return nil, err
	}
	if ops.MaxOpens > 0 {
		db.DB().SetMaxOpenConns(ops.MaxOpens)
	}
	if ops.MaxIdles > 0 {
		db.DB().SetMaxIdleConns(ops.MaxIdles)
	}
	if ops.MaxLifeTime > 0 {
		db.DB().SetConnMaxLifetime(time.Second * time.Duration(ops.MaxLifeTime))
	}

	db.LogMode(ops.LogMode)

	rdb := &DB{DB: db}
	if ops.KeepAlive > 0 {
		rdb.keepAlive(time.Duration(ops.KeepAlive) * time.Second)
	}
	return rdb, nil
}

/*
定时保活
 */
func (f *DB) keepAlive(d time.Duration) {
	f.ticker = time.NewTicker(d)
	go func() {
		for range f.ticker.C {
			if err := f.DB.DB().Ping(); err != nil {
				fmt.Printf("数据库断开连接，%v", err)
			}
		}
	}()
}

/*
Close 关闭数据库连接
 */
func (f *DB) Close() error {
	if f.ticker != nil {
		f.ticker.Stop()
	}
	return f.DB.Close()
}

/*
Resolve 将 SQL 查询结构 序列化到一个 map 中
 */
func Scan(rs *sql.Rows) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for rs.Next() {
		columns, err := rs.ColumnTypes()
		if err != nil {
			return nil, err
		}

		payload := make([]interface{}, 0, len(columns))
		for _, col := range columns {
			t, err := reflectType(col)
			if err != nil {
				return nil, err
			}
			payload = append(payload, reflect.New(t).Interface())
		}

		err = rs.Scan(payload...)
		if err != nil {
			return nil, err
		}

		items := make(map[string]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			items[columns[i].Name()] = reflect.ValueOf(payload[i]).Elem().Interface()
		}

		result = append(result, items)
	}

	return result, nil
}

/*
映射 数据库类型和 golang 类型
暂时不支持 复杂类型 如 enum 等
 */
func reflectType(s *sql.ColumnType) (reflect.Type, error) {
	var vty reflect.Type

	switch s.DatabaseTypeName() {
	case "BIT", "TINYINT", "BOOL":
		vty = reflect.TypeOf(false)
	case "DATE", "DATETIME", "TIME", "TIMESTAMP":
		vty = reflect.TypeOf(time.Now())
	case "TEXT", "BLOB", "LONGBLOB", "LONGTEXT",
		"MEDIUMBLOB", "TINYBLOB", "TINYTEXT",
		"MEDIUMTEXT", "BINARY", "CHAR", "VARBINARY",
		"VARCHAR", "NVARCHAR":
		vty = reflect.TypeOf("")
	case "INT", "MEDIUMINT", "SMALLINT":
		vty = reflect.TypeOf(0)
	case "BIGINT":
		vty = reflect.TypeOf(int64(0))
	case "DOUBLE", "FLOAT", "DECIMAL":
		vty = reflect.TypeOf(float64(0))
	default:
		return nil, errors.New("can't resolve db type")
	}

	if r, ok := s.Nullable(); r && ok {
		vty = reflect.TypeOf(reflect.New(vty).Interface())
	}

	return vty, nil
}
