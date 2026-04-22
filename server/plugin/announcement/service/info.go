package service

import (
	"errors"
	"github.com/icosmos-space/iadmin/server/global"
	commonReq "github.com/icosmos-space/iadmin/server/model/common/request"
	"github.com/icosmos-space/iadmin/server/plugin/announcement/model"
	"github.com/icosmos-space/iadmin/server/plugin/announcement/model/request"
)

var Info = new(info)

type info struct{}

// CreateInfo 创建公告记录
// Author [piexlmax](https://github.com/piexlmax)
func (s *info) CreateInfo(info *model.Info) (err error) {
	err = global.IADMIN_DB.Create(info).Error
	return err
}

// DeleteInfo 删除公告记录
// Author [piexlmax](https://github.com/piexlmax)
func (s *info) DeleteInfo(ID string) (err error) {
	err = global.IADMIN_DB.Delete(&model.Info{}, "id = ?", ID).Error
	return err
}

// DeleteInfoByIds 批量删除公告记录
// Author [piexlmax](https://github.com/piexlmax)
func (s *info) DeleteInfoByIds(IDs []string) (err error) {
	err = global.IADMIN_DB.Delete(&[]model.Info{}, "id in ?", IDs).Error
	return err
}

// UpdateInfo 更新公告记录
// Author [piexlmax](https://github.com/piexlmax)
func (s *info) UpdateInfo(info model.Info) (err error) {
	err = global.IADMIN_DB.Model(&model.Info{}).Where("id = ?", info.ID).Updates(&info).Error
	return err
}

// GetInfo 根据ID获取公告记录
// Author [piexlmax](https://github.com/piexlmax)
func (s *info) GetInfo(ID string) (info model.Info, err error) {
	err = global.IADMIN_DB.Where("id = ?", ID).First(&info).Error
	return
}

// GetInfoInfoList 分页获取公告记录
// Author [piexlmax](https://github.com/piexlmax)
func (s *info) GetInfoInfoList(info request.InfoSearch) (list []model.Info, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.IADMIN_DB.Model(&model.Info{})
	var infos []model.Info
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&infos).Error
	return infos, total, err
}
func (s *info) GetInfoDataSource(info commonReq.DataSourceQuery) (res commonReq.DataSourceResult, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	switch {
	case info.PageSize > 100:
		info.PageSize = 100
	case info.PageSize <= 0:
		info.PageSize = 20
	}
	offset := (info.Page - 1) * info.PageSize
	res.Page = info.Page
	res.PageSize = info.PageSize
	res.List = make([]map[string]any, 0)

	switch info.Field {
	case "userID":
		db := global.IADMIN_DB.Table("sys_users")
		if info.Keyword != "" {
			db = db.Where("nick_name LIKE ?", "%"+info.Keyword+"%")
		}
		if err = db.Count(&res.Total).Error; err != nil {
			return
		}
		err = db.Select("nick_name as label,id as value").Order("id desc").Offset(offset).Limit(info.PageSize).Scan(&res.List).Error
	default:
		err = errors.New("invalid datasource field")
		return
	}

	res.HasMore = int64(offset+len(res.List)) < res.Total
	return
}
