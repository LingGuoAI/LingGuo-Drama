// 统计相关类型定义

export interface StatisticsCard {
  title: string        // 标题
  number: string       // 数量
  upTrend?: string     // 上升趋势
  downTrend?: string   // 下降趋势
  leftType: string     // 左侧图标类型
}

export interface TableTrendData {
  date: string         // 日期 (YYYY-MM-DD)
  count: number        // 数量
}

export interface TableTrendParams {
  table_name: string   // 表名
}

export interface DetailStatistics {
  characters: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  } // 角色统计
  projects: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  } // 短剧项目统计
  scripts: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  } // 剧本统计
  shots: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  } // 镜头表统计
  admins: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  } // 系统管理员统计
  sys_base_menus: {
    today: number      // 今日新增
    yesterday: number  // 昨日新增
    this_week: number  // 本周新增
    last_week: number  // 上周新增
    this_month: number // 本月新增
    last_month: number // 上月新增
  } // 系统菜单统计
}

export interface IconConfig {
  name: string         // 图标名称
  color: string        // 图标颜色
  bgColor: string      // 背景颜色
}

export interface TrendFormat {
  text: string         // 显示文本
  type: 'up' | 'down' | 'none' // 趋势类型
}

export interface StatisticsResponse {
  code: number
  data: StatisticsCard[]
  msg: string
}

export interface TrendResponse {
  code: number
  data: TableTrendData[]
  msg: string
}

export interface DetailStatisticsResponse {
  code: number
  data: DetailStatistics
  msg: string
}

export interface StatisticsCardState {
  loading: boolean
  error: string | null
  data: StatisticsCard
}

export interface StatisticsPageState {
  cards: StatisticsCard[]
  detail: DetailStatistics | null
  loading: boolean
  error: string | null
  lastUpdateTime: string
}

// 统计相关常量
export const STATISTICS_CONFIG = {
  REFRESH_INTERVAL: 30000, // 30秒刷新间隔
  CARD_ANIMATION_DELAY: 100, // 卡片动画延迟
  TREND_DAYS: 7, // 趋势数据天数
} as const

// 图标类型枚举
export const ICON_TYPES = {
  ECHARTS_LINE: 'echarts-line',
  ECHARTS_BAR: 'echarts-bar',
  ICON_USERGROUP: 'icon-usergroup',
  ICON_FILE_PASTE: 'icon-file-paste',
  ICON_BANK: 'icon-bank',
  ICON_SHOPPING: 'icon-shopping',
} as const

export type IconType = typeof ICON_TYPES[keyof typeof ICON_TYPES]
