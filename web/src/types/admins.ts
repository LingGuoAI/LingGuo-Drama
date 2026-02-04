// 系统管理员相关类型定义

import { UploadFile } from './common'

export interface Admins {
  id?: number // 主键ID
  username: string // 用户名
  mobile: string // 手机号
  password: string // 密码
  email?: string // 邮箱
  authorityId?: number // 用户角色ID
  createdAt: string // 创建时间
  updatedAt: string // 更新时间
  deletedAt?: string // 删除时间
}

export interface CreateAdminsParams extends Omit<Admins, 'id' | 'createdAt' | 'updatedAt'> {}

export interface UpdateAdminsParams extends Partial<Admins> {
  id: number
}

export interface AdminsSearchParams {
  page?: number
  pageSize?: number
}
