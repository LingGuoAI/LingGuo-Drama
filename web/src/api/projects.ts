import { request } from '@/utils/request'
// @Tags Projects
// @Summary 创建短剧项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Projects true "创建短剧项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /projects [post]
export const createProjects = (data: any) => {
  return request.post({
    url: '/projects',
    data
  })
}

// @Tags Projects
// @Summary 删除短剧项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Projects true "删除短剧项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /projects/{id} [delete]
export const deleteProjects = (id: number | string) => {
  return request.delete({
    url: `/projects/${id}`
  })
}

// @Tags Projects
// @Summary 更新短剧项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Projects true "更新短剧项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /projects/{id} [put]
export const updateProjects = (id: number | string, data: any) => {
  return request.put({
    url: `/projects/${id}`,
    data
  })
}

// @Tags Projects
// @Summary 用id查询短剧项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Projects true "用id查询短剧项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /projects/{id} [get]
export const findProjects = (id: number | string) => {
  return request.get({
  url: `/projects/${id}`
  })
}
// @Tags Projects
// @Summary 分页获取短剧项目列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取短剧项目列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /projects [get]
export const getProjectsList = (params: any) => {
  return request.get({
    url: '/projects',
    params
  })
}