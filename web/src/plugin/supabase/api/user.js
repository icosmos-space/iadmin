import service from '@/utils/request'

export const getSupabaseUserList = (params) => {
  return service({
    url: '/supabaseUser/getUserList',
    method: 'get',
    params
  })
}

export const updateSupabaseUserPassword = (data) => {
  return service({
    url: '/supabaseUser/updatePassword',
    method: 'put',
    data
  })
}
