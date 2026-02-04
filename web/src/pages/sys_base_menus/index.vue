
<template>
  <div class="sys_base_menu-list">
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
          <!-- 父菜单搜索 -->
          <t-col :span="4">
            <t-form-item label="父菜单" name="parentId">
              <t-tree-select
                  v-model="searchInfo.parentId"
                  :data="sys_base_menusTreeData"
                  :keys="{ value: 'id', label: 'title', children: 'children' }"
                  placeholder="请选择父菜单"
                  clearable
                  :loading="sys_base_menusTreeLoading"
                  filterable
                  :filter="filterBase_menus"></t-tree-select>
            </t-form-item>
          </t-col>
          <!-- 路由名称搜索 -->
          <t-col :span="4">
            <t-form-item label="路由名称" name="name">
              <t-input
                  v-model="searchInfo.name"
                  placeholder="请输入路由名称"
                  clearable></t-input>
            </t-form-item>
          </t-col>
          <!-- 菜单标题搜索 -->
          <t-col :span="4"v-show="showAllQuery">
            <t-form-item label="菜单标题" name="title">
              <t-input
                  v-model="searchInfo.title"
                  placeholder="请输入菜单标题"
                  clearable></t-input>
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
              <t-button
                  variant="text"
                  @click="showAllQuery = !showAllQuery"
              >
                <template #icon>
                  <t-icon :name="showAllQuery ? 'chevron-up' : 'chevron-down'"></t-icon>
                </template>
                <span v-text="showAllQuery ? '收起' : '展开'"></span>
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
      <t-table
          ref="tableRef"
          :data="tableData"
          :columns="columns"
          :loading="loading"
          :pagination="pagination"
          row-key="id"
          @page-change="onPageChange"
          @page-size-change="onPageSizeChange"
          @sort-change="onSortChange"
          hover
      />
    </t-card>
    <!-- 新增/编辑抽屉 -->
    <t-drawer
        v-model:visible="drawerVisible"
        :header="drawerTitle"
        size="medium"
        :confirm-btn="{
                content: '确定',
                theme: 'primary',
                loading: submitLoading
            }"
        @confirm="onSubmit"
        @cancel="onCancel"
    >
      <t-form
          ref="formRef"
          :data="formData"
          :rules="rules"
          label-align="left"
          label-width="100px"
          @submit="onSubmit"
      >
        <t-form-item label="父菜单ID" name="parentId">
          <t-tree-select
              v-model="formData.parentId"
              :data="sys_base_menusTreeData"
              :keys="{ value: 'id', label: 'title', children: 'children' }"
              placeholder="请选择父菜单ID"
              clearable
              :loading="sys_base_menusTreeLoading"
              filterable
              :filter="filterBase_menus"
              empty="暂无数据"
          >
            <template #empty>
              <div style="text-align: center; padding: 20px; color: var(--td-text-color-placeholder);">
                <t-icon name="folder" size="24px" style="margin-bottom: 8px; display: block;"></t-icon>
                暂无数据
              </div>
            </template>
          </t-tree-select>
        </t-form-item>
        <t-form-item label="路由路径" name="path">
          <t-input
              v-model="formData.path"
              clearable
              placeholder="请输入路由路径"
              :status="!formData.path || (typeof formData.path === 'string' && formData.path.trim() === '') ? 'error' : 'default'"
              :maxlength="255"
              show-word-limit
          />
        </t-form-item>
        <t-form-item label="路由名称" name="name">
          <t-input
              v-model="formData.name"
              clearable
              placeholder="请输入路由名称"
              :status="!formData.name || (typeof formData.name === 'string' && formData.name.trim() === '') ? 'error' : 'default'"
              :maxlength="255"
              show-word-limit
          />
        </t-form-item>
        <t-form-item label="是否隐藏" name="hidden">
          <t-select
              v-model="formData.hidden"
              placeholder="请选择是否隐藏"
              clearable
          >
            <t-option :value="0" label="显示" />
            <t-option :value="1" label="隐藏" />
          </t-select>
        </t-form-item>
        <t-form-item label="组件路径" name="component">
          <t-input
              v-model="formData.component"
              clearable
              placeholder="请输入组件路径"
              :status="!formData.component || (typeof formData.component === 'string' && formData.component.trim() === '') ? 'error' : 'default'"
              :maxlength="255"
              show-word-limit
          />
        </t-form-item>
        <t-form-item label="排序" name="sort"><t-input
            v-model="formData.sort"
            placeholder="请输入排序"
            type="number"
            clearable
            :status="formData.sort === null || formData.sort === undefined ? 'error' : 'default'"
        />
        </t-form-item>
        <t-form-item label="菜单标题" name="title">
          <t-input
              v-model="formData.title"
              clearable
              placeholder="请输入菜单标题"
              :status="!formData.title || (typeof formData.title === 'string' && formData.title.trim() === '') ? 'error' : 'default'"
              :maxlength="255"
              show-word-limit
          />
        </t-form-item>
        <t-form-item label="菜单图标" name="icon">
          <icon v-model="formData.icon" />
        </t-form-item>
      </t-form>
    </t-drawer>

    <!-- 详情抽屉 -->
    <t-drawer
        v-model:visible="detailVisible"
        header="查看详情"
        size="medium"
    >
      <t-descriptions :column="1" layout="vertical" bordered :content-style="{ overflowWrap: 'break-word',whiteSpace:'normal' }">
        <t-descriptions-item label="父菜单ID">
          <div v-if="detailData.parent">
            <div style="font-weight: 500;">
              {{ detailData.parent.name }}
            </div>
            <div style="font-size: 12px; color: var(--td-text-color-placeholder); margin-top: 4px;">
              ID: {{ detailData.parent.id }}
            </div>
          </div>
          <span v-else style="color: var(--td-text-color-placeholder);">未选择</span>
        </t-descriptions-item>
        <t-descriptions-item label="路由路径">
          <span v-text="detailData.path"></span>
        </t-descriptions-item>
        <t-descriptions-item label="路由名称">
          <span v-text="detailData.name"></span>
        </t-descriptions-item>
        <t-descriptions-item label="是否隐藏">
          <t-tag
              :theme="Number(detailData.hidden) === 1 ? 'warning' : 'success'"
              variant="light"
              size="small"
          >
            {{ Number(detailData.hidden) === 1 ? '隐藏' : '显示' }}
          </t-tag>
        </t-descriptions-item>
        <t-descriptions-item label="组件路径">
          <span v-text="detailData.component"></span>
        </t-descriptions-item>
        <t-descriptions-item label="排序">
          <span v-text="detailData.sort"></span>
        </t-descriptions-item>
        <t-descriptions-item label="菜单标题">
          <span v-text="detailData.title"></span>
        </t-descriptions-item>
        <t-descriptions-item label="菜单图标">
          <span v-text="detailData.icon"></span>
        </t-descriptions-item>
      </t-descriptions>
    </t-drawer>
  </div>
