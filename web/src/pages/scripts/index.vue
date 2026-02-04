<template>
    <div class="scripts-list">
        <!-- 搜索区域 -->
        <t-card class="search-form">
            <t-form ref="searchFormRef" :model="searchInfo" label-align="left" label-width="100px" @submit="onSearch">
                <t-row :gutter="[24, 24]">
                    <!-- 所属项目ID搜索 -->
                    <t-col :span="4">
                        <t-form-item label="短剧项目" name="projectId">
                            <t-select v-model="searchInfo.projectId" placeholder="请选择短剧项目" clearable
                                :loading="projectsSelectLoading">
                                <t-option v-for="item in projectsSelectData" :key="item.id" :label="item.title"
                                    :value="item.id"></t-option>
                            </t-select>
                        </t-form-item>
                    </t-col>
                    <!-- 分集标题搜索 -->
                    <t-col :span="4">
                        <t-form-item label="分集标题" name="title">
                            <t-input v-model="searchInfo.title" placeholder="请输入分集标题" clearable>
                            </t-input>
                        </t-form-item>
                    </t-col>
                    <!-- 第几集搜索 -->
                    <t-col :span="4">
                        <t-form-item label="第几集" name="episodeNo">
                            <t-input v-model="searchInfo.episodeNo" placeholder="请输入第几集" clearable type="number">
                            </t-input>
                        </t-form-item>
                    </t-col>
                    <!-- 是否定稿搜索 -->
                    <t-col :span="4">
                        <t-form-item label="是否定稿" name="isLocked">
                            <t-select v-model="searchInfo.isLocked" placeholder="请选择是否定稿" clearable>
                                <t-option v-for="item in isLockedStatusOptions" :key="item.value" :label="item.label"
                                    :value="item.value"></t-option>
                            </t-select>
                        </t-form-item>
                    </t-col>
                </t-row>

                <!-- 搜索按钮 -->
                <t-row style="margin-top: 20px;display: flex;justify-content: flex-end;">
                    <t-form-item label=" ">
                        <t-space>
                            <t-button theme="primary" @click="onSearch">
                                <template #icon><t-icon name="search"></t-icon></template>
                                搜索
                            </t-button>
                            <t-button variant="outline" @click="onReset">
                                <template #icon><t-icon name="refresh"></t-icon></template>
                                重置
                            </t-button>
                        </t-space>
                    </t-form-item>
                </t-row>
            </t-form>
        </t-card>

        <!-- 表格区域 -->
        <t-card>
            <!-- 操作按钮 -->
            <div>
                <t-space>
                    <t-button theme="primary" @click="onCreate">
                        <template #icon><t-icon name="add"></t-icon></template>
                        新增
                    </t-button>
                    <t-button theme="default" variant="outline" @click="onRefresh">
                        <template #icon><t-icon name="refresh"></t-icon></template>
                        刷新
                    </t-button>
                </t-space>
            </div>

            <!-- 表格 -->
            <t-table ref="tableRef" :data="tableData" :columns="columns" :loading="loading" :pagination="pagination"
                row-key="id" @page-change="onPageChange" @page-size-change="onPageSizeChange" hover />
        </t-card>
        <!-- 新增/编辑抽屉 -->
        <t-dialog v-model:visible="drawerVisible" :header="drawerTitle" width="600px" size="large" :confirm-btn="{
            content: '确定',
            theme: 'primary',
            loading: submitLoading
        }" @confirm="onSubmit" @cancel="onCancel">
            <t-form ref="formRef" :data="formData" :rules="rules" label-align="left" label-width="100px"
                @submit="onSubmit">
                <t-form-item label="短剧项目" name="projectId">
                    <t-select v-model="formData.projectId" placeholder="请选择短剧项目" clearable
                        :loading="projectsSelectLoading" filterable :status="!formData.projectId ? 'error' : 'default'">
                        <t-option v-for="item in projectsSelectData" :key="item.id" :label="item.title"
                            :value="item.id"></t-option>
                    </t-select>
                </t-form-item>
                <t-form-item label="分集标题" name="title">
                    <t-input v-model="formData.title" clearable placeholder="请输入分集标题" :maxlength="255"
                        show-word-limit />
                </t-form-item>
                <t-form-item label="剧本正文" name="content">
                    <!-- 富文本编辑器 -->
                    <div class="richtext-editor-container">
                        <t-textarea v-model="formData.content" placeholder="请输入剧本正文"
                            :autosize="{ minRows: 6, maxRows: 20 }" />
                    </div>
                </t-form-item>
                <t-form-item label="大纲/简介" name="outline">
                    <t-input v-model="formData.outline" clearable placeholder="请输入大纲/简介" :maxlength="255"
                        show-word-limit />
                </t-form-item>
                <t-form-item label="第几集" name="episodeNo"><t-input v-model="formData.episodeNo" placeholder="请输入第几集"
                        type="number" clearable />
                </t-form-item>
                <t-form-item label="是否定稿" name="isLocked">
                    <t-select v-model="formData.isLocked" placeholder="请选择是否定稿">
                        <t-option v-for="item in isLockedStatusOptions" :key="item.value" :label="item.label"
                            :value="item.value"></t-option>
                    </t-select>
                </t-form-item>
            </t-form>
        </t-dialog>


        <t-dialog v-model:visible="detailVisible" header="查看详情" width="600px" size="large" :footer="false"
            :close-btn="true" :show-overlay="true" @close="detailVisible = false">
            <t-descriptions :column="1" layout="vertical" bordered
                :content-style="{ overflowWrap: 'break-word', whiteSpace: 'normal' }">
                <t-descriptions-item label="短剧项目">
                    <span v-if="detailData.projects">
                        {{ detailData.projects.title }}
                    </span>
                    <span v-else style="color: var(--td-text-color-placeholder);">--</span>
                </t-descriptions-item>
                <t-descriptions-item label="分集标题">
                    <span v-if="detailData.title !== null && detailData.title !== undefined && detailData.title !== ''">
                        {{ detailData.title }}
                    </span>
                    <span v-else style="color: var(--td-text-color-placeholder);">--</span>
                </t-descriptions-item>
                <t-descriptions-item label="剧本正文">
                    <div v-if="detailData.content" class="richtext-content-preview">
                        <div v-html="detailData.content"></div>
                    </div>
                    <span v-else style="color: var(--td-text-color-placeholder);">无内容</span>
                </t-descriptions-item>
                <t-descriptions-item label="大纲/简介">
                    <span
                        v-if="detailData.outline !== null && detailData.outline !== undefined && detailData.outline !== ''">
                        {{ detailData.outline }}
                    </span>
                    <span v-else style="color: var(--td-text-color-placeholder);">--</span>
                </t-descriptions-item>
                <t-descriptions-item label="第几集">
                    <span
                        v-if="detailData.episodeNo !== null && detailData.episodeNo !== undefined && detailData.episodeNo !== ''">
                        {{ detailData.episodeNo }}
                    </span>
                    <span v-else style="color: var(--td-text-color-placeholder);">--</span>
                </t-descriptions-item>
                <t-descriptions-item label="是否定稿">
                    <t-tag :theme="getStatusTagTheme(detailData.isLocked, isLockedStatusOptions)" variant="light">
                        {{ getStatusLabel(detailData.isLocked, isLockedStatusOptions) }}
                    </t-tag>
                </t-descriptions-item>
            </t-descriptions>
        </t-dialog>
    </div>
