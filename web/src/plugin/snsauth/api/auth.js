import service from '@/utils/request'

export const getSnsProviderList = () => {
  return service({
    url: '/snsAuth/getProviderList',
    method: 'get'
  })
}

export const getEnabledSnsProviderList = () => {
  return service({
    url: '/snsAuth/getEnabledProviderList',
    method: 'get'
  })
}

export const updateSnsProviderConfig = (data) => {
  return service({
    url: '/snsAuth/updateProviderConfig',
    method: 'put',
    data
  })
}

export const getSnsLoginURL = (params) => {
  return service({
    url: '/snsAuth/getLoginURL',
    method: 'get',
    params
  })
}

export const getSnsBindURL = (params) => {
  return service({
    url: '/snsAuth/getBindURL',
    method: 'get',
    params
  })
}

export const getMySnsBindings = () => {
  return service({
    url: '/snsAuth/getMyBindings',
    method: 'get'
  })
}

export const unbindSns = (params) => {
  return service({
    url: '/snsAuth/unbind',
    method: 'delete',
    params
  })
}

export const telegramLogin = (data) => {
  return service({
    url: '/snsAuth/telegramLogin',
    method: 'post',
    data
  })
}

export const telegramBind = (data) => {
  return service({
    url: '/snsAuth/telegramBind',
    method: 'post',
    data
  })
}
