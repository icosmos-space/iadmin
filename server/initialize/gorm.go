package initialize

import (
	"os"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/model/example"
	"github.com/icosmos-space/iadmin/server/model/system"
	appclientmodel "github.com/icosmos-space/iadmin/server/plugin/appclient/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.IADMIN_CONFIG.System.DbType {
	case "mysql":
		global.IADMIN_ACTIVE_DBNAME = &global.IADMIN_CONFIG.Mysql.Dbname
		return GormMysql()
	case "pgsql":
		global.IADMIN_ACTIVE_DBNAME = &global.IADMIN_CONFIG.Pgsql.Dbname
		return GormPgSql()
	case "oracle":
		global.IADMIN_ACTIVE_DBNAME = &global.IADMIN_CONFIG.Oracle.Dbname
		return GormOracle()
	case "mssql":
		global.IADMIN_ACTIVE_DBNAME = &global.IADMIN_CONFIG.Mssql.Dbname
		return GormMssql()
	case "sqlite":
		global.IADMIN_ACTIVE_DBNAME = &global.IADMIN_CONFIG.Sqlite.Dbname
		return GormSqlite()
	default:
		global.IADMIN_ACTIVE_DBNAME = &global.IADMIN_CONFIG.Mysql.Dbname
		return GormMysql()
	}
}

func RegisterTables() {
	if global.IADMIN_CONFIG.System.DisableAutoMigrate {
		global.IADMIN_LOG.Info("auto-migrate is disabled, skipping table registration")
		return
	}

	db := global.IADMIN_DB
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysIgnoreApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCodePackage{},
		system.SysExportTemplate{},
		system.Condition{},
		system.JoinTemplate{},
		system.SysParams{},
		system.SysVersion{},
		system.SysError{},
		system.SysApiToken{},
		system.SysLoginLog{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},
		example.ExaAttachmentCategory{},

		appclientmodel.AppUser{},
	)
	if err != nil {
		global.IADMIN_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = bizModel()

	if err != nil {
		global.IADMIN_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.IADMIN_LOG.Info("register table success")
}
