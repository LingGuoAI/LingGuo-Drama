import { request } from '@/utils/request'
// @Tags Admins
// @Summary 创建系统管理员
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Admins true "创建系统管理员"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /admins [post]
export const createAdmins = (data: any) => {
  return request.post({
    url: '/admins',
    data
  })
}

// @Tags Admins
// @Summary 删除系统管理员
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Admins true "删除系统管理员"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /admins/{id} [delete]
export const deleteAdmins = (id: number | string) => {
  return request.delete({
    url: `/admins/${id}`
  })
}

// @Tags Admins
// @Summary 更新系统管理员
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Admins true "更新系统管理员"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /admins/{id} [put]
export const updateAdmins = (id: number | string, data: any) => {
  return request.put({
    url: `/admins/${id}`,
    data
  })
}

// @Tags Admins
// @Summary 用id查询系统管理员
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Admins true "用id查询系统管理员"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /admins/{id} [get]
export const findAdmins = (id: number | string) => {
  return request.get({
  url: `/admins/${id}`
  })
}
// @Tags Admins
// @Summary 分页获取系统管理员列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取系统管理员列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /admins [get]
export const getAdminsList = (params: any) => {
  return request.get({
    url: '/admins',
    params
  })
}