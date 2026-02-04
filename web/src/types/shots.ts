// 镜头表相关类型定义

import { UploadFile } from './common'

export interface Shots {
  id?: number // 主键ID
  projectId: number // 所属项目ID, 外键约束(project_id) -> projects(id)
  scriptId: number // 所属剧本/分集ID, 外键约束(script_id) -> scripts(id)
  sequenceNo: number // 镜头序号
  shotType?: string // 景别: 全景/特写/中景
  cameraMovement?: string // 运镜: 推/拉/摇/移
  angle?: string // 视角: 俯视/平视
  dialogue?: string // 台词/旁白
  visualDesc?: string // 画面描述
  atmosphere?: string // 氛围/环境描述
  imagePrompt?: string // 绘画Prompt
  videoPrompt?: string // 视频生成Prompt
  audioPrompt?: string // 音效/BGM提示词
  imageUrl?: string // 分镜图
  videoUrl?: string // 最终视频片段
  audioUrl?: string // 配音/音效
  durationMs?: number // 时长(毫秒, 原duration*1000)
  status?: any // 状态 0-Pending 1-Done 2-Fail
  createdAt?: string // 添加时间
  updatedAt?: string // 修改时间
  deletedAt?: string // 删除时间
}

export interface CreateShotsParams extends Omit<Shots, 'id' | 'createdAt' | 'updatedAt'> {}

export interface UpdateShotsParams extends Partial<Shots> {
  id: number
}

export interface ShotsSearchParams {
  page?: number
  pageSize?: number
  id?: number
  projectId?: number
  scriptId?: number
  sequenceNo?: number
}
