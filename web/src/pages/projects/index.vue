<template>
    <div class="projects-list">
        <div class="welcome-section">
            <div class="welcome-text">
                <h2 class="welcome-title">灵果短剧工作台</h2>
                <p class="welcome-desc">一站式 AI 短剧创作平台，从剧本灵感到视频成片</p>
            </div>
            <div class="welcome-stat">
            </div>
        </div>

        <t-card :bordered="false" class="operation-bar">
            <div class="operation-wrapper">
                <div class="left-actions">
                    <t-button theme="primary" size="large" @click="onCreate">
                        <template #icon><t-icon name="add" /></template>
                        新建项目
                    </t-button>

                    <t-radio-group v-model="viewType" variant="default-filled">
                        <t-radio-button value="card"><t-icon name="app" /> 卡片</t-radio-button>
                        <t-radio-button value="list"><t-icon name="view-list" /> 列表</t-radio-button>
                    </t-radio-group>
                </div>

                <div class="right-search">
                    <t-input v-model="searchInfo.title" placeholder="搜索项目名称..." clearable @enter="onSearch"
                        style="width: 280px">
                        <template #suffix-icon>
                            <t-icon name="search" style="cursor: pointer" @click="onSearch" />
                        </template>
                    </t-input>
                    <t-button theme="default" variant="text" @click="onRefresh">
                        <t-icon name="refresh" />
                    </t-button>
                </div>
            </div>
        </t-card>

        <div class="content-area" v-loading="loading">

            <div v-if="viewType === 'card'" class="card-view">
                <t-row :gutter="[24, 24]">
                    <t-col :xs="12" :sm="6" :md="4" :lg="3" :xl="3">
                        <div class="project-card create-card" @click="onCreate">
                            <div class="create-inner">
                                <div class="icon-box"><t-icon name="add" /></div>
                                <span>创建新短剧</span>
                            </div>
                        </div>
                    </t-col>

                    <t-col v-for="item in tableData" :key="item.id" :xs="12" :sm="6" :md="4" :lg="3" :xl="3">
                        <t-card class="project-card" :bordered="false" hover-shadow>
                            <div class="card-cover-wrapper" @click="enterStudio(item)">
                                <t-image :src="item.image || ''" fit="cover" class="card-image"
                                    :style="{ aspectRatio: getAspectRatio(item.settings) }" :lazy="true">
                                    <template #error>
                                        <div class="image-placeholder">
                                            <t-icon name="image" size="24px" />
                                        </div>
                                    </template>
                                </t-image>

                                <div class="status-tag" :class="`status-${item.status}`">
                                    {{ getStatusLabel(item.status) }}
                                </div>

                                <div class="hover-overlay">
                                    <t-button shape="circle" theme="primary" @click.stop="enterStudio(item)">
                                        <t-icon name="tools" />
                                    </t-button>
                                    <t-button shape="circle" variant="outline" theme="default" class="btn-white"
                                        @click.stop="onEdit(item)">
                                        <t-icon name="edit" />
                                    </t-button>
                                    <t-popconfirm content="确认删除该项目？" @confirm="onDelete(item)">
                                        <t-button shape="circle" variant="outline" theme="danger" class="btn-white"
                                            @click.stop>
                                            <t-icon name="delete" />
                                        </t-button>
                                    </t-popconfirm>
                                </div>
                            </div>

                            <div class="card-info">
                                <div class="info-title" :title="item.title">{{ item.title }}</div>
                                <div class="info-meta">
                                    <span class="meta-time">{{ formatDate(item.createdAt) }}</span>
                                    <span class="meta-duration" v-if="item.totalDuration">{{
                                        formatDuration(item.totalDuration)
                                        }}</span>
                                </div>
                            </div>
                        </t-card>
                    </t-col>
                </t-row>

                <div v-if="tableData.length === 0 && !loading" class="empty-container">
                    <t-empty description="暂无项目，快去创建一个吧" />
                </div>

                <div class="pagination-container" v-if="tableData.length > 0">
                    <t-pagination v-model="pagination.current" v-model:pageSize="pagination.pageSize"
                        :total="pagination.total" @change="onPageChange" show-jumper />
                </div>
            </div>

            <t-card v-else :bordered="false" class="list-view-card">
                <t-table :data="tableData" :columns="columns" :pagination="pagination" row-key="id"
                    @page-change="onPageChange" hover>
                    <template #image="{ row }">
                        <t-image :src="row.image" fit="cover" shape="round"
                            style="width: 48px; height: 48px; border: 1px solid var(--td-component-stroke);" />
                    </template>
                    <template #status="{ row }">
                        <t-tag :theme="getStatusTheme(row.status)" variant="light-outline" shape="round">
                            {{ getStatusLabel(row.status) }}
                        </t-tag>
                    </template>
                    <template #action="{ row }">
                        <t-link theme="primary" hover="color" @click="enterStudio(row)">
                            <t-icon name="tools" style="margin-right: 4px" />去创作
                        </t-link>
                        <t-divider layout="vertical" />
                        <t-link theme="default" hover="color" @click="onEdit(row)">设置</t-link>
                        <t-divider layout="vertical" />
                        <t-popconfirm content="确认删除？" @confirm="onDelete(row)">
                            <t-link theme="danger" hover="color">删除</t-link>
                        </t-popconfirm>
                    </template>
                </t-table>
            </t-card>
        </div>

        <t-dialog v-model:visible="drawerVisible" :header="formType === 'create' ? '创建新项目' : '项目设置'" width="600px"
            :close-on-overlay-click="false" :confirm-btn="{
                content: formType === 'create' ? '立即创建' : '保存修改',
                theme: 'primary',
                loading: submitLoading
            }" @confirm="onSubmit" @close="onCancel">
            <t-form ref="formRef" :data="formData" :rules="rules" label-align="top" class="project-form">

                <t-form-item label="项目名称" name="title">
                    <t-input v-model="formData.title" placeholder="例如：霸道总裁爱上我 (第1季)" size="large" />
                </t-form-item>

                <t-row :gutter="24">
                    <t-col :span="6">
                        <t-form-item label="视频比例" name="ratio">
                            <div class="ratio-select-group">
                                <div class="ratio-option" :class="{ active: getFormDataSettings().ratio === '9:16' }"
                                    @click="updateSettings('ratio', '9:16')">
                                    <div class="ratio-preview vertical"></div>
                                    <span class="ratio-label">9:16 (竖屏)</span>
                                </div>
                                <div class="ratio-option" :class="{ active: getFormDataSettings().ratio === '16:9' }"
                                    @click="updateSettings('ratio', '16:9')">
                                    <div class="ratio-preview horizontal"></div>
                                    <span class="ratio-label">16:9 (横屏)</span>
                                </div>
                            </div>
                        </t-form-item>
                    </t-col>

                    <t-col :span="6">
                        <t-form-item label="项目封面" name="image">
                            <t-upload v-model="tempFileList" :action="uploadConfig.action"
                                :headers="uploadConfig.headers" theme="image" accept="image/*" :max="1"
                                :show-image-filename="false" @success="handleImageUploadSuccess"
                                @fail="handleUploadFail">
                            </t-upload>
                        </t-form-item>
                    </t-col>
                </t-row>

                <t-form-item label="剧情简介" name="description">
                    <t-textarea v-model="formData.description" placeholder="简单描述一下故事大纲..."
                        :autosize="{ minRows: 3, maxRows: 5 }" />
                </t-form-item>

                <input type="hidden" v-model="formData.adminId" />
            </t-form>
        </t-dialog>
    </div>
