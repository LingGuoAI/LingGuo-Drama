// 剧本搜索相关类型定义

export interface ScriptsSearchParams {
  page?: number
  pageSize?: number
  sort?: string
  order?: 'asc' | 'desc'
  id?: number // 主键ID
  projectId?: number // 所属项目ID
  title?: string // 分集标题
  episodeNo?: number // 第几集
  isLocked?: any // 是否定稿
}

export interface ScriptsSearchForm {
  id?: number
  projectId?: number
  title?: string
  episodeNo?: number
  isLocked?: any
}
