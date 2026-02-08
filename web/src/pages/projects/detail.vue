<template>
    <div class="project-detail">
        <div class="detail-header">
            <div class="header-left">
                <t-button variant="text" shape="circle" @click="$router.back()">
                    <template #icon><t-icon name="arrow-left" /></template>
                </t-button>
                <div class="header-info">
                    <div class="title-row">
                        <h1 class="page-title">{{ project?.title || '加载中...' }}</h1>
                        <t-tag :theme="getStatusTheme(project?.status)" variant="light">{{
                            getStatusText(project?.status)
                        }}</t-tag>
                    </div>
                    <div class="page-desc text-ellipsis">{{ project?.description || '暂无简介' }}</div>
                </div>
            </div>
            <div class="header-right">
                <t-space>
                    <t-button theme="default" variant="outline" @click="refreshData">
                        <template #icon><t-icon name="refresh" /></template>
                        刷新
                    </t-button>
                    <t-button theme="primary" @click="editProject">
                        <template #icon><t-icon name="edit" /></template>
                        编辑项目
                    </t-button>
                </t-space>
            </div>
        </div>

        <div class="detail-content" v-loading="loading">
            <t-tabs v-model="activeTab" theme="normal" class="main-tabs">

                <t-tab-panel value="overview" label="项目概览">
                    <div class="tab-pane-wrapper">
                        <t-row :gutter="[16, 16]" class="stats-row">
                            <t-col :span="3">
                                <t-card :bordered="false" class="stat-card">
                                    <div class="stat-icon is-blue"><t-icon name="file-paste" /></div>
                                    <div class="stat-info">
                                        <div class="stat-label">剧本章节</div>
                                        <div class="stat-value">{{ scriptsCount }}</div>
                                    </div>
                                </t-card>
                            </t-col>
                            <t-col :span="3">
                                <t-card :bordered="false" class="stat-card">
                                    <div class="stat-icon is-green"><t-icon name="user" /></div>
                                    <div class="stat-info">
                                        <div class="stat-label">角色数量</div>
                                        <div class="stat-value">{{ charactersCount }}</div>
                                    </div>
                                </t-card>
                            </t-col>
                            <t-col :span="3">
                                <t-card :bordered="false" class="stat-card">
                                    <div class="stat-icon is-orange"><t-icon name="image" /></div>
                                    <div class="stat-info">
                                        <div class="stat-label">场景库</div>
                                        <div class="stat-value">{{ scenesCount }}</div>
                                    </div>
                                </t-card>
                            </t-col>
                            <t-col :span="3">
                                <t-card :bordered="false" class="stat-card">
                                    <div class="stat-icon is-purple"><t-icon name="gift" /></div>
                                    <div class="stat-info">
                                        <div class="stat-label">道具库</div>
                                        <div class="stat-value">{{ propsCount }}</div>
                                    </div>
                                </t-card>
                            </t-col>
                        </t-row>

                        <t-card title="项目信息" :bordered="false" class="mt-4">
                            <t-descriptions :column="2" bordered>
                                <t-descriptions-item label="项目名称">{{ project?.title }}</t-descriptions-item>
                                <t-descriptions-item label="创建时间">{{ formatDate(project?.created_at)
                                }}</t-descriptions-item>
                                <t-descriptions-item label="业务流水号">{{ project?.serialNo || '--' }}</t-descriptions-item>
                                <t-descriptions-item label="视频比例">{{ getRatioLabel(project?.settings)
                                }}</t-descriptions-item>
                                <t-descriptions-item label="剧情简介" :span="2">
                                    {{ project?.description || '暂无简介' }}
                                </t-descriptions-item>
                            </t-descriptions>
                        </t-card>
                    </div>
                </t-tab-panel>

                <t-tab-panel value="scripts" label="剧本章节">
                    <div class="tab-pane-wrapper">
                        <div class="pane-action-bar">
                            <div class="bar-title">章节列表 ({{ scriptsCount }})</div>
                            <t-button theme="primary" @click="createNewScript">
                                <template #icon><t-icon name="add" /></template>
                                新建章节
                            </t-button>
                        </div>

                        <t-table v-if="scriptsCount > 0" :data="sortedScripts" :columns="scriptColumns" row-key="id"
                            stripe hover :loading="tableLoading">
                            <template #status="{ row }">
                                <t-tag :theme="getScriptStatusTheme(row)" variant="light-outline" shape="round">
                                    {{ getScriptStatusText(row) }}
                                </t-tag>
                            </template>
                            <template #operation="{ row }">
                                <t-link theme="primary" hover="color" @click="enterScriptWorkflow(row)">
                                    <t-icon name="tools" style="margin-right:4px" />去创作
                                </t-link>
                                <t-divider layout="vertical" />
                                <t-popconfirm content="确认删除该章节？相关分镜也将被删除。" @confirm="deleteScript(row)">
                                    <t-link theme="danger" hover="color">删除</t-link>
                                </t-popconfirm>
                            </template>
                        </t-table>

                        <div v-else class="empty-placeholder">
                            <t-empty title="暂无章节" description="点击右上角创建你的第一集剧本" />
                        </div>
                    </div>
                </t-tab-panel>

                <t-tab-panel value="characters" label="角色库">
                    <div class="tab-pane-wrapper">
                        <div class="pane-action-bar">
                            <div class="bar-title">角色列表</div>
                            <t-space>
                                <t-button variant="outline" @click="openExtractDialog('character')">
                                    <template #icon><t-icon name="file-paste" /></template> 从剧本提取
                                </t-button>
                                <t-button theme="primary" @click="openAddCharacterDialog">
                                    <template #icon><t-icon name="add" /></template> 添加角色
                                </t-button>
                            </t-space>
                        </div>

                        <div v-if="charactersCount > 0" class="resource-grid">
                            <t-card v-for="char in project?.characters" :key="char.id" class="resource-card"
                                :bordered="false" hover-shadow>
                                <div class="resource-cover portrait">
                                    <t-image :src="getImageUrl(char)" fit="cover" class="cover-img">
                                        <template #error>
                                            <div class="img-error"><t-icon name="user" size="32px" /></div>
                                        </template>
                                    </t-image>
                                    <div class="resource-overlay">
                                        <t-button shape="circle" size="small" @click="editCharacter(char)"><t-icon
                                                name="edit" /></t-button>
                                        <t-tooltip content="AI生成形象">
                                            <t-button shape="circle" size="small" theme="warning"
                                                @click="generateCharacterImage(char)"><t-icon name="magic" /></t-button>
                                        </t-tooltip>
                                        <t-popconfirm content="确认删除?" @confirm="deleteCharacter(char)">
                                            <t-button shape="circle" size="small" theme="danger"><t-icon
                                                    name="delete" /></t-button>
                                        </t-popconfirm>
                                    </div>
                                </div>
                                <div class="resource-info">
                                    <div class="name-row">
                                        <span class="name">{{ char.name }}</span>
                                        <t-tag size="small" :theme="getRoleTagTheme(char.role)">{{
                                            getRoleLabel(char.role)
                                            }}</t-tag>
                                    </div>
                                    <div class="desc" :title="char.appearance">{{ char.appearance || '暂无外貌描述' }}</div>
                                </div>
                            </t-card>
                        </div>
                        <div v-else class="empty-placeholder"><t-empty description="暂无角色数据" /></div>
                    </div>
                </t-tab-panel>

                <t-tab-panel value="scenes" label="场景库">
                    <div class="tab-pane-wrapper">
                        <div class="pane-action-bar">
                            <div class="bar-title">场景列表</div>
                            <t-space>
                                <t-button variant="outline" @click="openExtractDialog('scene')">
                                    <template #icon><t-icon name="file-paste" /></template> 从剧本提取
                                </t-button>
                                <t-button theme="primary" @click="openAddSceneDialog">
                                    <template #icon><t-icon name="add" /></template> 添加场景
                                </t-button>
                            </t-space>
                        </div>

                        <div v-if="scenesCount > 0" class="resource-grid">
                            <t-card v-for="scene in scenes" :key="scene.id" class="resource-card" :bordered="false"
                                hover-shadow>
                                <div class="resource-cover landscape">
                                    <t-image :src="getImageUrl(scene)" fit="cover" class="cover-img">
                                        <template #error>
                                            <div class="img-error"><t-icon name="image" size="32px" /></div>
                                        </template>
                                    </t-image>
                                    <div class="resource-overlay">
                                        <t-button shape="circle" size="small" @click="editScene(scene)"><t-icon
                                                name="edit" /></t-button>
                                        <t-tooltip content="AI生成场景图">
                                            <t-button shape="circle" size="small" theme="warning"
                                                @click="generateSceneImage(scene)"><t-icon name="magic" /></t-button>
                                        </t-tooltip>
                                        <t-popconfirm content="确认删除?" @confirm="deleteScene(scene)">
                                            <t-button shape="circle" size="small" theme="danger"><t-icon
                                                    name="delete" /></t-button>
                                        </t-popconfirm>
                                    </div>
                                </div>
                                <div class="resource-info">
                                    <div class="name-row">
                                        <span class="name">{{ scene.name || scene.location }}</span>
                                    </div>
                                    <div class="desc" :title="scene.prompt">{{ scene.prompt || scene.description ||
                                        '暂无描述' }}
                                    </div>
                                </div>
                            </t-card>
                        </div>
                        <div v-else class="empty-placeholder"><t-empty description="暂无场景数据" /></div>
                    </div>
                </t-tab-panel>

                <t-tab-panel value="props" label="道具库">
                    <div class="tab-pane-wrapper">
                        <div class="pane-action-bar">
                            <div class="bar-title">道具列表</div>
                            <t-space>
                                <t-button variant="outline" @click="openExtractDialog('prop')">
                                    <template #icon><t-icon name="file-paste" /></template> 从剧本提取
                                </t-button>
                                <t-button theme="primary" @click="openAddPropDialog">
                                    <template #icon><t-icon name="add" /></template> 添加道具
                                </t-button>
                            </t-space>
                        </div>

                        <div v-if="propsCount > 0" class="resource-grid">
                            <t-card v-for="prop in project?.props" :key="prop.id" class="resource-card"
                                :bordered="false" hover-shadow>
                                <div class="resource-cover square">
                                    <t-image :src="getImageUrl(prop)" fit="cover" class="cover-img">
                                        <template #error>
                                            <div class="img-error"><t-icon name="gift" size="32px" /></div>
                                        </template>
                                    </t-image>
                                    <div class="resource-overlay">
                                        <t-button shape="circle" size="small" @click="editProp(prop)"><t-icon
                                                name="edit" /></t-button>
                                        <t-tooltip content="AI生成道具图">
                                            <t-button shape="circle" size="small" theme="warning"
                                                @click="generatePropImage(prop)"><t-icon name="magic" /></t-button>
                                        </t-tooltip>
                                        <t-popconfirm content="确认删除?" @confirm="deleteProp(prop)">
                                            <t-button shape="circle" size="small" theme="danger"><t-icon
                                                    name="delete" /></t-button>
                                        </t-popconfirm>
                                    </div>
                                </div>
                                <div class="resource-info">
                                    <div class="name-row">
                                        <span class="name">{{ prop.name }}</span>
                                        <t-tag size="small" variant="outline" v-if="prop.type">{{ prop.type }}</t-tag>
                                    </div>
                                    <div class="desc" :title="prop.description">{{ prop.description || '暂无描述' }}</div>
                                </div>
                            </t-card>
                        </div>
                        <div v-else class="empty-placeholder"><t-empty description="暂无道具数据" /></div>
                    </div>
                </t-tab-panel>

            </t-tabs>
        </div>

    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import dayjs from 'dayjs'