</template>

<script setup lang="tsx">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import dayjs from 'dayjs'

import {
    createProjects,
    deleteProjects,
    updateProjects,
    getProjectsList,
    findProjects
} from '@/api/projects'

const router = useRouter()

// ========== 1. 状态定义 ==========
const viewType = ref('card') // 'card' | 'list'
const loading = ref(false)
const tableData = ref([])
const pagination = reactive({
    current: 1,
    pageSize: 12,
    total: 0,
    showJumper: true
})
const searchInfo = ref({ title: '' })

// 状态字典
const statusMap = {
    0: { label: '草稿', theme: 'default' },
    1: { label: '生成中', theme: 'primary' },
    2: { label: '已完成', theme: 'success' },
    '-1': { label: '失败', theme: 'danger' }
}

// ========== 2. 数据获取 ==========
const getTableData = async () => {
    loading.value = true
    try {
        const params = {
            page: pagination.current,
            pageSize: pagination.pageSize,
            title: searchInfo.value.title || undefined
        }
        const res = await getProjectsList(params)
        if (res.code === 0) {
            // 兼容后端返回结构
            const list = res.data.list || res.data || []
            tableData.value = list
            pagination.total = res.data.total || list.length || 0
        }
    } catch (e) {
        MessagePlugin.error('数据加载失败')
    } finally {
        loading.value = false
    }
}

