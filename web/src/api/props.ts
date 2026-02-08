import { request } from "@/utils/request";
// @Tags Props
// @Summary 创建道具项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Scenes true "创建道具项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /props [post]
export const createProps = (data: any) => {
  return request.post({
    url: "/props",
    data,
  });
};

// @Tags Props
// @Summary 删除道具项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Props true "删除道具项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /props/{id} [delete]
export const deleteProps = (id: number | string) => {
  return request.delete({
    url: `/props/${id}`,
  });
};

// @Tags Props
// @Summary 更新道具项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Projects true "更新道具项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /props/{id} [put]
export const updateProps = (id: number | string, data: any) => {
  return request.put({
    url: `/props/${id}`,
    data,
  });
};

// @Tags Props
// @Summary 用id查询道具项目
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Scenes true "用id查询道具项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /props/{id} [get]
export const findProps = (id: number | string) => {
  return request.get({
    url: `/props/${id}`,
  });
};
// @Tags Props
// @Summary 分页获取道具项目列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取道具项目列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /props [get]
export const getPropsList = (params: any) => {
  return request.get({
    url: "/props",
    params,
  });
};
