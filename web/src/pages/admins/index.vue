
<template>
    <div class="admins-list">
        <!-- 搜索区域 -->
        <t-card class="search-form">
            <t-form
                ref="searchFormRef"
                :model="searchInfo"
                label-align="left"
                label-width="100px"
                @submit="onSearch"
            >
                <t-row :gutter="[24, 24]">
                            <!-- 用户名搜索 -->
                    <t-col :span="4">
                        <t-form-item label="用户名" name="username">
                            <t-input
                                v-model="searchInfo.username" placeholder="请输入用户名" clearable>
                            </t-input>
                        </t-form-item>
                    </t-col>
                            <!-- 手机号搜索 -->
                    <t-col :span="4">
                        <t-form-item label="手机号" name="mobile">
                            <t-input
                                v-model="searchInfo.mobile" placeholder="请输入手机号" clearable>
                            </t-input>
                        </t-form-item>
                    </t-col>
                            <!-- 邮箱搜索 -->
                    <t-col :span="4">
                        <t-form-item label="邮箱" name="email">
                            <t-input
                                v-model="searchInfo.email" placeholder="请输入邮箱" clearable>
                            </t-input>
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
    <t-dialog v-model:visible="drawerVisible" :header="drawerTitle" width="600px" size="medium" :confirm-btn="{
            content: '确定',
            theme: 'primary',
            loading: submitLoading
         }" @confirm="onSubmit" @cancel="onCancel">
        <t-form ref="formRef" :data="formData" :rules="rules" label-align="left" label-width="100px" @submit="onSubmit">
            <t-form-item label="用户名" name="username">
                <t-input v-model="formData.username" clearable placeholder="请输入用户名"
                    :status="!formData.username || (typeof formData.username === 'string' && formData.username.trim() === '') ? 'error' : 'default'" :maxlength="120" show-word-limit />
            </t-form-item>
            <t-form-item label="手机号" name="mobile">
                <t-input v-model="formData.mobile" clearable placeholder="请输入手机号"
                    :status="!formData.mobile || (typeof formData.mobile === 'string' && formData.mobile.trim() === '') ? 'error' : 'default'" :maxlength="11" show-word-limit />
            </t-form-item>
            <t-form-item label="密码" name="password">
                <t-input v-model="formData.password" type="password" clearable
                    :placeholder="formType === 'update' ? '请重新输入密码' : '请输入密码'" show-password
                    autocomplete="new-password">
                </t-input>
            </t-form-item>
            <t-form-item label="邮箱" name="email">
                <t-input v-model="formData.email" clearable placeholder="请输入邮箱" :maxlength="80" show-word-limit />
            </t-form-item>
        </t-form>
    </t-dialog>

        
