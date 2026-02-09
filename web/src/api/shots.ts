import { request } from "@/utils/request";
// @Tags Shots
// @Summary 创建镜头表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Shots true "创建镜头表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /shots [post]
export const createShots = (data: any) => {
  return request.post({
    url: "/shots",
    data,
  });
};

// @Tags Shots
// @Summary 删除镜头表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Shots true "删除镜头表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /shots/{id} [delete]
export const deleteShots = (id: number | string) => {
  return request.delete({
    url: `/shots/${id}`,
  });
};

// @Tags Shots
// @Summary 更新镜头表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Shots true "更新镜头表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /shots/{id} [put]
export const updateShots = (id: number | string, data: any) => {
  return request.put({
    url: `/shots/${id}`,
    data,
  });
};

// @Tags Shots
// @Summary 用id查询镜头表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Shots true "用id查询镜头表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /shots/{id} [get]
export const findShots = (id: number | string) => {
  return request.get({
    url: `/shots/${id}`,
  });
};
// @Tags Shots
// @Summary 分页获取镜头表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取镜头表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /shots [get]
export const getShotsList = (params: any) => {
  return request.get({
    url: "/shots",
    params,
  });
};
// @Tags Shots
// @Summary 获取短剧项目选择列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态筛选"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /shots/getProjectsSelectList [get]
export const getProjectsSelectList = (params?: any) => {
  return request.get({
    url: "/shots/getProjectsSelectList",
    params,
  });
};
// @Tags Shots
// @Summary 获取剧本选择列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态筛选"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /shots/getScriptsSelectList [get]
export const getScriptsSelectList = (params?: any) => {
  return request.get({
    url: "/shots/getScriptsSelectList",
    params,
  });
};