// 1. 引入 API
import { findProjects } from '@/api/projects'
import { deleteScripts, getScriptsList } from '@/api/scripts'
import { getImageUrl } from '@/utils/format' // 确保你有这个工具函数，如果没有，使用上面的本地定义

const router = useRouter()
const route = useRoute()

// ========== 状态管理 ==========
const loading = ref(false)
const tableLoading = ref(false)
const activeTab = ref('overview')
const project = ref<any>({})
const scenes = ref<any[]>([])
const scripts = ref<any[]>([])
let pollingTimer: any = null

// ========== 统计计算 ==========
const scriptsCount = computed(() => scripts.value.length)
const charactersCount = computed(() => project.value?.characters?.length || 0)
const scenesCount = computed(() => scenes.value.length)
const propsCount = computed(() => project.value?.props?.length || 0)

const sortedScripts = computed(() => {
    return [...scripts.value].sort((a, b) => (a.episodeNo || 0) - (b.episodeNo || 0))
})

// ========== 数据加载 ==========
const refreshData = () => {
    loadProjectData()
    loadScripts()
}

const loadProjectData = async () => {
    loading.value = true
    try {
        const id = route.params.id as string
        const res = await findProjects(id)
        if (res.code === 0) {
            project.value = res.data
            scenes.value = res.data.scenes || []
        } else {
            MessagePlugin.error(res.message || '获取项目详情失败')
        }
    } catch (e) {
        console.error(e)
        MessagePlugin.error('加载数据异常')
    } finally {
        loading.value = false
    }
}