</template>

<script setup lang="tsx">
import {ref, reactive, computed, onMounted, nextTick} from 'vue'
import {useRouter} from 'vue-router'
import {MessagePlugin, DialogPlugin} from 'tdesign-vue-next'
import {
  createSysBaseMenus,
  deleteSysBaseMenus,
  updateSysBaseMenus,
  findSysBaseMenus,
  getSysBaseMenusList,
  getSysBaseMenusTreeList
} from '@/api/sys_base_menus'
import icon from '@/pages/sys_base_menus/icon.vue'
import {formatDate} from '@/utils/format'

defineOptions({
  name: 'SysBaseMenuList'
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
  parentId: undefined,
  name: undefined,
  title: undefined,
})

// 计算是否有搜索项和隐藏搜索项
const hasSearchItems = computed(() => {
  return true
})

const hasHiddenSearchItems = computed(() => {
  return true
})

// 系统菜单(父级)树形数据相关
const sys_base_menusTreeData = ref([])
const sys_base_menusTreeLoading = ref(false)

// 获取系统菜单(父级)树形数据
const getSys_base_menusTreeData = async () => {
  sys_base_menusTreeLoading.value = true
  try {
    const res = await getSysBaseMenusTreeList()
    if (res.code === 0) {
      const treeData = res.data || []

      // 在树形数据前面添加"无父级"选项
      sys_base_menusTreeData.value = [
        {
          id: 0,
          title: '无父级（顶级菜单）',
          name: '无父级（顶级菜单）',
          children: []
        },
        ...treeData
      ]
    } else {
      sys_base_menusTreeData.value = [
        {
          id: 0,
          title: '无父级（顶级菜单）',
          name: '无父级（顶级菜单）',
          children: []
        }
      ]
      console.error('获取系统菜单(父级)数据失败:', res.message)
    }
  } catch (error) {
    console.error('获取系统菜单(父级)数据失败:', error)
    sys_base_menusTreeData.value = [
      {
        id: 0,
        title: '无父级（顶级菜单）',
        name: '无父级（顶级菜单）',
        children: []
      }
    ]
  } finally {
    sys_base_menusTreeLoading.value = false
  }
}

// 系统菜单(父级)过滤方法
const filterBase_menus = (inputValue, node) => {
  return node.name.toLowerCase().includes(inputValue.toLowerCase())
}

