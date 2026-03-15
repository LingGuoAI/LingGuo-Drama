import { request } from "@/utils/request";

// @Tags AiConfig
// @Summary 创建AI服务配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AiConfig true "创建AI服务配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /ai-config [post]
export const createAiConfig = (data: any) => {
  return request.post({
    url: "/ai-config",
    data,
  });
};

// @Tags AiConfig
// @Summary 删除AI服务配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path int true "删除AI服务配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ai-config/{id} [delete]
export const deleteAiConfig = (id: number | string) => {
  return request.delete({
    url: `/ai-config/${id}`,
  });
};

// @Tags AiConfig
// @Summary 更新AI服务配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path int true "更新AI服务配置的ID"
// @Param data body model.AiConfig true "更新AI服务配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ai-config/{id} [put]
export const updateAiConfig = (id: number | string, data: any) => {
  return request.put({
    url: `/ai-config/${id}`,
    data,
  });
};

// @Tags AiConfig
// @Summary 用id查询AI服务配置详情
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path int true "用id查询AI服务配置详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /ai-config/{id} [get]
export const findAiConfig = (id: number | string) => {
  return request.get({
    url: `/ai-config/${id}`,
  });
};

// @Tags AiConfig
// @Summary 分页获取AI服务配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取AI服务配置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ai-config [get]
export const getAiConfigList = (params: any) => {
  return request.get({
    url: "/ai-config",
    params,
  });
};

// @Tags AiConfig
// @Summary 测试AI服务连通性 (如果你后端加了这个接口的话，可以留着；如果没加可以不放到里面)
// @Router /ai-config/test [post]
export const testAiConfigConnection = (data: any) => {
  return request.post({
    url: "/ai-config/test",
    data,
  });
};

// ==========================================
// 💡 AI 连通性测试专属接口 (挂载在 tasks 路由下)
// ==========================================

// @Tags AiConfig
// @Summary 测试文本大模型配置连通性
// @Router /tasks/testTextConfig [post]
export const testTextConfig = (data: any) => {
  return request.post({
    url: "/tasks/testTextConfig",
    data,
    timeout: 60000,
  });
};

// @Tags AiConfig
// @Summary 测试生图大模型配置连通性
// @Router /tasks/testImageConfig [post]
export const testImageConfig = (data: any) => {
  return request.post({
    url: "/tasks/testImageConfig",
    data,
    timeout: 60000,
  });
};

// @Tags AiConfig
// @Summary 测试视频大模型配置连通性
// @Router /tasks/testVideoConfig [post]
export const testVideoConfig = (data: any) => {
  return request.post({
    url: "/tasks/testVideoConfig",
    data,
    timeout: 60000,
  });
};