const loadScripts = async () => {
    tableLoading.value = true
    try {
        const params = {
            projectId: route.params.id,
            page: 1,
            pageSize: 100
        }
        const res = await getScriptsList(params)
        if (res.code === 0) {
            const list = Array.isArray(res.data) ? res.data : (res.data?.list || [])
            scripts.value = list
        }
    } catch (e) {
        console.error(e)
    } finally {
        tableLoading.value = false
    }
}

// ========== 辅助函数 (修复了这里) ==========
const formatDate = (val: any) => val ? dayjs(val).format('YYYY-MM-DD HH:mm') : '--'

// 修复点：函数名改为 getStatusText
const getStatusText = (s: any) => {
    const map: any = { 0: '草稿', 1: '生成中', 2: '已完成', '-1': '失败' }
    return map[s] || '未知'
}

const getStatusTheme = (s: any) => {
    if (s === 2) return 'success'
    if (s === 1) return 'primary'
    if (s === -1) return 'danger'
    return 'default'
}

const getRatioLabel = (settingsStr: any) => {
    try {
        if (!settingsStr) return '默认'
        const s = typeof settingsStr === 'string' ? JSON.parse(settingsStr) : settingsStr
        return s?.ratio || '默认'
    } catch { return '--' }
}

// ========== 业务逻辑：剧本章节 ==========
const scriptColumns = [
    {
        colKey: 'episodeNo',
        title: '集数',
        width: 80,
        align: 'center',
        cell: (h, { row }) => `第${row.episodeNo}集`
    },
    { colKey: 'title', title: '标题', ellipsis: true },
    { colKey: 'status', title: '状态', width: 120, cell: 'status' },
    {
        colKey: 'updatedAt',
        title: '更新时间',
        width: 180,
        cell: (h, { row }) => formatDate(row.updatedAt)
    },
    { colKey: 'operation', title: '操作', width: 200, fixed: 'right', cell: 'operation' }
]