</template>

<script setup lang="tsx">
import { ref, reactive, computed, onMounted, nextTick, onBeforeUnmount, shallowRef, watch } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import {
    createScripts,
    deleteScripts,
    updateScripts,
    findScripts,
    getScriptsList,
    getProjectsSelectList
} from '@/api/scripts'
import { formatDate, getImageUrl } from '@/utils/format'

defineOptions({
    name: 'ScriptsList'
})

const router = useRouter()

// ========== 状态选项定义 ==========
// 是否定稿状态选项
const isLockedStatusOptions = ref([{ "value": 0, "label": "否" }, { "value": 1, "label": "是" }])

// 表格相关
const tableRef = ref()
const loading = ref(false)
const tableData = ref([])
const showAllQuery = ref(false)

// 分页
const pagination = reactive({
    current: 1,
    pageSize: 10,
    total: 0,
    showJumper: true,
    showSizeChanger: true
})

// 搜索表单
const searchFormRef = ref()
const searchInfo = ref({
    projectId: undefined,
    title: undefined,
    episodeNo: undefined,
    isLocked: undefined,
})

// 计算是否有搜索项和隐藏搜索项
const hasSearchItems = computed(() => {
    return true
})

const hasHiddenSearchItems = computed(() => {
    return false
})

// 状态处理辅助函数 - 支持动态状态选项
const getStatusLabel = (value, options) => {
    const option = options.find(item => item.value === value)
    return option ? option.label : value
}

