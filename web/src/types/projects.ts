// 短剧项目相关类型定义

import { UploadFile } from './common'

export interface Projects {
  id?: number // 主键ID
  adminId: number // 归属用户ID(默认1)
  serialNo?: string // 业务流水号
  title: string // 项目名称/短剧标题
  description?: string // 项目简介
  status: any // 状态 0-Draft 1-Generating 2-Completed
  image?: string // 封面图
  totalDuration?: number // 总时长(秒)
  settings?: any // 生成配置快照
  createdAt?: string // 添加时间
  updatedAt?: string // 修改时间
  deletedAt?: string // 删除时间
}

export interface CreateProjectsParams extends Omit<Projects, 'id' | 'createdAt' | 'updatedAt'> {}

export interface UpdateProjectsParams extends Partial<Projects> {
  id: number
}

export interface ProjectsSearchParams {
  page?: number
  pageSize?: number
  id?: number
  adminId?: number
  serialNo?: string
  title?: string
  status?: any
}