const getScriptStatusText = (row: any) => row.shots?.length > 0 ? "已拆分" : (row.content ? "已生成剧本" : "草稿")
const getScriptStatusTheme = (row: any) => row.shots?.length > 0 ? "success" : (row.content ? "warning" : "default")

const createNewScript = () => {
    const nextNum = scriptsCount.value + 1
    router.push({
        name: 'ProjectChapterCreate',
        params: {
            id: project.value.id,
            episodeNumber: nextNum
        }
    })
}

const enterScriptWorkflow = (row: any) => {
    router.push({
        name: 'ProjectChapterCreate',
        params: {
            id: project.value.id,
            episodeNumber: row.episodeNo
        }
    })
}

const deleteScript = async (row: any) => {
    const confirm = DialogPlugin.confirm({
        header: '确认删除',
        body: `确定要删除第 ${row.episodeNo} 集吗？`,
        onConfirm: async () => {
            try {
                const res = await deleteScripts(row.id)
                if (res.code === 0) {
                    MessagePlugin.success('删除成功')
                    loadScripts()
                } else {
                    MessagePlugin.error(res.message || '删除失败')
                }
            } catch (e) {
                MessagePlugin.error('删除失败')
            }
            confirm.hide()
        }
    })
}

// ========== 角色/场景逻辑 (占位) ==========
const getRoleLabel = (role: string) => ({ main: '主角', supporting: '配角', minor: '路人' }[role] || role)
const getRoleTagTheme = (role: string) => ({ main: 'danger', supporting: 'primary', minor: 'default' }[role] || 'default')

// 对话框状态
const addCharacterDialogVisible = ref(false)
const addSceneDialogVisible = ref(false)
const extractDialogVisible = ref(false)
const extractType = ref('character')
const selectedExtractScriptId = ref(null)
const editingCharacter = ref(null)
const editingScene = ref(null)
const tempFileList = ref([])
const newCharacter = ref({ name: "", role: "supporting", appearance: "", personality: "", description: "", image_url: "", local_path: "" })
const newScene = ref({ location: "", prompt: "", image_url: "", local_path: "" })
const formLoading = ref(false)
const uploadConfig = { action: '', headers: {} } // 根据实际配置

