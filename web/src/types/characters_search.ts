// 角色搜索相关类型定义

export interface CharactersSearchParams {
  page?: number
  pageSize?: number
  sort?: string
  order?: 'asc' | 'desc'
  id?: number // 主键ID
  projectId?: number // 所属项目ID
  name?: string // 角色名
  roleType?: string // 角色类型: main/supporting/minor
  gender?: string // 性别(需从appearance解析或留空)
}

export interface CharactersSearchForm {
  id?: number
  projectId?: number
  name?: string
  roleType?: string
  gender?: string
}