const getStatusTagTheme = (value, options) => {
    if (options.length === 2) {
        // 二元状态：0-成功色，1-警告色
        return value === 0 ? 'success' : 'warning'
    } else {
        // 多元状态：根据索引分配颜色
        const themes = ['default', 'warning', 'success', 'danger', 'primary']
        return themes[value % themes.length]
    }
}
// 短剧项目选择数据相关
const projectsSelectData = ref([])
const projectsSelectLoading = ref(false)

// 获取短剧项目选择数据
const getProjectsSelectData = async () => {
    projectsSelectLoading.value = true
    try {
        const res = await getProjectsSelectList()
        if (res.code === 0) {
            // 不需要再进行数据格式化，因为后端已经返回了正确格式
            projectsSelectData.value = res.data || []
        } else {
            projectsSelectData.value = []
            console.error('获取短剧项目数据失败:', res.message)
        }
    } catch (error) {
        console.error('获取短剧项目数据失败:', error)
        projectsSelectData.value = []
    } finally {
        projectsSelectLoading.value = false
    }
}

// 表格列配置
const columns = computed(() => [
    {
        title: '短剧项目',
        colKey: 'projectId',
        sorter: false,
        cell: (h, { row }) => {
            const relationObj = row.projects
            if (relationObj && relationObj.title) {
                return relationObj.title
            }
            return '--'
        }
    },
    {
        title: '分集标题',
        colKey: 'title',
        cell: (h, { row }) => row.title || '--'
    },
    {
        title: '大纲/简介',
        colKey: 'outline',
        cell: (h, { row }) => row.outline || '--'
    },
    {
        title: '第几集',
        colKey: 'episodeNo',
        sorter: false,
        cell: (h, { row }) => row.episodeNo || '--'
    },
    {
        title: '是否定稿',
        colKey: 'isLocked',
        sorter: false,
        cell: (h, { row }) => {
            const option = isLockedStatusOptions.value.find(item => item.value === row.isLocked)
            if (option) {
                return (
                    <t-tag shape="round" theme={getStatusTagTheme(row.isLocked, isLockedStatusOptions.value)} variant="light">
                        {option.label}
                    </t-tag>
                );
            }
            return row.isLocked
        }
    },
    {
        title: '创建时间',
        colKey: 'created_at',
        width: 180,
        cell: (h, { row }) => formatDate(row.createdAt)
    },
    {
        title: '操作',
        colKey: 'action',
        width: 200,
        fixed: 'right',
        cell: (h, { row }) => h('t-space', { size: 'small' }, [
            h('t-button', {
                variant: 'text',
                size: 'small',
                style: {
                    margin: '8px',
                    cursor: 'pointer',
                    color: 'var(--td-brand-color)',
                    '--ripple-color': 'var(--td-brand-color)'
                },
                onClick: () => onView(row)
            }, '查看'),
            h('t-button', {
                variant: 'text',
                size: 'small',
                style: {
                    margin: '8px',
                    cursor: 'pointer',
                    color: 'var(--td-brand-color)',
                    '--ripple-color': 'var(--td-brand-color)'
                },
                onClick: () => onEdit(row)
            }, '编辑'),
            h('t-button', {
                variant: 'text',
                size: 'small',
                style: {
                    margin: '8px',
                    cursor: 'pointer',
                    color: 'var(--td-error-color)',
                    '--ripple-color': 'var(--td-error-color)'
                },
                onClick: () => onDelete(row)
            }, '删除')
        ])
    }
])

// 详情相关
const detailVisible = ref(false)
const detailData = ref({})
// 表单相关
const formRef = ref()
const drawerVisible = ref(false)
const drawerTitle = ref('')
const submitLoading = ref(false)
const formType = ref('create')

// 表单数据初始化，确保字符串字段为空字符串而不是null或undefined
const formData = ref({
    projectId: null,
    title: '',
    content: '',
    outline: '',
    episodeNo: null,
    isLocked: 0, // 状态字段默认第一个选项
})

