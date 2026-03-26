//go:build driver_mssql || !driver_custom

package system

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/gookit/color"
	"github.com/icosmos-space/iadmin/server/config"
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/model/system/request"
	"github.com/icosmos-space/iadmin/server/utils"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type MssqlInitHandler struct{}

func NewMssqlInitHandler() *MssqlInitHandler {
	return &MssqlInitHandler{}
}

// WriteConfig mssql回写配置
func (h MssqlInitHandler) WriteConfig(ctx context.Context) error {
	c, ok := ctx.Value("config").(config.Mssql)
	if !ok {
		return errors.New("mssql config invalid")
	}
	global.IADMIN_CONFIG.System.DbType = "mssql"
	global.IADMIN_CONFIG.Mssql = c
	global.IADMIN_CONFIG.JWT.SigningKey = uuid.New().String()
	cs := utils.StructToMap(global.IADMIN_CONFIG)
	for k, v := range cs {
		global.IADMIN_VP.Set(k, v)
	}
	global.IADMIN_ACTIVE_DBNAME = &c.Dbname
	return global.IADMIN_VP.WriteConfig()
}

// EnsureDB 创建数据库并初始化 mssql
func (h MssqlInitHandler) EnsureDB(ctx context.Context, conf *request.InitDB) (next context.Context, err error) {
	if s, ok := ctx.Value("dbtype").(string); !ok || s != "mssql" {
		return ctx, ErrDBTypeMismatch
	}

	c := conf.ToMssqlConfig()
	next = context.WithValue(ctx, "config", c)
	if c.Dbname == "" {
		return ctx, nil
	} // 如果没有数据库名, 则跳出初始化数据

	masterDsn := conf.MssqlMasterDsn()
	dbnameForID := strings.ReplaceAll(c.Dbname, "'", "''")
	dbnameForCreate := strings.ReplaceAll(c.Dbname, "]", "]]")
	createSql := fmt.Sprintf("IF DB_ID(N'%s') IS NULL CREATE DATABASE [%s];", dbnameForID, dbnameForCreate)
	if err = createDatabase(masterDsn, "sqlserver", createSql); err != nil {
		return nil, err
	}

	dsn := conf.MssqlEmptyDsn()

	mssqlConfig := sqlserver.Config{
		DSN:               dsn, // DSN data source name
		DefaultStringSize: 191, // string 类型字段的默认长度
	}

	var db *gorm.DB

	if db, err = gorm.Open(sqlserver.New(mssqlConfig), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
		return nil, err
	}

	global.IADMIN_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	next = context.WithValue(next, "db", db)
	return next, err
}

func (h MssqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	return createTables(ctx, inits)
}

func (h MssqlInitHandler) InitData(ctx context.Context, inits initSlice) error {
	next, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, init := range inits {
		if init.DataInserted(next) {
			color.Info.Printf(InitDataExist, Mssql, init.InitializerName())
			continue
		}
		if n, err := init.InitializeData(next); err != nil {
			color.Info.Printf(InitDataFailed, Mssql, init.InitializerName(), err)
			return err
		} else {
			next = n
			color.Info.Printf(InitDataSuccess, Mssql, init.InitializerName())
		}
	}
	color.Info.Printf(InitSuccess, Mssql)
	return nil
}
