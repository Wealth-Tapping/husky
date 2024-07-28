package husky

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DatabaseTxCall func(tx *gorm.DB) error

type DatabaseTxCalls interface {
	AddCall(call DatabaseTxCall)
	Run(db *gorm.DB) error
}

type _DatabaseTxCalls []DatabaseTxCall

func (calls *_DatabaseTxCalls) AddCall(call DatabaseTxCall) {
	*calls = append(*calls, call)
}

func (calls *_DatabaseTxCalls) Run(db *gorm.DB) error {
	if len(*calls) == 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, call := range *calls {
			if err := call(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

func NewDatabaseTxCalls() DatabaseTxCalls {
	return new(_DatabaseTxCalls)
}

type DatabaseConfig struct {
	Url         string `toml:"url"`
	MaxIdleConn *int   `toml:"maxIdleConn"`
	MaxOpenConn *int   `toml:"maxOpenConn"`
	MaxLifetime *int64 `toml:"maxLifetime"`
	ShowSql     bool   `toml:"showSql"`
	SlowTimeMs  *int64 `toml:"slowTime"`
}

var _DatabaseIns map[string]*gorm.DB

func init() {
	_DatabaseIns = make(map[string]*gorm.DB)
}

type _DatabaseLog struct{}

func (*_DatabaseLog) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func InitDatabase(config *DatabaseConfig, key ...string) {
	dbConfig := postgres.Config{
		DSN:                  config.Url,
		PreferSimpleProtocol: true,
	}

	dbLog := logger.New(
		&_DatabaseLog{},
		logger.Config{
			SlowThreshold:             time.Millisecond * time.Duration(NilDefault(config.SlowTimeMs, 200)),
			LogLevel:                  If(config.ShowSql, logger.Info, logger.Error),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	db, err := gorm.Open(postgres.New(dbConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		CreateBatchSize:        1000,
		Logger:                 dbLog,
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
		AllowGlobalUpdate:      true,
	})

	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if config.MaxIdleConn != nil {
		sqlDB.SetMaxIdleConns(*config.MaxIdleConn)
	}
	if config.MaxOpenConn != nil {
		sqlDB.SetMaxOpenConns(*config.MaxOpenConn)
	}
	if config.MaxLifetime != nil {
		sqlDB.SetConnMaxIdleTime(time.Duration(*config.MaxLifetime) * time.Second)
	}

	if len(key) == 0 {
		_DatabaseIns[""] = db
	} else {
		_DatabaseIns[key[0]] = db
	}
}

func Database(key ...string) *gorm.DB {
	if len(key) == 0 {
		return _DatabaseIns[""]
	} else {
		return _DatabaseIns[key[0]]
	}
}
