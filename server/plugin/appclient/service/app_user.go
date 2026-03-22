package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/appclient/model"
	appReq "github.com/icosmos-space/iadmin/server/plugin/appclient/model/request"
	appRes "github.com/icosmos-space/iadmin/server/plugin/appclient/model/response"
	"github.com/icosmos-space/iadmin/server/utils"
	"gorm.io/gorm"
)

type AppUserService struct{}

var AppUserServiceApp = new(AppUserService)

func toPublic(u *model.AppUser) appRes.AppUserPublic {
	return appRes.AppUserPublic{
		ID:       u.ID,
		UUID:     u.UUID,
		Username: u.Username,
		Nickname: u.Nickname,
		Phone:    u.Phone,
	}
}

func (s *AppUserService) issueToken(u *model.AppUser) (string, error) {
	j := utils.NewAppJWT()
	cl := j.CreateClaims(u.ID, u.Username, u.UUID)
	return j.CreateToken(cl)
}

// Register 注册
func (s *AppUserService) Register(req appReq.AppRegister) (*appRes.AppLoginResponse, error) {
	var n int64
	global.IADMIN_DB.Model(&model.AppUser{}).Where("username = ?", req.Username).Count(&n)
	if n > 0 {
		return nil, errors.New("用户名已存在")
	}
	u := model.AppUser{
		UUID:     uuid.New(),
		Username: req.Username,
		Password: utils.BcryptHash(req.Password),
		Nickname: req.Nickname,
		Enable:   1,
	}
	if u.Nickname == "" {
		u.Nickname = req.Username
	}
	if err := global.IADMIN_DB.Create(&u).Error; err != nil {
		return nil, err
	}
	token, err := s.issueToken(&u)
	if err != nil {
		return nil, err
	}
	return &appRes.AppLoginResponse{Token: token, User: toPublic(&u)}, nil
}

// Login 登录
func (s *AppUserService) Login(req appReq.AppLogin) (*appRes.AppLoginResponse, error) {
	var u model.AppUser
	err := global.IADMIN_DB.Where("username = ?", req.Username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户不存在或密码错误")
	}
	if err != nil {
		return nil, err
	}
	if u.Enable != 1 {
		return nil, errors.New("账号已禁用")
	}
	if ok := utils.BcryptCheck(req.Password, u.Password); !ok {
		return nil, errors.New("用户不存在或密码错误")
	}
	token, err := s.issueToken(&u)
	if err != nil {
		return nil, err
	}
	return &appRes.AppLoginResponse{Token: token, User: toPublic(&u)}, nil
}

// GetByID 根据主键查询（脱敏）
func (s *AppUserService) GetByID(id uint) (*appRes.AppUserPublic, error) {
	var u model.AppUser
	if err := global.IADMIN_DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	if u.Enable != 1 {
		return nil, errors.New("账号已禁用")
	}
	p := toPublic(&u)
	return &p, nil
}
