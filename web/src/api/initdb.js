import service from '@/utils/request'
// @Tags InitDB
// @Summary 初始化用户数据库
// @Produce  application/json
// @Param data body request.InitDB true "初始化数据库参数"
// @Success 200 {string} string "{"code":0,"data":{},"msg":"自动创建数据库成功"}"
// @Router /init/initdb [post]
export const initDB = (data) => {
  return service({
    url: '/init/initdb',
    method: 'post',
    data,
    donNotShowLoading: true
  })
}

// @Tags CheckDB
// @Summary 初始化用户数据库
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"探测完成"}"
// @Router /init/checkdb [post]
export const checkDB = () => {
  return service({
    url: '/init/checkdb',
    method: 'post'
  })
}

// 与后端 driver_custom / driver_* 一致的可选数据库类型
export const getInitDbTypes = () => {
  return service({
    url: '/init/db-types',
    method: 'get',
    donNotShowLoading: true
  })
}
