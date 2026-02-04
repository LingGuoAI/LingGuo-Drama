import { request } from '@/utils/request'

// 统计卡片数据结构
export interface StatisticsCard {
  title: string        // 标题
  number: string       // 数量
  upTrend?: string     // 上升趋势
  downTrend?: string   // 下降趋势
  leftType: string     // 左侧图标类型
}

// 表趋势数据结构
export interface TableTrendData {
  date: string         // 日期
  count: number        // 数量
}

// 详细统计数据结构
export interface DetailStatistics {
  
  characters: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  }
  
  projects: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  }
  
  scripts: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  }
  
  shots: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  }
  
  admins: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  }
  
  sys_base_menus: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  }
  
}

// 统计API请求参数
export interface StatisticsParams {
  // 基础统计不需要参数
}

// 表趋势请求参数
export interface TableTrendParams {
  table_name: string   // 表名
}

/**
 * 获取统计数据
 * @returns 统计卡片数据数组
 */
export const getStatisticsData = (): Promise<{
  code: number
  data: StatisticsCard[]
  msg: string
}> => {
  return request.get({
    url: '/statistics/statistics'
  })
}

/**
 * 获取表趋势数据
 * @param params 请求参数
 * @returns 趋势数据数组
 */
export const getTableTrend = (params: TableTrendParams): Promise<{
  code: number
  data: TableTrendData[]
  msg: string
}> => {
  return request.get({
    url: '/statistics/trend',
    params
  })
}

/**
 * 获取详细统计信息
 * @returns 详细统计数据
 */
export const getDetailStatistics = (): Promise<{
  code: number
  data: DetailStatistics
  msg: string
}> => {
  return request.get({
    url: '/statistics/detail'
  })
}


/**
 * 获取角色趋势数据
 * @returns 角色7天趋势数据
 */
export const getCharactersTrend = (): Promise<{
  code: number
  data: TableTrendData[]
  msg: string
}> => {
  return request.get({
    url: '/statistics/trend',
    params: { table_name: 'characters' }
  })
}

/**
 * 获取短剧项目趋势数据
 * @returns 短剧项目7天趋势数据
 */
export const getProjectsTrend = (): Promise<{
  code: number
  data: TableTrendData[]
  msg: string
}> => {
  return request.get({
    url: '/statistics/trend',
    params: { table_name: 'projects' }
  })
}

/**
 * 获取剧本趋势数据
 * @returns 剧本7天趋势数据
 */
export const getScriptsTrend = (): Promise<{
  code: number
  data: TableTrendData[]
  msg: string
}> => {
  return request.get({
    url: '/statistics/trend',
    params: { table_name: 'scripts' }
  })
}

/**
 * 获取镜头表趋势数据
 * @returns 镜头表7天趋势数据
 */
export const getShotsTrend = (): Promise<{
  code: number
  data: TableTrendData[]
  msg: string
}> => {
  return request.get({
    url: '/statistics/trend',
    params: { table_name: 'shots' }
  })
}

/**
 * 获取系统管理员趋势数据
 * @returns 系统管理员7天趋势数据
 */
export const getAdminsTrend = (): Promise<{
  code: number
  data: TableTrendData[]
  msg: string
}> => {
  return request.get({
    url: '/statistics/trend',
    params: { table_name: 'admins' }
  })
}

/**
 * 获取系统菜单趋势数据
 * @returns 系统菜单7天趋势数据
 */
export const getSysBaseMenusTrend = (): Promise<{
  code: number
  data: TableTrendData[]
  msg: string
}> => {
  return request.get({
    url: '/statistics/trend',
    params: { table_name: 'sys_base_menus' }
  })
}


// 统计数据辅助函数

/**
 * 格式化统计数字
 * @param num 数字
 * @returns 格式化后的字符串
 */
export const formatStatNumber = (num: number): string => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

/**
 * 格式化趋势百分比
 * @param trend 趋势字符串
 * @returns 格式化后的趋势显示
 */
export const formatTrendPercent = (trend?: string): {
  text: string
  type: 'up' | 'down' | 'none'
} => {
  if (!trend) {
    return { text: '-', type: 'none' }
  }

  const isUp = !trend.startsWith('-')
  return {
    text: trend,
    type: isUp ? 'up' : 'down'
  }
}

/**
 * 获取图标类型配置
 * @param iconType 图标类型
 * @returns 图标配置
 */
export const getIconConfig = (iconType: string): {
  name: string
  color: string
  bgColor: string
} => {
  const iconConfigs: Record<string, { name: string; color: string; bgColor: string }> = {
    'echarts-line': {
      name: 'chart-line',
      color: '#0052d9',
      bgColor: 'rgba(0, 82, 217, 0.1)'
    },
    'echarts-bar': {
      name: 'chart-bar',
      color: '#00a870',
      bgColor: 'rgba(0, 168, 112, 0.1)'
    },
    'icon-usergroup': {
      name: 'usergroup',
      color: '#e37318',
      bgColor: 'rgba(227, 115, 24, 0.1)'
    },
    'icon-file-paste': {
      name: 'file-paste',
      color: '#951ce8',
      bgColor: 'rgba(149, 28, 232, 0.1)'
    },
    'icon-bank': {
      name: 'wallet',
      color: '#c9353f',
      bgColor: 'rgba(201, 53, 63, 0.1)'
    },
    'icon-shopping': {
      name: 'shop',
      color: '#0594fa',
      bgColor: 'rgba(5, 148, 250, 0.1)'
    }
  }

  return iconConfigs[iconType] || {
    name: 'chart',
    color: '#666666',
    bgColor: 'rgba(102, 102, 102, 0.1)'
  }
}

// 导出统计相关常量
export const STATISTICS_REFRESH_INTERVAL = 30000 // 30秒刷新一次

export const TREND_CHART_COLORS = {
  primary: '#0052d9',
  success: '#00a870',
  warning: '#e37318',
  danger: '#d54941'
}

// 统计卡片默认配置
export const DEFAULT_CARD_CONFIG = {
  loading: false,
  error: null,
  data: {
    title: '加载中...',
    number: '0',
    leftType: 'chart'
  }
}