// 简单的空函数占位，防止模板报错
const openAddCharacterDialog = () => { addCharacterDialogVisible.value = true }
const openAddSceneDialog = () => { addSceneDialogVisible.value = true }
const openExtractDialog = (type: string) => { extractType.value = type; extractDialogVisible.value = true }
const editCharacter = (char: any) => { }
const deleteCharacter = (char: any) => { }
const generateCharacterImage = (char: any) => { }
const editScene = (scene: any) => { }
const deleteScene = (scene: any) => { }
const generateSceneImage = (scene: any) => { }
const openAddPropDialog = () => { }
const editProp = (prop: any) => { }
const deleteProp = (prop: any) => { }
const generatePropImage = (prop: any) => { }
const saveCharacter = () => { }
const handleUploadSuccess = () => { }
const saveScene = () => { }
const confirmExtract = () => { }

// 编辑项目本身
const editProject = () => {
    MessagePlugin.info('请在项目列表页进行设置')
}

// 初始化
onMounted(() => {
    if (route.query.tab) activeTab.value = route.query.tab as string
    loadProjectData()
    loadScripts()
})
</script>

<style scoped lang="less">
/* 基础变量 */
@bg-color: #f2f3f5;
@card-radius: 8px;

.project-detail {
    background-color: @bg-color;
    min-height: 100vh;
}

/* 1. Header 样式 */
.detail-header {
    background: #fff;
    padding: 16px 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

    .header-left {
        display: flex;
        align-items: center;
        gap: 16px;

        .header-info {
            .title-row {
                display: flex;
                align-items: center;
                gap: 12px;
                margin-bottom: 4px;

                .page-title {
                    font-size: 20px;
                    font-weight: 700;
                    color: var(--td-text-color-primary);
                    margin: 0;
                }
            }

            .page-desc {
                font-size: 13px;
                color: var(--td-text-color-secondary);
                max-width: 600px;
            }
        }
    }
}

/* 2. Content 样式 */
.detail-content {
    padding: 24px;
}

.main-tabs {
    background: #fff;
    border-radius: @card-radius;
    padding: 16px 0;
    min-height: 600px;

    :deep(.t-tabs__nav-container) {
        padding: 0 24px;
    }

    :deep(.t-tabs__content) {
        padding: 24px;
    }
}

/* 统计卡片 */
.stat-card {
    text-align: center;
    cursor: default;
    transition: transform 0.2s;
    background: #f9f9f9;

    &:hover {
        transform: translateY(-2px);
    }

    .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 50%;
        margin: 0 auto 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 24px;

        &.is-blue {
            background: rgba(0, 82, 217, 0.1);
            color: #0052D9;
        }

        &.is-green {
            background: rgba(43, 164, 113, 0.1);
            color: #2BA471;
        }

        &.is-orange {
            background: rgba(237, 123, 47, 0.1);
            color: #ED7B2F;
        }

        &.is-purple {
            background: rgba(114, 46, 209, 0.1);
            color: #722ED1;
        }
    }

    .stat-value {
        font-size: 24px;
        font-weight: 700;
        color: var(--td-text-color-primary);
    }

    .stat-label {
        font-size: 13px;
        color: var(--td-text-color-secondary);
    }
}

/* 通用操作栏 */
.pane-action-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .bar-title {
        font-size: 16px;
        font-weight: 600;
        color: var(--td-text-color-primary);
        padding-left: 12px;
        border-left: 4px solid var(--td-brand-color);
        line-height: 1;
    }
}

/* 资源网格 (角色/场景/道具) */
.resource-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 20px;
}

.resource-card {
    overflow: hidden;
    border-radius: 8px;

    .resource-cover {
        position: relative;
        background: #f3f3f3;
        overflow: hidden;

        &.portrait {
            aspect-ratio: 3/4;
        }

        &.landscape {
            aspect-ratio: 16/9;
        }

        &.square {
            aspect-ratio: 1/1;
        }

        .cover-img {
            width: 100%;
            height: 100%;
            display: block;
        }

        .img-error {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #ccc;
        }

        .resource-overlay {
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0, 0, 0, 0.4);
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
            opacity: 0;
            transition: opacity 0.2s;
        }

        &:hover .resource-overlay {
            opacity: 1;
        }
    }

    .resource-info {
        padding: 12px;

        .name-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 6px;

            .name {
                font-weight: 600;
                font-size: 14px;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
            }
        }

        .desc {
            font-size: 12px;
            color: var(--td-text-color-placeholder);
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
            height: 36px;
        }
    }
}

.empty-placeholder {
    padding: 40px 0;
    display: flex;
    justify-content: center;
}

.text-ellipsis {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.mt-4 {
    margin-top: 24px;
}
</style>