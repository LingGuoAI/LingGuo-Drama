import { request } from '@/utils/request'
// @Tags Scripts
// @Summary 创建剧本
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Scripts true "创建剧本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /scripts [post]
export const createScripts = (data: any) => {
  return request.post({
    url: '/scripts',
    data
  })
}

// @Tags Scripts
// @Summary 删除剧本
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Scripts true "删除剧本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /scripts/{id} [delete]
export const deleteScripts = (id: number | string) => {
  return request.delete({
    url: `/scripts/${id}`
  })
}

// @Tags Scripts
// @Summary 更新剧本
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Scripts true "更新剧本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /scripts/{id} [put]
export const updateScripts = (id: number | string, data: any) => {
  return request.put({
    url: `/scripts/${id}`,
    data
  })
}

// @Tags Scripts
// @Summary 用id查询剧本
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Scripts true "用id查询剧本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /scripts/{id} [get]
export const findScripts = (id: number | string) => {
  return request.get({
  url: `/scripts/${id}`
  })
}
// @Tags Scripts
// @Summary 分页获取剧本列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取剧本列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /scripts [get]
export const getScriptsList = (params: any) => {
  return request.get({
    url: '/scripts',
    params
  })
}
// @Tags Scripts
// @Summary 获取短剧项目选择列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态筛选"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /scripts/getProjectsSelectList [get]
export const getProjectsSelectList = (params?: any) => {
  return request.get({
  url: '/scripts/getProjectsSelectList',
  params
  })
}