<t-dialog v-model:visible="detailVisible" header="查看详情" width="600px" size="large" :footer="false" :close-btn="true"
                  :show-overlay="true" @close="detailVisible = false">
    <t-descriptions :column="1" layout="vertical" bordered
        :content-style="{ overflowWrap: 'break-word',whiteSpace:'normal' }">
        <t-descriptions-item label="用户名">
            <span v-if="detailData.username !== null && detailData.username !== undefined && detailData.username !== ''">
    {{ detailData.username }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="手机号">
            <span v-if="detailData.mobile !== null && detailData.mobile !== undefined && detailData.mobile !== ''">
    {{ detailData.mobile }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="密码">
            <span>******</span>
        </t-descriptions-item>
        <t-descriptions-item label="邮箱">
            <span v-if="detailData.email !== null && detailData.email !== undefined && detailData.email !== ''">
    {{ detailData.email }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
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
    createAdmins,
    deleteAdmins,
    updateAdmins,
    findAdmins,
    getAdminsList
} from '@/api/admins'
import { formatDate, getImageUrl } from '@/utils/format'

    defineOptions({
        name: 'AdminsList'
    })

    const router = useRouter()

    // ========== 状态选项定义 ==========

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
        username: undefined,
        mobile: undefined,
        email: undefined,
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

    // 表格列配置
    const columns = computed(() => [
        {
            title: '用户名',
            colKey: 'username',
            cell: (h, { row }) => row.username || '--'
        },
        {
            title: '手机号',
            colKey: 'mobile',
            cell: (h, { row }) => row.mobile || '--'
        },
        {
            title: '邮箱',
            colKey: 'email',
            cell: (h, { row }) => row.email || '--'
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
        username: '',
        mobile: '',
        password: '', // 密码字段
        email: '',
    })

    // === 验证规则 ===
    const rules = reactive({
        username: [
            { required: true, message: '请输入用户名', trigger: ['blur', 'change'] },
            { whitespace: true, message: '用户名不能只包含空格', trigger: 'blur' },
            { max: 120, message: '用户名长度不能超过120个字符', trigger: ['blur', 'change'] },
            { min: 2, message: '用户名长度不能少于2个字符', trigger: ['blur', 'change'] },
            { pattern: /^[a-zA-Z0-9_]{3,20}$/, message: '用户名只能包含字母、数字、下划线，长度3-20位', trigger: ['blur', 'change'] }

        ],
        mobile: [
            { required: true, message: '请输入手机号', trigger: ['blur', 'change'] },
            { whitespace: true, message: '手机号不能只包含空格', trigger: 'blur' },
            { max: 11, message: '手机号长度不能超过11个字符', trigger: ['blur', 'change'] },
            { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: ['blur', 'change'] }

        ],
        password: [
            // 密码字段验证规则 - 新增和编辑都必填
            {
                validator: (val) => {
                    // 无论新增还是编辑都必须填写密码
                    if (!val || val.trim() === '') {
                        return { result: false, message: '请输入密码', type: 'error' }
                    }
                    // 检查密码长度
                    if (val.length < 6) {
                        return { result: false, message: '密码长度不能少于6个字符', type: 'error' }
                    }
                    if (val.length > 50) {
                        return { result: false, message: '密码长度不能超过50个字符', type: 'error' }
                    }
                    return { result: true }
                },
                trigger: ['blur', 'change']
            }

        ],
        email: [
            { whitespace: true, message: '邮箱不能只包含空格', trigger: 'blur' },
            { max: 80, message: '邮箱长度不能超过80个字符', trigger: ['blur', 'change'] },
            { type: 'email', message: '请输入正确的邮箱格式', trigger: ['blur', 'change'] }

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
            const res = await getAdminsList(params)
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
        // 用户名
        if (searchInfo.value.username !== undefined && searchInfo.value.username !== '') {
            params.username = searchInfo.value.username
        }
        // 手机号
        if (searchInfo.value.mobile !== undefined && searchInfo.value.mobile !== '') {
            params.mobile = searchInfo.value.mobile
        }
        // 邮箱
        if (searchInfo.value.email !== undefined && searchInfo.value.email !== '') {
            params.email = searchInfo.value.email
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
            username: undefined,
            mobile: undefined,
            email: undefined,
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
            const res = await findAdmins(row.id)
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
        drawerTitle.value = '新增系统管理员'
        resetForm()
        drawerVisible.value = true
    }

    // 编辑
    const onEdit = async (row) => {
        try {
            const res = await findAdmins(row.id)
            if (res.code === 0) {
                formType.value = 'update'
                drawerTitle.value = '编辑系统管理员'

                // 处理返回的数据，确保上传字段格式正确
                const data = res.data
                // 密码字段编辑时不回显，保持为空
                data.password = ''

                // 确保所有字符串字段都有默认值，避免null或undefined导致的trim()错误
                if (data.username === null || data.username === undefined) {
                    data.username = ''
                }
                if (data.mobile === null || data.mobile === undefined) {
                    data.mobile = ''
                }
                if (data.email === null || data.email === undefined) {
                    data.email = ''
                }

                formData.value = data
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
                    const res = await deleteAdmins(row.id)
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
            // 密码字段处理：编辑时如果为空则删除该字段
            if (formType.value === 'update' && (!submitData.password || submitData.password.trim() === '')) {
                delete submitData.password
            }

            // 处理空字符串字段，但不要将其设为null，以避免数据类型不匹配
            Object.keys(submitData).forEach(key => {
                if (typeof submitData[key] === 'string' && submitData[key].trim() === '') {
                    // 对于字符串字段，保留空字符串而不是设为null
                    submitData[key] = ''
                }
            })

            let res
            if (formType.value === 'create') {
                res = await createAdmins(submitData)
            } else {
                res = await updateAdmins(submitData.id, submitData)
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
            username: '',
            mobile: '',
            password: '', // 密码字段
            email: '',
        }

        // 重置所有临时上传列表

        // 清除验证状态
        nextTick(() => {
            formRef.value?.clearValidate()
        })
    }

    // 初始化
    const init = async () => {
        getTableData()
    }

    onMounted(() => {
        init()
    })
</script>

<style scoped>
    .admins-list {
        padding: 20px;
    }

    .search-form {
        margin-bottom: 20px;
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
        .admins-list {
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