// 表格列配置
const columns = computed(() => [
  {
    title: '父菜单ID',
    colKey: 'parentId',
    cell: (h, {row}) => {
      const relationObj = row.parent
      if (relationObj && relationObj.name) {
        return h('t-space', {direction: 'vertical', size: 'small'}, [
          h('span', {style: {fontWeight: '500'}}, relationObj.name),
          h('span', {style: {fontSize: '12px', color: 'var(--td-text-color-placeholder)'}}, `ID: ${relationObj.id}`)
        ])
      }
      return h('span', {style: {color: 'var(--td-text-color-placeholder)'}}, '未选择')
    }
  },
  {
    title: '路由路径',
    colKey: 'path',
    cell: (h, {row}) => row.path || '--'
  },
  {
    title: '路由名称',
    colKey: 'name',
    cell: (h, {row}) => row.name || '--'
  },
  {
    title: '是否隐藏',
    colKey: 'hidden',
    cell: (h, {row}) => {
      const isHidden = Number(row.hidden) === 1
      return h('t-tag', {
        theme: isHidden ? 'warning' : 'success',
        variant: 'light',
        size: 'small'
      }, isHidden ? '隐藏' : '显示')
    }
  },
  {
    title: '组件路径',
    colKey: 'component',
    cell: (h, {row}) => row.component || '--'
  },
  {
    title: '排序',
    colKey: 'sort',
    cell: (h, {row}) => row.sort || '--'
  },
  {
    title: '菜单标题',
    colKey: 'title',
    cell: (h, {row}) => row.title || '--'
  },
  {
    title: '菜单图标',
    colKey: 'icon',
    cell: (h, {row}) => row.icon || '--'
  },
  {
    title: '创建时间',
    colKey: 'createdAt',
    width: 180,
    sorter: false,
    cell: (h, {row}) => formatDate(row.createdAt)
  },
  {
    title: '操作',
    colKey: 'action',
    width: 200,
    fixed: 'right',
    cell: (h, {row}) => h('t-space', {size: 'small'}, [
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
  parentID: undefined,
  parentId: 0,
  path: '',
  name: '',
  hidden: 0,
  component: '',
  sort: 0,
  title: '',
  icon: '',
})

// === 验证规则 ===
const rules = reactive({
  parentId: [
    // 非必填字段的基础验证
    {type: 'number', message: '父菜单ID必须是数字', trigger: ['blur', 'change']}
  ],
  path: [
    {required: true, message: '请输入路由路径', trigger: ['blur', 'change']},
    {whitespace: true, message: '路由路径不能只包含空格', trigger: 'blur'}, {
      max: 255,
      message: '路由路径长度不能超过255个字符',
      trigger: ['blur', 'change']
    }

  ],
  name: [
    {required: true, message: '请输入路由名称', trigger: ['blur', 'change']},
    {whitespace: true, message: '路由名称不能只包含空格', trigger: 'blur'}, {
      max: 255,
      message: '路由名称长度不能超过255个字符',
      trigger: ['blur', 'change']
    }
    , {min: 2, message: '路由名称长度不能少于2个字符', trigger: ['blur', 'change']}

  ],
  hidden: [
    {required: true, message: '请输入是否隐藏', trigger: ['blur', 'change']}
  ],
  component: [
    {required: true, message: '请输入组件路径', trigger: ['blur', 'change']},
    {whitespace: true, message: '组件路径不能只包含空格', trigger: 'blur'}, {
      max: 255,
      message: '组件路径长度不能超过255个字符',
      trigger: ['blur', 'change']
    }

  ],
  sort: [
    {required: true, message: '请输入排序', trigger: ['blur', 'change']},
    {type: 'number', message: '排序必须是数字', trigger: ['blur', 'change']}
    , {min: 0, max: 9999, message: '排序值必须在0-9999之间', trigger: ['blur', 'change']}

  ],
  title: [
    {required: true, message: '请输入菜单标题', trigger: ['blur', 'change']},
    {whitespace: true, message: '菜单标题不能只包含空格', trigger: 'blur'}, {
      max: 255,
      message: '菜单标题长度不能超过255个字符',
      trigger: ['blur', 'change']
    }
    , {min: 2, message: '菜单标题长度不能少于2个字符', trigger: ['blur', 'change']}

  ],
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
    const res = await getSysBaseMenusList(params)
    if (res.code === 0) {
      if (res.data && typeof res.data === 'object') {
        if (Array.isArray(res.data.list)) {
          tableData.value = res.data.list
          pagination.total = res.data.total || 0
        } else if (Array.isArray(res.data)) {
          tableData.value = res.data
          pagination.total = res.data.length
        } else {
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
  if (searchInfo.value.parentId !== undefined && searchInfo.value.parentId !== '') {
    params.parentId = searchInfo.value.parentId
  }

  if (searchInfo.value.name !== undefined && searchInfo.value.name !== '') {
    params.name = searchInfo.value.name
  }

  if (searchInfo.value.title !== undefined && searchInfo.value.title !== '') {
    params.title = searchInfo.value.title
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
    parentId: undefined,
    name: undefined,
    title: undefined,
  }
  getTableData()
}

// 分页
const onPageChange = ({current, pageSize}) => {
  pagination.pageSize = pageSize
  pagination.current = current
  getTableData()
}

const onPageSizeChange = ({pageSize}) => {
  pagination.pageSize = pageSize
  pagination.current = 1
  getTableData()
}

// 查看详情
const onView = async (row) => {
  try {
    const res = await findSysBaseMenus(row.id)
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
    MessagePlugin.close(Promise.resolve(undefined))
    MessagePlugin.success('数据已刷新')
  })
}
// 新增
const onCreate = () => {
  formType.value = 'create'
  drawerTitle.value = '新增系统菜单'
  resetForm()
  getSys_base_menusTreeData()
  drawerVisible.value = true
}

// 编辑
const onEdit = async (row) => {
  try {
    const res = await findSysBaseMenus(row.id)
    if (res.code === 0) {
      formType.value = 'update'
      drawerTitle.value = '编辑系统菜单'

      // 处理返回的数据，确保上传字段格式正确
      const data = res.data

      // 确保所有字符串字段都有默认值，避免null或undefined导致的trim()错误
      if (data.path === null || data.path === undefined) {
        data.path = ''
      }
      if (data.name === null || data.name === undefined) {
        data.name = ''
      }
      if (data.component === null || data.component === undefined) {
        data.component = ''
      }

      if (data.title === null || data.title === undefined) {
        data.title = ''
      }
      if (data.icon === null || data.icon === undefined) {
        data.icon = ''
      }

      formData.value = data
      getSys_base_menusTreeData()
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
        const res = await deleteSysBaseMenus(row.id)
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
    const submitData = {...formData.value}

    // 处理空字符串字段，但不要将其设为null，以避免数据类型不匹配
    Object.keys(submitData).forEach(key => {
      if (typeof submitData[key] === 'string' && submitData[key].trim() === '') {
        // 对于字符串字段，保留空字符串而不是设为null
        submitData[key] = ''
      }
    })

    let res
    if (formType.value === 'create') {
      res = await createSysBaseMenus(submitData)
    } else {
      res = await updateSysBaseMenus(submitData.id, submitData)
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
    parentID: undefined,
    parentId: 0,
    path: '',
    name: '',
    hidden: 0,
    component: '',
    sort: 0,
    title: '',
    icon: '',
  }

  // 重置所有临时上传列表

  // 清除验证状态
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

// 初始化
const init = async () => {
  getSys_base_menusTreeData()
  getTableData()
}

onMounted(() => {
  init()
})
</script>

<style scoped>
.sys_base_menu-list {
  padding: 20px;
}

.search-form {
  margin-bottom: 20px;
}

.richtext-content h1,
.richtext-content h2,
.richtext-content h3,
.richtext-content h4,
.richtext-content h5,
.richtext-content h6 {
  margin-top: 0.5em;
  margin-bottom: 0.5em;
}

.richtext-content p {
  margin-bottom: 1em;
}

.richtext-content img {
  max-width: 100%;
  height: auto;
}

.richtext-content ul,
.richtext-content ol {
  padding-left: 2em;
  margin-bottom: 1em;
}


.image-preview-wrapper:hover .image-overlay {
  opacity: 1;
}

/* 上传区域样式 */
.upload-area :deep(.t-upload__trigger) {
  width: 100%;
  min-height: 120px;
  border: 2px dashed var(--td-border-level-2-color);
  border-radius: 8px;
  background: var(--td-bg-color-container);
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-area :deep(.t-upload__trigger:hover) {
  border-color: var(--td-brand-color);
  background: var(--td-brand-color-light);
}

.upload-trigger .t-icon {
  color: var(--td-text-color-placeholder);
  margin-bottom: 12px;
  display: block;
}

/* 多图上传按钮 */
.multi-upload-btn :deep(.t-upload__trigger) {
  width: 100px;
  height: 100px;
  border: 2px dashed var(--td-border-level-2-color);
  border-radius: 8px;
  background: var(--td-bg-color-container);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.multi-upload-btn :deep(.t-upload__trigger:hover) {
  border-color: var(--td-brand-color);
  background: var(--td-brand-color-light);
}

.reupload-component :deep(.t-upload__trigger) {
  width: auto;
  border: none;
  background: none;
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
  .sys_base_menu-list {
    padding: 10px;
  }

  .search-form :deep(.t-col) {
    flex: 0 0 100% !important;
    max-width: 100% !important;
  }
}
</style>