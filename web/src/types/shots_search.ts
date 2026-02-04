// 镜头表搜索相关类型定义

export interface ShotsSearchParams {
  page?: number
  pageSize?: number
  sort?: string
  order?: 'asc' | 'desc'
  id?: number // 主键ID
  projectId?: number // 所属项目ID
  scriptId?: number // 所属剧本/分集ID
  sequenceNo?: number // 镜头序号
}

export interface ShotsSearchForm {
  id?: number
  projectId?: number
  scriptId?: number
  sequenceNo?: number
}