// === 验证规则 ===
const rules = reactive({
    projectId: [
        { required: true, message: '请输入所属项目ID', trigger: ['blur', 'change'] },
        { type: 'number', message: '所属项目ID必须是数字', trigger: ['blur', 'change'] }

    ],
    title: [
        { whitespace: true, message: '分集标题不能只包含空格', trigger: 'blur' },
        { max: 255, message: '分集标题长度不能超过255个字符', trigger: ['blur', 'change'] },
        { min: 2, message: '分集标题长度不能少于2个字符', trigger: ['blur', 'change'] }

    ],
    content: [
        { max: 10000, message: '剧本正文长度不能超过10000个字符', trigger: ['blur', 'change'] }

    ],
    outline: [
        { whitespace: true, message: '大纲/简介不能只包含空格', trigger: 'blur' },
        { max: 255, message: '大纲/简介长度不能超过255个字符', trigger: ['blur', 'change'] }

    ],
    episodeNo: [
        { type: 'number', message: '第几集必须是数字', trigger: ['blur', 'change'] }

    ],
    isLocked: [

    ]
})

// 获取表格数据
const getTableData = async () => {
    loading.value = true
    try {
        const params = {
            page: pagination.current,
            pageSize: pagination.pageSize,
            ...processSearchParams()
        }
        const res = await getScriptsList(params)
        if (res.code === 0) {
            if (res.data && typeof res.data === 'object') {
                if (Array.isArray(res.data.list)) {
                    tableData.value = res.data.list
                    pagination.total = res.data.total || 0
                }
                else if (Array.isArray(res.data)) {
                    tableData.value = res.data
                    pagination.total = res.data.length
                }
                else {
                    tableData.value = []
                    pagination.total = 0
                }
            } else {
                tableData.value = []
                pagination.total = 0
            }
        } else {
            tableData.value = []
            pagination.total = 0
            MessagePlugin.error(res.message || '获取数据失败')
        }
    } catch (error) {
        console.error('获取数据失败:', error)
        tableData.value = []
        pagination.total = 0
        MessagePlugin.error('获取数据失败')
    } finally {
        loading.value = false
    }
}

// 处理搜索参数
const processSearchParams = () => {
    const params = {}
    // 所属项目ID
    if (searchInfo.value.projectId !== undefined && searchInfo.value.projectId !== '') {
        params.projectId = searchInfo.value.projectId
    }
    // 分集标题
    if (searchInfo.value.title !== undefined && searchInfo.value.title !== '') {
        params.title = searchInfo.value.title
    }
    // 第几集
    if (searchInfo.value.episodeNo !== undefined && searchInfo.value.episodeNo !== '') {
        params.episodeNo = searchInfo.value.episodeNo
    }
    // 是否定稿
    if (searchInfo.value.isLocked !== undefined && searchInfo.value.isLocked !== '') {
        params.isLocked = searchInfo.value.isLocked
    }

    return params
}

// 搜索
const onSearch = () => {
    pagination.current = 1
    getTableData()
}

// 重置搜索
const onReset = () => {
    searchInfo.value = {
        projectId: undefined,
        title: undefined,
        episodeNo: undefined,
        isLocked: undefined,
    }
    getTableData()
}

// 分页
const onPageChange = ({ current, pageSize }) => {
    pagination.pageSize = pageSize
    pagination.current = current
    getTableData()
}

const onPageSizeChange = ({ pageSize }) => {
    pagination.pageSize = pageSize
    pagination.current = 1
    getTableData()
}

// 查看详情
const onView = async (row) => {
    try {
        const res = await findScripts(row.id)
        if (res.code === 0) {
            let data = res.data

            detailData.value = data
            detailVisible.value = true
        } else {
            MessagePlugin.error(res.message || '获取数据失败')
        }
    } catch (error) {
        console.error('获取数据失败:', error)
        MessagePlugin.error('获取数据失败')
    }
}

const onRefresh = () => {
    MessagePlugin.loading('正在刷新数据...')
    getTableData().finally(() => {
        MessagePlugin.close()
        MessagePlugin.success('数据已刷新')
    })
}
// 新增
const onCreate = () => {
    formType.value = 'create'
    drawerTitle.value = '新增剧本'
    resetForm()
    getProjectsSelectData()
    drawerVisible.value = true
}

