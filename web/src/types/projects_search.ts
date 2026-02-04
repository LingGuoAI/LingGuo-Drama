// 短剧项目搜索相关类型定义

export interface ProjectsSearchParams {
  page?: number
  pageSize?: number
  sort?: string
  order?: 'asc' | 'desc'
  id?: number // 主键ID
  adminId?: number // 归属用户ID(默认1)
  serialNo?: string // 业务流水号
  title?: string // 项目名称/短剧标题
  status?: any // 状态
}

export interface ProjectsSearchForm {
  id?: number
  adminId?: number
  serialNo?: string
  title?: string
  status?: any
}