const onSearch = () => {
    pagination.current = 1
    getTableData()
}

const onRefresh = () => {
    getTableData()
    MessagePlugin.success('刷新成功')
}

const onPageChange = (pageInfo: any) => {
    pagination.current = pageInfo.current
    pagination.pageSize = pageInfo.pageSize
    getTableData()
}

// ========== 3. 表单与弹窗逻辑 ==========
const drawerVisible = ref(false)
const submitLoading = ref(false)
const formType = ref('create')
const formRef = ref()
const tempFileList = ref([]) 
// 表单数据模型
const formData = ref({
    id: null,
    adminId: 1, // 默认归属用户
    title: '',
    description: '',
    status: 0,
    image: '', // 最终提交给后端的图片URL
    totalDuration: 0,
    settings: JSON.stringify({ ratio: '9:16' }), // 默认配置
    serialNo: ''
})

const rules = {
    title: [{ required: true, message: '项目名称不能为空' }]
}

// 解析/获取 Settings 对象
const getFormDataSettings = () => {
    try {
        return typeof formData.value.settings === 'string'
            ? JSON.parse(formData.value.settings)
            : formData.value.settings || { ratio: '9:16' }
    } catch (e) {
        return { ratio: '9:16' }
    }
}

// 更新 Settings 中的某个字段
const updateSettings = (key, value) => {
    const currentSettings = getFormDataSettings()
    currentSettings[key] = value
    formData.value.settings = JSON.stringify(currentSettings)
}

// 打开新建
const onCreate = () => {
    formType.value = 'create'
    formData.value = {
        id: null,
        adminId: 1,
        title: '',
        description: '',
        status: 0,
        image: '',
        totalDuration: 0,
        settings: JSON.stringify({ ratio: '9:16' }),
        serialNo: generateSerialNo() // 生成一个临时流水号
    }
    tempFileList.value = []
    drawerVisible.value = true
}

// 打开编辑
const onEdit = async (row: any) => {
    try {
        // 获取最新详情
        const res = await findProjects(row.id)
        if (res.code === 0) {
            formType.value = 'update'
            const data = res.data

            // 数据回显
            formData.value = {
                ...data,
                // 确保 settings 是字符串
                settings: typeof data.settings === 'object' ? JSON.stringify(data.settings) : (data.settings || JSON.stringify({ ratio: '9:16' }))
            }

            // 图片回显处理
            if (data.image) {
                tempFileList.value = [{ url: data.image, name: '封面图' }]
            } else {
                tempFileList.value = []
            }

            drawerVisible.value = true
        }
    } catch (e) {
        MessagePlugin.error('获取详情失败')
    }
}

// ========== 4. 图片上传修复核心逻辑 ==========
const getAuthToken = () => localStorage.getItem('token')

const uploadConfig = reactive({
    action: import.meta.env.VITE_API_URL + '/admin/v1/upload/singleUpload',
    headers: computed(() => ({ 'Authorization': `${getAuthToken()}` })),
})