// 编辑
const onEdit = async (row) => {
    try {
        const res = await findScripts(row.id)
        if (res.code === 0) {
            formType.value = 'update'
            drawerTitle.value = '编辑剧本'

            // 处理返回的数据，确保上传字段格式正确
            const data = res.data

            // 确保所有字符串字段都有默认值，避免null或undefined导致的trim()错误
            if (data.title === null || data.title === undefined) {
                data.title = ''
            }
            if (data.outline === null || data.outline === undefined) {
                data.outline = ''
            }

            formData.value = data
            getProjectsSelectData()
            drawerVisible.value = true
        } else {
            MessagePlugin.error(res.message || '获取数据失败')
        }
    } catch (error) {
        console.error('获取数据失败:', error)
        MessagePlugin.error('获取数据失败')
    }
}

// 删除
const onDelete = async (row) => {
    const confirmDialog = DialogPlugin.confirm({
        header: '确认删除',
        body: '确定要删除这条数据吗？',
        onConfirm: async () => {
            try {
                const res = await deleteScripts(row.id)
                if (res.code === 0) {
                    MessagePlugin.success('删除成功')
                    getTableData()
                } else {
                    MessagePlugin.error(res.message || '删除失败')
                }
            } catch (error) {
                console.error('删除失败:', error)
                MessagePlugin.error('删除失败')
            }
            confirmDialog.destroy()
        }
    })
}

// === 表单提交方法 ===
const onSubmit = async () => {
    const valid = await formRef.value?.validate()

    if (valid !== true) {
        MessagePlugin.warning('请检查表单填写是否正确')
        return
    }

    submitLoading.value = true
    try {
        const submitData = { ...formData.value }
        // 处理上传字段数据

        // 处理空字符串字段，但不要将其设为null，以避免数据类型不匹配
        Object.keys(submitData).forEach(key => {
            if (typeof submitData[key] === 'string' && submitData[key].trim() === '') {
                // 对于字符串字段，保留空字符串而不是设为null
                submitData[key] = ''
            }
        })

        let res
        if (formType.value === 'create') {
            res = await createScripts(submitData)
        } else {
            res = await updateScripts(submitData.id, submitData)
        }

        if (res.code === 0) {
            MessagePlugin.success('操作成功')
            drawerVisible.value = false
            getTableData()
        } else {
            MessagePlugin.error(res.message || '操作失败')
        }
    } catch (error) {
        console.error('提交失败:', error)
        MessagePlugin.error('操作失败，请重试')
    } finally {
        submitLoading.value = false
    }
}

// 取消
const onCancel = () => {
    drawerVisible.value = false
    resetForm()
}

// === 重置表单时同时清除验证状态 ===
const resetForm = () => {
    formData.value = {
        projectId: null,
        title: '',
        content: '',
        outline: '',
        episodeNo: null,
        isLocked: 0, // 状态字段默认第一个选项
    }

    // 重置所有临时上传列表

    // 清除验证状态
    nextTick(() => {
        formRef.value?.clearValidate()
    })
}

// 初始化
const init = async () => {
    getProjectsSelectData()
    getTableData()
}

onMounted(() => {
    init()
})
</script>

<style scoped>
.scripts-list {
    padding: 20px;
}

.search-form {
    margin-bottom: 20px;
}

.richtext-editor-container {
    width: 100%;
}

:root {
    --td-brand-color: #0052D9;
    --td-error-color: #D54941;
    --td-brand-color-light: rgba(0, 82, 217, 0.05);
}

/* 操作按钮样式增强 */
:deep(.t-table .t-button--variant-text) {
    padding: 4px 8px;
    min-width: auto;
    font-size: 14px;
}

:deep(.t-table .t-button--variant-text:hover) {
    background-color: rgba(0, 82, 217, 0.05);
}

/* 删除按钮悬停效果 */
:deep(.t-table .t-button--variant-text[style*="--td-error-color"]:hover) {
    background-color: rgba(213, 73, 65, 0.05) !important;
}

/* 搜索表单样式 */
.search-form .t-form-item {
    margin-bottom: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .scripts-list {
        padding: 10px;
    }

    .search-form :deep(.t-col) {
        flex: 0 0 100% !important;
        max-width: 100% !important;
    }

    .uploaded-item {
        flex-direction: column;
        align-items: stretch;
    }

    .image-preview-wrapper {
        margin-right: 0;
        margin-bottom: 8px;
        align-self: center;
    }

    .uploaded-images {
        justify-content: center;
    }
}
</style>