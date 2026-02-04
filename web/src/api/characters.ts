import { request } from '@/utils/request'
// @Tags Characters
// @Summary 创建角色
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Characters true "创建角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /characters [post]
export const createCharacters = (data: any) => {
  return request.post({
    url: '/characters',
    data
  })
}

// @Tags Characters
// @Summary 删除角色
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Characters true "删除角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /characters/{id} [delete]
export const deleteCharacters = (id: number | string) => {
  return request.delete({
    url: `/characters/${id}`
  })
}

// @Tags Characters
// @Summary 更新角色
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Characters true "更新角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /characters/{id} [put]
export const updateCharacters = (id: number | string, data: any) => {
  return request.put({
    url: `/characters/${id}`,
    data
  })
}

// @Tags Characters
// @Summary 用id查询角色
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Characters true "用id查询角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /characters/{id} [get]
export const findCharacters = (id: number | string) => {
  return request.get({
  url: `/characters/${id}`
  })
}
// @Tags Characters
// @Summary 分页获取角色列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取角色列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /characters [get]
export const getCharactersList = (params: any) => {
  return request.get({
    url: '/characters',
    params
  })
}
// @Tags Characters
// @Summary 获取短剧项目选择列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态筛选"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /characters/getProjectsSelectList [get]
export const getProjectsSelectList = (params?: any) => {
  return request.get({
  url: '/characters/getProjectsSelectList',
  params
  })
}