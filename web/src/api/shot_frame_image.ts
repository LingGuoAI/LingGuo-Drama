import { request } from "@/utils/request";

// @Tags ShotFrameImages
// @Summary 创建分镜图片
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.ShotFrameImages true "创建分镜图片信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /shot_frame_images [post]
export const createShotFrameImages = (data: any) => {
  return request.post({
    url: "/shot_frame_images",
    data,
  });
};

// @Tags ShotFrameImages
// @Summary 删除分镜图片
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path string true "分镜图片ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /shot_frame_images/{id} [delete]
export const deleteShotFrameImages = (id: number | string) => {
  return request.delete({
    url: `/shot_frame_images/${id}`,
  });
};
