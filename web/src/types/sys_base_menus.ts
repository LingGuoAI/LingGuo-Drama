// 系统菜单相关类型定义

import { UploadFile } from './common'

export interface SysBaseMenus {
  id?: number // 菜单ID
  parentId?: number // 父菜单ID
  path: string // 路由path
  name: string // 路由name
  hidden?: any // 是否在列表隐藏
  component?: string // 对应前端文件路径
  sort?: number // 排序标记
  title: string // 菜单名
  icon?: string // 菜单图标
  createdAt: string // 创建时间
  updatedAt: string // 更新时间
  deletedAt?: string // 删除时间
}

export interface CreateSysBaseMenusParams extends Omit<SysBaseMenus, 'id' | 'createdAt' | 'updatedAt'> {}

export interface UpdateSysBaseMenusParams extends Partial<SysBaseMenus> {
  id: number
}

export interface SysBaseMenusSearchParams {
  page?: number
  pageSize?: number
}
