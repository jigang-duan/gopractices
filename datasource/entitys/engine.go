package entitys

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/jigang-duan/gopractices/bootstrap"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
)

var ORM *xorm.EngineGroup

func NewEngine(driverName string, dataSourceNames ...string) (*xorm.EngineGroup, error) {
	return xorm.NewEngineGroup(driverName, dataSourceNames)
}

func SyncTable(b *bootstrap.Bootstrapper, tables ...interface{}) {
	b.Logger().Info("同步表结构")
	for _, table := range tables  {
		if err := ORM.Sync2(table); err != nil {
			b.Logger().Fatalf("orm未能初始化表: %v", err)
		}
	}
}

func Configure(b *bootstrap.Bootstrapper) {
	driverName := viper.GetString("database.mysql.master.driver")
	masterDataSourceName := viper.GetString("database.mysql.master.dataSource")
	eg, err := NewEngine(driverName, masterDataSourceName)
	if err != nil {
		b.Logger().Fatalf("orm创建数据库引擎出错: %v", err)
	}
	b.Logger().Info("创建XORM引擎")
	ORM = eg

	iris.RegisterOnInterrupt(func() {
		_ = ORM.Close()
		ORM = nil
	})

	eg.ShowSQL(true)
	eg.Logger().SetLevel(core.LOG_INFO)

	needSyncTable := viper.GetBool("database.sync")

	if needSyncTable {
		SyncTable(b, new(User))
	}
}