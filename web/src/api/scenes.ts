import { request } from "@/utils/request";
// @Tags Scenes
// @Summary 创建场景项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Scenes true "创建场景项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /scenes [post]
export const createScenes = (data: any) => {
  return request.post({
    url: "/scenes",
    data,
  });
};

// @Tags Scenes
// @Summary 删除场景项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Scenes true "删除场景项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /scenes/{id} [delete]
export const deleteScenes = (id: number | string) => {
  return request.delete({
    url: `/scenes/${id}`,
  });
};

// @Tags Scenes
// @Summary 更新场景项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Projects true "更新场景项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /scenes/{id} [put]
export const updateScenes = (id: number | string, data: any) => {
  return request.put({
    url: `/scenes/${id}`,
    data,
  });
};

// @Tags Scenes
// @Summary 用id查询场景项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Scenes true "用id查询场景项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /scenes/{id} [get]
export const findScenes = (id: number | string) => {
  return request.get({
    url: `/scenes/${id}`,
  });
};
// @Tags Scenes
// @Summary 分页获取场景项目列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取场景项目列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /scenes [get]
export const getScenesList = (params: any) => {
  return request.get({
    url: "/scenes",
    params,
  });
};