// 监听上传成功
const handleImageUploadSuccess = (context: any) => {
    console.log('Upload Success Context:', context)
    // TDesign 的 context.response 是后端返回的原始数据
    const rawResponse = context.response

    // 根据你的后端结构解析 { code: 0, data: { file_url: '...' } }
    if (rawResponse && (rawResponse.code === 0 || rawResponse.code === 200)) {
        const url = rawResponse.data?.file_url || rawResponse.data?.url
        if (url) {
            // 核心修复：直接赋值给 formData.image
            formData.value.image = url
            MessagePlugin.success('封面上传成功')
        } else {
            MessagePlugin.warning('上传成功但未获取到链接')
        }
    } else {
        MessagePlugin.error(rawResponse?.message || '上传失败')
    }
}

const handleUploadFail = () => {
    MessagePlugin.error('网络错误，上传失败')
}

// ========== 5. 提交逻辑 ==========
const onSubmit = async () => {
    const valid = await formRef.value.validate()
    if (valid !== true) return

    submitLoading.value = true
    try {
        // 构造 Payload
        const payload = { ...formData.value }

        // 确保 settings 是字符串 (后端要求可能是字符串)
        if (typeof payload.settings === 'object') {
            payload.settings = JSON.stringify(payload.settings)
        }

        // 如果 tempFileList 为空，说明用户清除了图片
        if (tempFileList.value.length === 0) {
            payload.image = ''
        }

        let res
        if (formType.value === 'create') {
            res = await createProjects(payload)
        } else {
            res = await updateProjects(payload.id, payload)
        }

        if (res.code === 0) {
            MessagePlugin.success('保存成功')
            drawerVisible.value = false
            getTableData()
        } else {
            MessagePlugin.error(res.message || '操作失败')
        }
    } catch (e) {
        console.error(e)
        MessagePlugin.error('系统异常')
    } finally {
        submitLoading.value = false
    }
}

const onDelete = async (row: any) => {
    const res = await deleteProjects(row.id)
    if (res.code === 0) {
        MessagePlugin.success('已删除')
        getTableData()
    }
}

// ========== 6. 辅助/格式化 ==========
const formatDate = (val: any) => val ? dayjs(val).format('YYYY-MM-DD HH:mm') : '--'
const formatDuration = (seconds: any) => {
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    return `${m}:${s.toString().padStart(2, '0')}`
}
const getStatusLabel = (s: any) => statusMap[s]?.label || '未知'
const getStatusTheme = (s: any) => statusMap[s]?.theme || 'default'

const getAspectRatio = (settingsStr: any) => {
    try {
        const s = JSON.parse(settingsStr || '{}')
        return s.ratio === '16:9' ? '16/9' : '9/16'
    } catch {
        return '9/16'
    }
}

const generateSerialNo = () => 'SN' + dayjs().format('YYYYMMDDHHmmss')

const enterStudio = (row: any) => {
    router.push({
        name: 'ProjectDetail',
        params: {
            id: row.id         // 传递项目ID
        }
    })
}

// 列表列定义
const columns = [
    { colKey: 'image', title: '封面', width: 80, cell: 'image' },
    { colKey: 'title', title: '项目名称', ellipsis: true },
    { colKey: 'status', title: '状态', width: 100, cell: 'status' },
    { colKey: 'totalDuration', title: '时长', width: 100, cell: (h, { row }) => formatDuration(row.totalDuration) },
    { colKey: 'created_at', title: '创建时间', width: 180, cell: (h, { row }) => formatDate(row.createdAt) },
    { colKey: 'action', title: '操作', width: 200, fixed: 'right', cell: 'action' }
]

onMounted(() => {
    getTableData()
})
</script>

<style scoped lang="less">
/* 基础变量 */
@bg-color: #f2f3f5;
@card-radius: 12px;
@primary-color: #0052D9;

.projects-list {
    padding: 24px;
    background-color: @bg-color;
    min-height: 100vh;
}

/* 欢迎区 */
.welcome-section {
    margin-bottom: 24px;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;

    .welcome-title {
        font-size: 24px;
        font-weight: 700;
        color: var(--td-text-color-primary);
        margin-bottom: 8px;
    }

    .welcome-desc {
        color: var(--td-text-color-secondary);
        font-size: 14px;
    }
}

