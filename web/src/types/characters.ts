// 角色相关类型定义

import { UploadFile } from './common'

export interface Characters {
  id?: number // 主键ID
  projectId: number // 所属项目ID, 外键约束(project_id) -> projects(id)
  name: string // 角色名
  roleType?: string // 角色类型: main/supporting/minor
  gender?: string // 性别(需从appearance解析或留空)
  ageGroup?: string // 年龄段
  personality?: string // 性格描述
  appearanceDesc?: string // 外貌长文本描述(原appearance)
  visualPrompt?: string // AI绘画专用Prompt(从appearance提取)
  avatarUrl?: string // 头像/定妆照
  voiceId?: string // TTS音色ID
  createdAt?: string // 添加时间
  updatedAt?: string // 修改时间
  deletedAt?: string // 删除时间
}

export interface CreateCharactersParams extends Omit<Characters, 'id' | 'createdAt' | 'updatedAt'> {}

export interface UpdateCharactersParams extends Partial<Characters> {
  id: number
}

export interface CharactersSearchParams {
  page?: number
  pageSize?: number
  id?: number
  projectId?: number
  name?: string
  roleType?: string
  gender?: string
}
