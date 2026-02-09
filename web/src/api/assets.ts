// src/api/assets.ts
import { request } from "@/utils/request";

// 获取素材列表
export const getAssetsList = (params: any) => {
  return request.get({
    url: "/assets",
    params,
  });
};

// 创建素材 (将分镜视频转为素材)
export const createAsset = (data: any) => {
  return request.post({
    url: "/assets",
    data,
  });
};

// 删除素材
export const deleteAsset = (id: number | string) => {
  return request.delete({
    url: `/assets/${id}`,
  });
};
