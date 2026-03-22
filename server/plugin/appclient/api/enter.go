package api

type ApiGroup struct {
	AppAuth    appAuth
	AppProfile appProfile
}

var ApiGroupApp = new(ApiGroup)
