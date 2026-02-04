// 剧本相关类型定义

import { UploadFile } from './common'

export interface Scripts {
  id?: number // 主键ID
  projectId: number // 所属项目ID, 外键约束(project_id) -> projects(id)
  title?: string // 分集标题
  content?: string // 剧本正文
  outline?: string // 大纲/简介
  episodeNo?: number // 第几集
  isLocked?: any // 是否定稿 0-否 1-是
  createdAt?: string // 添加时间
  updatedAt?: string // 修改时间
  deletedAt?: string // 删除时间
}

export interface CreateScriptsParams extends Omit<Scripts, 'id' | 'createdAt' | 'updatedAt'> {}

export interface UpdateScriptsParams extends Partial<Scripts> {
  id: number
}

export interface ScriptsSearchParams {
  page?: number
  pageSize?: number
  id?: number
  projectId?: number
  title?: string
  episodeNo?: number
  isLocked?: any
}