/* 操作栏 */
.operation-bar {
    margin-bottom: 24px;
    border-radius: 8px;

    .operation-wrapper {
        display: flex;
        justify-content: space-between;
        align-items: center;
        flex-wrap: wrap;
        gap: 16px;
    }

    .left-actions {
        display: flex;
        align-items: center;
        gap: 20px;
    }

    .right-search {
        display: flex;
        align-items: center;
        gap: 12px;
    }
}

/* 卡片视图样式 */
.project-card {
    border-radius: @card-radius;
    overflow: hidden;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    height: 100%;

    &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);

        .hover-overlay {
            opacity: 1 !important;
        }
    }
}

/* 新建卡片 */
.create-card {
    border: 2px dashed var(--td-component-stroke);
    background: #fff;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 280px;

    &:hover {
        border-color: @primary-color;
        background: var(--td-brand-color-light);

        .create-inner {
            color: @primary-color;
        }
    }

    .create-inner {
        text-align: center;
        color: var(--td-text-color-placeholder);
        transition: color 0.2s;

        .icon-box {
            font-size: 40px;
            margin-bottom: 12px;
        }

        span {
            font-weight: 500;
        }
    }
}

/* 卡片封面区域 */
.card-cover-wrapper {
    position: relative;
    background: #000;
    cursor: pointer;
    overflow: hidden;

    .card-image {
        width: 100%;
        display: block;
    }

    .image-placeholder {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f3f3f3;
        color: #ccc;
        aspect-ratio: 9/16;
    }

    .status-tag {
        position: absolute;
        top: 10px;
        left: 10px;
        padding: 2px 8px;
        border-radius: 4px;
        font-size: 12px;
        color: #fff;
        backdrop-filter: blur(4px);

        &.status-0 {
            background: rgba(0, 0, 0, 0.5);
        }

        &.status-1 {
            background: rgba(0, 82, 217, 0.8);
        }

        &.status-2 {
            background: rgba(43, 164, 113, 0.9);
        }

        &.status--1 {
            background: rgba(213, 73, 65, 0.9);
        }
    }

    /* 悬浮遮罩 */
    .hover-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.4);
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 12px;
        opacity: 0;
        transition: opacity 0.3s;

        .btn-white {
            color: #fff;
            border-color: #fff;

            &:hover {
                background: rgba(255, 255, 255, 0.2);
            }
        }
    }
}

/* 卡片信息区 */
.card-info {
    padding: 12px 16px;

    .info-title {
        font-weight: 600;
        font-size: 16px;
        margin-bottom: 8px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        color: var(--td-text-color-primary);
    }

    .info-meta {
        display: flex;
        justify-content: space-between;
        font-size: 12px;
        color: var(--td-text-color-placeholder);
    }
}

/* 分页 */
.pagination-container {
    margin-top: 24px;
    display: flex;
    justify-content: flex-end;
}

/* 比例选择器样式 */
.ratio-select-group {
    display: flex;
    gap: 16px;

    .ratio-option {
        cursor: pointer;
        border: 1px solid var(--td-component-stroke);
        padding: 12px;
        border-radius: 8px;
        text-align: center;
        transition: all 0.2s;
        flex: 1;

        &:hover {
            border-color: @primary-color;
        }

        &.active {
            background: var(--td-brand-color-light);
            border-color: @primary-color;
            color: @primary-color;

            .ratio-preview {
                border-color: @primary-color;
                background: rgba(0, 82, 217, 0.1);
            }
        }

        .ratio-preview {
            border: 2px solid var(--td-text-color-placeholder);
            margin: 0 auto 8px;

            &.vertical {
                width: 24px;
                height: 42px;
            }

            &.horizontal {
                width: 42px;
                height: 24px;
                margin-top: 9px;
                margin-bottom: 9px;
            }
        }

        .ratio-label {
            font-size: 12px;
        }
    }
}

/* 空状态 */
.empty-container {
    min-height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>