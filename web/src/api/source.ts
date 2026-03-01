import { request } from "@/utils/request";

// @Tags Source
// @Summary 创建素材
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Source true "创建素材信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /source [post]
export const createSource = (data: any) => {
  return request.post({
    url: "/source",
    data,
  });
};

// @Tags Source
// @Summary 删除素材
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path string true "素材ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /source/{id} [delete]
export const deleteSource = (id: number | string) => {
  return request.delete({
    url: `/source/${id}`,
  });
};
