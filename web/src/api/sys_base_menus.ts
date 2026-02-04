import { request } from '@/utils/request'
// @Tags SysBaseMenus
// @Summary 创建系统菜单
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysBaseMenus true "创建系统菜单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sys_base_menuses [post]
export const createSysBaseMenus = (data: any) => {
  return request.post({
    url: '/sys_base_menuses',
    data
  })
}

// @Tags SysBaseMenus
// @Summary 删除系统菜单
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysBaseMenus true "删除系统菜单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sys_base_menuses/{id} [delete]
export const deleteSysBaseMenus = (id: number | string) => {
  return request.delete({
    url: `/sys_base_menuses/${id}`
  })
}

// @Tags SysBaseMenus
// @Summary 更新系统菜单
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysBaseMenus true "更新系统菜单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sys_base_menuses/{id} [put]
export const updateSysBaseMenus = (id: number | string, data: any) => {
  return request.put({
    url: `/sys_base_menuses/${id}`,
    data
  })
}

// @Tags SysBaseMenus
// @Summary 用id查询系统菜单
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysBaseMenus true "用id查询系统菜单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sys_base_menuses/{id} [get]
export const findSysBaseMenus = (id: number | string) => {
  return request.get({
  url: `/sys_base_menuses/${id}`
  })
}
// @Tags SysBaseMenus
// @Summary 分页获取系统菜单列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取系统菜单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys_base_menuses [get]
export const getSysBaseMenusList = (params: any) => {
  return request.get({
    url: '/sys_base_menuses',
    params
  })
}
// @Tags SysBaseMenus
// @Summary 获取系统菜单树结构
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys_base_menuses/getSysBaseMenusTreeList [get]
export const getSysBaseMenusTreeList = () => {
  return request.get({
  url: '/sys_base_menuses/getSysBaseMenusTreeList'
  })
}