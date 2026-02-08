<template>
    <div class="project-detail-container">
        <div class="detail-header">
            <div class="header-left">
                <t-button variant="text" shape="circle" @click="goBack">
                    <template #icon><t-icon name="arrow-left" /></template>
                </t-button>
                <div class="project-info">
                    <div class="title-row">
                        <span class="title">{{ project.title || '加载中...' }}</span>
                        <t-tag :theme="getStatusTheme(project.status)" variant="light">{{ getStatusText(project.status)
                        }}</t-tag>
                    </div>
                    <div class="desc">创建时间: {{ formatDate(project.createdAt) }}</div>
                </div>
            </div>
            <div class="header-right">
                <t-button theme="default" variant="outline" @click="init"><template #icon><t-icon
                            name="refresh" /></template>刷新数据</t-button>
            </div>
        </div>

        <div class="detail-content" v-loading="loading">
            <t-tabs v-model="activeTab" theme="card" class="detail-tabs">

                <t-tab-panel value="overview" label="项目概览">
                    <div class="tab-panel-content overview-panel">
                        <t-row :gutter="24">
                            <t-col :span="4">
                                <div class="cover-section">
                                    <t-image :src="getImageUrl(project.image)" fit="cover" class="project-cover-large"
                                        :style="{ aspectRatio: getAspectRatio(project.settings) }" />
                                </div>
                            </t-col>
                            <t-col :span="8">
                                <t-card title="基本信息" :bordered="false" class="info-card">
                                    <t-descriptions :column="1" layout="vertical">
                                        <t-descriptions-item label="项目名称">{{ project.title }}</t-descriptions-item>
                                        <t-descriptions-item label="剧情简介">{{ project.description || '暂无简介'
                                        }}</t-descriptions-item>
                                        <t-descriptions-item label="视频比例">{{ getRatioLabel(project.settings)
                                        }}</t-descriptions-item>
                                        <t-descriptions-item label="总时长">{{ formatDuration(project.totalDuration)
                                        }}</t-descriptions-item>
                                    </t-descriptions>
                                </t-card>

                                <t-card title="数据统计" :bordered="false" class="info-card" style="margin-top: 16px">
                                    <t-row :gutter="16">
                                        <t-col :span="4">
                                            <div class="stat-item">
                                                <div class="label">剧集数</div>
                                                <div class="num">{{ episodeList.length }}</div>
                                            </div>
                                        </t-col>
                                        <t-col :span="4">
                                            <div class="stat-item">
                                                <div class="label">角色数</div>
                                                <div class="num">{{ characterList.length }}</div>
                                            </div>
                                        </t-col>
                                        <t-col :span="4">
                                            <div class="stat-item">
                                                <div class="label">场景数</div>
                                                <div class="num">{{ sceneList.length }}</div>
                                            </div>
                                        </t-col>
                                    </t-row>
                                </t-card>
                            </t-col>
                        </t-row>
                    </div>
                </t-tab-panel>

                <t-tab-panel value="episodes" label="剧集列表">
                    <div class="tab-panel-content">
                        <div class="action-bar">
                            <t-button theme="primary" @click="handleAddEpisode"><template #icon><t-icon
                                        name="add" /></template>新建剧集</t-button>
                        </div>
                        <div class="episode-grid">
                            <t-card v-for="ep in episodeList" :key="ep.id" class="episode-card" hover-shadow
                                @click="enterEpisode(ep)">
                                <div class="ep-cover">
                                    <div class="ep-no">第 {{ ep.episodeNo }} 集</div>
                                    <div class="ep-status" v-if="ep.status === 2"><t-icon name="check-circle" /> 已生成
                                    </div>
                                </div>
                                <div class="ep-info">
                                    <div class="ep-title">{{ ep.title }}</div>
                                    <div class="ep-meta">分镜数: {{ ep.shotsCount || 0 }}</div>
                                </div>
                            </t-card>
                        </div>
                        <t-empty v-if="episodeList.length === 0" description="暂无剧集，点击新建开始创作" />
                    </div>
                </t-tab-panel>

                <t-tab-panel value="characters" label="角色库">
                    <div class="tab-panel-content">
                        <div class="action-bar">
                            <div class="left-actions">
                                <t-button theme="primary" @click="openCharacterDialog('create')"><template #icon><t-icon
                                            name="user-add" /></template>添加角色</t-button>
                                <t-button theme="default" variant="outline" :disabled="selectedCharIds.length === 0"
                                    @click="batchGenerate('char')" :loading="batchGeneratingChar">
                                    <template #icon><t-icon name="magic" /></template>批量生图
                                </t-button>
                            </div>

                            <t-input v-model="charSearch" placeholder="搜索角色..." style="width: 200px"
                                @enter="loadCharacters">
                                <template #suffix-icon><t-icon name="search" @click="loadCharacters" /></template>
                            </t-input>
                        </div>

                        <div class="selection-bar" v-if="characterList.length > 0">
                            <t-checkbox :checked="checkAllChars" :indeterminate="isCharIndeterminate"
                                @change="handleSelectAllChars">全选 ({{ selectedCharIds.length }})</t-checkbox>
                        </div>

                        <div class="resource-grid">
                            <t-card v-for="char in characterList" :key="char.id" class="resource-card"
                                :class="{ 'is-selected': selectedCharIds.includes(char.id) }" :bordered="false">
                                <div class="card-select">
                                    <t-checkbox :checked="selectedCharIds.includes(char.id)"
                                        @change="() => toggleCharSelection(char.id)" />
                                </div>

                                <div class="res-cover">
                                    <t-image :src="getImageUrl(char.avatarUrl)" fit="cover" class="res-img" />

                                    <div v-if="generatingCharIds.includes(char.id)" class="loading-mask">
                                        <t-loading size="small" text="生成中..." />
                                    </div>

                                    <div class="res-overlay">
                                        <t-tooltip content="AI生成形象">
                                            <t-button shape="circle" theme="success" size="small"
                                                @click="singleGenerate('char', char)">
                                                <t-icon name="magic" />
                                            </t-button>
                                        </t-tooltip>

                                        <t-button shape="circle" theme="primary" size="small"
                                            @click="openCharacterDialog('edit', char)"><t-icon name="edit" /></t-button>
                                        <t-popconfirm content="确认删除该角色？" @confirm="handleDeleteCharacter(char.id)">
                                            <t-button shape="circle" theme="danger" size="small"><t-icon
                                                    name="delete" /></t-button>
                                        </t-popconfirm>
                                    </div>
                                </div>
                                <div class="res-info">
                                    <div class="res-name">{{ char.name }}</div>
                                    <t-tag size="small" :theme="char.roleType === 'main' ? 'primary' : 'default'"
                                        variant="light">{{ char.roleType === 'main' ? '主角' : '配角' }}</t-tag>
                                </div>
                                <div class="res-desc text-ellipsis-2">{{ char.appearanceDesc || '暂无外貌描述' }}</div>
                            </t-card>
                        </div>
                        <t-pagination v-if="charPagination.total > 0" v-model="charPagination.current"
                            v-model:pageSize="charPagination.pageSize" :total="charPagination.total"
                            @change="loadCharacters" style="margin-top: 20px" />
                        <t-empty v-if="characterList.length === 0" description="暂无角色数据" />
                    </div>
                </t-tab-panel>

                <t-tab-panel value="scenes" label="场景库">
                    <div class="tab-panel-content">
                        <div class="action-bar">
                            <div class="left-actions">
                                <t-button theme="primary" @click="openSceneDialog('create')"><template #icon><t-icon
                                            name="image" /></template>添加场景</t-button>
                                <t-button theme="default" variant="outline" :disabled="selectedSceneIds.length === 0"
                                    @click="batchGenerate('scene')" :loading="batchGeneratingScene">
                                    <template #icon><t-icon name="magic" /></template>批量生图
                                </t-button>
                            </div>

                            <t-input v-model="sceneSearch" placeholder="搜索场景..." style="width: 200px"
                                @enter="loadScenes">
                                <template #suffix-icon><t-icon name="search" @click="loadScenes" /></template>
                            </t-input>
                        </div>

                        <div class="selection-bar" v-if="sceneList.length > 0">
                            <t-checkbox :checked="checkAllScenes" :indeterminate="isSceneIndeterminate"
                                @change="handleSelectAllScenes">全选 ({{ selectedSceneIds.length }})</t-checkbox>
                        </div>

                        <div class="resource-grid">
                            <t-card v-for="scene in sceneList" :key="scene.id" class="resource-card"
                                :class="{ 'is-selected': selectedSceneIds.includes(scene.id) }" :bordered="false">
                                <div class="card-select">
                                    <t-checkbox :checked="selectedSceneIds.includes(scene.id)"
                                        @change="() => toggleSceneSelection(scene.id)" />
                                </div>

                                <div class="res-cover scene-cover">
                                    <t-image :src="getImageUrl(scene.visualPrompt)" fit="cover" class="res-img" />

                                    <div v-if="generatingSceneIds.includes(scene.id)" class="loading-mask">
                                        <t-loading size="small" text="生成中..." />
                                    </div>

                                    <div class="res-overlay">
                                        <t-tooltip content="AI生成图片">
                                            <t-button shape="circle" theme="success" size="small"
                                                @click="singleGenerate('scene', scene)">
                                                <t-icon name="magic" />
                                            </t-button>
                                        </t-tooltip>

                                        <t-button shape="circle" theme="primary" size="small"
                                            @click="openSceneDialog('edit', scene)"><t-icon name="edit" /></t-button>
                                        <t-popconfirm content="确认删除该场景？" @confirm="handleDeleteScene(scene.id)">
                                            <t-button shape="circle" theme="danger" size="small"><t-icon
                                                    name="delete" /></t-button>
                                        </t-popconfirm>
                                    </div>
                                </div>
                                <div class="res-info">
                                    <div class="res-name">{{ scene.name }}</div>
                                    <span class="scene-meta">{{ scene.time }} · {{ scene.location }}</span>
                                </div>
                                <div class="res-desc text-ellipsis-2">{{ scene.atmosphere || '暂无描述' }}</div>
                            </t-card>
                        </div>
                        <t-pagination v-if="scenePagination.total > 0" v-model="scenePagination.current"
                            v-model:pageSize="scenePagination.pageSize" :total="scenePagination.total"
                            @change="loadScenes" style="margin-top: 20px" />
                        <t-empty v-if="sceneList.length === 0" description="暂无场景数据" />
                    </div>
                </t-tab-panel>

                <t-tab-panel value="props" label="道具库">
                    <div class="tab-panel-content">
                        <div class="action-bar">
                            <t-button theme="primary" @click="openPropDialog('create')"><template #icon><t-icon
                                        name="gift" /></template>添加道具</t-button>
                            <t-input v-model="propSearch" placeholder="搜索道具..." style="width: 200px" @enter="loadProps">
                                <template #suffix-icon><t-icon name="search" @click="loadProps" /></template>
                            </t-input>
                        </div>

                        <div class="resource-grid">
                            <t-card v-for="prop in propList" :key="prop.id" class="resource-card" :bordered="false">
                                <div class="res-cover">
                                    <t-image :src="getImageUrl(prop.imageUrl)" fit="contain" class="res-img"
                                        style="padding: 10px; background: #f9f9f9;" />
                                    <div class="res-overlay">
                                        <t-button shape="circle" theme="primary" size="small"
                                            @click="openPropDialog('edit', prop)"><t-icon name="edit" /></t-button>
                                        <t-popconfirm content="确认删除该道具？" @confirm="handleDeleteProp(prop.id)">
                                            <t-button shape="circle" theme="danger" size="small"><t-icon
                                                    name="delete" /></t-button>
                                        </t-popconfirm>
                                    </div>
                                </div>
                                <div class="res-info">
                                    <div class="res-name">{{ prop.name }}</div>
                                    <t-tag size="small" variant="outline">{{ prop.type || '通用' }}</t-tag>
                                </div>
                                <div class="res-desc text-ellipsis-2">{{ prop.description || '暂无描述' }}</div>
                            </t-card>
                        </div>
                        <t-pagination v-if="propPagination.total > 0" v-model="propPagination.current"
                            v-model:pageSize="propPagination.pageSize" :total="propPagination.total" @change="loadProps"
                            style="margin-top: 20px" />
                        <t-empty v-if="propList.length === 0" description="暂无道具数据" />
                    </div>
                </t-tab-panel>

            </t-tabs>
        </div>

        <t-dialog v-model:visible="charDialog.visible" :header="charDialog.mode === 'create' ? '添加角色' : '编辑角色'"
            :confirm-btn="{ loading: charDialog.loading }" @confirm="submitCharacter">
            <t-form :data="charFormData" label-align="top">
                <t-form-item label="头像" name="avatarUrl">
                    <t-upload v-model="charFileList" :action="uploadConfig.action" :headers="uploadConfig.headers"
                        :show-file-list="false" accept="image/*" @success="(ctx) => handleUploadSuccess(ctx, 'char')">
                        <div class="upload-box" v-if="!charFormData.avatarUrl"><t-icon name="add" /></div>
                        <t-image v-else :src="getImageUrl(charFormData.avatarUrl)" class="upload-preview" fit="cover" />
                    </t-upload>
                </t-form-item>
                <t-form-item label="名称" name="name" required><t-input v-model="charFormData.name" /></t-form-item>
                <t-form-item label="类型" name="roleType"><t-select v-model="charFormData.roleType"
                        :options="[{ label: '主角', value: 'main' }, { label: '配角', value: 'supporting' }]" /></t-form-item>
                <t-form-item label="外貌描述" name="appearanceDesc"><t-textarea
                        v-model="charFormData.appearanceDesc" /></t-form-item>
            </t-form>
        </t-dialog>

        <t-dialog v-model:visible="sceneDialog.visible" :header="sceneDialog.mode === 'create' ? '添加场景' : '编辑场景'"
            :confirm-btn="{ loading: sceneDialog.loading }" @confirm="submitScene">
            <t-form :data="sceneFormData" label-align="top">
                <t-form-item label="参考图" name="visualPrompt">
                    <t-upload v-model="sceneFileList" :action="uploadConfig.action" :headers="uploadConfig.headers"
                        :show-file-list="false" accept="image/*" @success="(ctx) => handleUploadSuccess(ctx, 'scene')">
                        <div class="upload-box scene" v-if="!sceneFormData.visualPrompt"><t-icon name="add" /></div>
                        <t-image v-else :src="getImageUrl(sceneFormData.visualPrompt)" class="upload-preview scene"
                            fit="cover" />
                    </t-upload>
                </t-form-item>
                <t-form-item label="名称" name="name" required><t-input v-model="sceneFormData.name"
                        placeholder="例如：街道-白天" /></t-form-item>
                <t-row :gutter="16">
                    <t-col :span="6"><t-form-item label="地点" name="location"><t-input
                                v-model="sceneFormData.location" /></t-form-item></t-col>
                    <t-col :span="6"><t-form-item label="时间" name="time"><t-input
                                v-model="sceneFormData.time" /></t-form-item></t-col>
                </t-row>
                <t-form-item label="氛围描述" name="atmosphere"><t-textarea
                        v-model="sceneFormData.atmosphere" /></t-form-item>
            </t-form>
        </t-dialog>

        <t-dialog v-model:visible="propDialog.visible" :header="propDialog.mode === 'create' ? '添加道具' : '编辑道具'"
            :confirm-btn="{ loading: propDialog.loading }" @confirm="submitProp">
            <t-form :data="propFormData" label-align="top">
                <t-form-item label="图片" name="imageUrl">
                    <t-upload v-model="propFileList" :action="uploadConfig.action" :headers="uploadConfig.headers"
                        :show-file-list="false" accept="image/*" @success="(ctx) => handleUploadSuccess(ctx, 'prop')">
                        <div class="upload-box" v-if="!propFormData.imageUrl"><t-icon name="add" /></div>
                        <t-image v-else :src="getImageUrl(propFormData.imageUrl)" class="upload-preview" fit="contain"
                            style="background:#f9f9f9" />
                    </t-upload>
                </t-form-item>
                <t-form-item label="名称" name="name" required><t-input v-model="propFormData.name" /></t-form-item>
                <t-form-item label="类型" name="type"><t-input v-model="propFormData.type"
                        placeholder="例如：交通工具、武器" /></t-form-item>
                <t-form-item label="AI提示词 (Image Prompt)" name="imagePrompt"><t-textarea
                        v-model="propFormData.imagePrompt" placeholder="用于AI生成的英文提示词" /></t-form-item>
                <t-form-item label="描述" name="description"><t-textarea
                        v-model="propFormData.description" /></t-form-item>
            </t-form>
        </t-dialog>

    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import {
    ArrowLeftIcon, AddIcon, SearchIcon, EditIcon, DeleteIcon,
    RefreshIcon, UserAddIcon, ImageIcon, GiftIcon, CheckCircleIcon, MagicIcon
} from 'tdesign-icons-vue-next'
import dayjs from 'dayjs'

// API
import { findProjects } from '@/api/projects'
import { getScriptsList, createScripts } from '@/api/scripts'
import { getCharactersList, createCharacters, updateCharacters, deleteCharacters } from '@/api/characters'
import { getScenesList, createScenes, updateScenes, deleteScenes } from '@/api/scenes'
import { getPropsList, createProps, updateProps, deleteProps } from '@/api/props'
// 任务API
import {
    generateCharacterImageTask,
    batchGenerateCharacterImagesTask,
    generateSceneImageTask,
    batchGenerateSceneImagesTask,
    findTasks
} from '@/api/tasks'

import { getImageUrl } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const projectId = route.params.id as string

// === 状态 ===
const loading = ref(false)
const project = ref<any>({})
const activeTab = ref('overview')

// 剧集
const episodeList = ref<any[]>([])

// 角色
const characterList = ref<any[]>([])
const charSearch = ref('')
const charPagination = reactive({ current: 1, pageSize: 12, total: 0 })
const charDialog = reactive({ visible: false, mode: 'create', loading: false })
const charFormData = ref<any>({})
const charFileList = ref([])
// 角色生图状态
const selectedCharIds = ref<number[]>([])
const generatingCharIds = ref<number[]>([])
const batchGeneratingChar = ref(false)

// 场景
const sceneList = ref<any[]>([])
const sceneSearch = ref('')
const scenePagination = reactive({ current: 1, pageSize: 12, total: 0 })
const sceneDialog = reactive({ visible: false, mode: 'create', loading: false })
const sceneFormData = ref<any>({})
const sceneFileList = ref([])
// 场景生图状态
const selectedSceneIds = ref<number[]>([])
const generatingSceneIds = ref<number[]>([])
const batchGeneratingScene = ref(false)

// 道具
const propList = ref<any[]>([])
const propSearch = ref('')
const propPagination = reactive({ current: 1, pageSize: 12, total: 0 })
const propDialog = reactive({ visible: false, mode: 'create', loading: false })
const propFormData = ref<any>({})
const propFileList = ref([])

const getAuthToken = () => localStorage.getItem('token')
const uploadConfig = reactive({
    action: import.meta.env.VITE_API_URL + '/admin/v1/upload/singleUpload',
    headers: computed(() => ({ 'Authorization': `${getAuthToken()}` })),
})

// === 初始化 ===
const init = async () => {
    loading.value = true
    try {
        await Promise.all([
            loadProjectInfo(),
            loadEpisodes(),
            loadCharacters(),
            loadScenes(),
            loadProps()
        ])
    } catch (e) { console.error(e) } finally { loading.value = false }
}

const loadProjectInfo = async () => {
    const res = await findProjects(projectId)
    if (res.code === 0) project.value = res.data
}

// === 剧集逻辑 ===
const loadEpisodes = async () => {
    const res = await getScriptsList({ projectId, page: 1, pageSize: 100 })
    if (res.code === 0) {
        const list = res.data.list || []
        episodeList.value = list.sort((a: any, b: any) => Number(a.episodeNo) - Number(b.episodeNo))
    }
}
const handleAddEpisode = () => {
    const maxEpisodeNo = episodeList.value.reduce((max, item) => Math.max(max, Number(item.episodeNo || 0)), 0)
    const nextEp = maxEpisodeNo + 1
    router.push(`/admin/projects/chapter/${projectId}/${nextEp}`)
}
const enterEpisode = (ep: any) => {
    router.push(`/admin/projects/chapter/${projectId}/${ep.episodeNo}`)
}

// === 生图通用逻辑 (轮询) ===
const pollTask = async (taskId: string, onSuccess: () => void, onFail: () => void) => {
    const timer = setInterval(async () => {
        try {
            const res = await findTasks(taskId)
            const data = res.data?.data || res.data
            const status = data?.status
            if (status === 'completed' || status === 2) {
                clearInterval(timer); onSuccess()
            } else if (status === 'failed' || status === 3) {
                clearInterval(timer); MessagePlugin.error('任务执行失败'); onFail()
            }
        } catch { clearInterval(timer); onFail() }
    }, 2000)
}

// 单个生图入口：明确使用 generateCharacterImageTask / generateSceneImageTask
const singleGenerate = async (type: 'char' | 'scene', item: any) => {
    const generatingList = type === 'char' ? generatingCharIds : generatingSceneIds
    const loadFunc = type === 'char' ? loadCharacters : loadScenes

    // 防止重复点击
    if (generatingList.value.includes(item.id)) return
    generatingList.value.push(item.id)

    try {
        let res
        if (type === 'char') {
            // 使用单个角色生图接口
            // 注意：请确保 api/tasks.ts 中定义了 generateCharacterImageTask 且参数名匹配后端 (通常是 characterId)
            res = await generateCharacterImageTask({ characterId: item.id })
        } else {
            // 使用单个场景生图接口
            res = await generateSceneImageTask({ sceneId: item.id })
        }

        // 兼容不同的后端返回结构
        const taskId = res.data?.data?.task_id || res.data?.taskId || res.data?.task_id

        if (taskId) {
            MessagePlugin.success('任务已提交')
            pollTask(taskId,
                () => { // Success
                    const idx = generatingList.value.indexOf(item.id)
                    if (idx > -1) generatingList.value.splice(idx, 1)
                    loadFunc()
                },
                () => { // Fail
                    const idx = generatingList.value.indexOf(item.id)
                    if (idx > -1) generatingList.value.splice(idx, 1)
                }
            )
        } else {
            throw new Error('未获取到任务ID')
        }
    } catch (e) {
        MessagePlugin.error('任务提交失败')
        const idx = generatingList.value.indexOf(item.id)
        if (idx > -1) generatingList.value.splice(idx, 1)
    }
}

// 批量生图处理：仅处理批量逻辑
const processBatchGeneration = async (type: 'char' | 'scene', ids: number[]) => {
    const generatingList = type === 'char' ? generatingCharIds : generatingSceneIds
    const loadFunc = type === 'char' ? loadCharacters : loadScenes
    const apiFunc = type === 'char' ? batchGenerateCharacterImagesTask : batchGenerateSceneImagesTask

    // 标记状态
    ids.forEach(id => { if (!generatingList.value.includes(id)) generatingList.value.push(id) })

    try {
        const res = await apiFunc(type === 'char' ? { characterIds: ids } : { sceneIds: ids })
        const taskList = res.data?.data || []

        if (taskList.length > 0) {
            MessagePlugin.success(`已提交 ${taskList.length} 个生图任务`)
            taskList.forEach((item: any) => {
                const id = type === 'char' ? item.character_id : item.scene_id
                const taskId = item.task_id

                pollTask(taskId,
                    () => {
                        const idx = generatingList.value.indexOf(id)
                        if (idx > -1) generatingList.value.splice(idx, 1)
                        loadFunc()
                    },
                    () => {
                        const idx = generatingList.value.indexOf(id)
                        if (idx > -1) generatingList.value.splice(idx, 1)
                    }
                )
            })
        }
    } catch (e) {
        MessagePlugin.error('任务提交失败')
        ids.forEach(id => {
            const idx = generatingList.value.indexOf(id)
            if (idx > -1) generatingList.value.splice(idx, 1)
        })
    }
}

// 批量生图入口
const batchGenerate = async (type: 'char' | 'scene') => {
    const selectedIds = type === 'char' ? selectedCharIds.value : selectedSceneIds.value
    const generatingRef = type === 'char' ? batchGeneratingChar : batchGeneratingScene

    if (selectedIds.length === 0) return
    generatingRef.value = true

    // 过滤掉已经在生成中的
    const idsToGen = selectedIds.filter(id => {
        const list = type === 'char' ? generatingCharIds.value : generatingSceneIds.value
        return !list.includes(id)
    })

    try {
        if (idsToGen.length > 0) {
            await processBatchGeneration(type, idsToGen)
        }
    } finally {
        generatingRef.value = false
    }
}

// === 角色逻辑 ===
const loadCharacters = async () => {
    const res = await getCharactersList({
        projectId, page: charPagination.current, pageSize: charPagination.pageSize, name: charSearch.value || undefined
    })
    if (res.code === 0) {
        characterList.value = res.data.list || []
        charPagination.total = res.data.total || 0
    }
}
// 角色选择
const checkAllChars = computed(() => characterList.value.length > 0 && selectedCharIds.value.length === characterList.value.length)
const isCharIndeterminate = computed(() => selectedCharIds.value.length > 0 && selectedCharIds.value.length < characterList.value.length)
const handleSelectAllChars = (checked: boolean) => { selectedCharIds.value = checked ? characterList.value.map((c: any) => c.id) : [] }
const toggleCharSelection = (id: number) => { const idx = selectedCharIds.value.indexOf(id); idx > -1 ? selectedCharIds.value.splice(idx, 1) : selectedCharIds.value.push(id) }

// ... 角色增删改 (保持原样) ...
const openCharacterDialog = (mode: string, row?: any) => {
    charDialog.mode = mode; charDialog.visible = true; charFileList.value = []
    charFormData.value = mode === 'edit' && row ? { ...row } : { projectId: Number(projectId), name: '', roleType: 'main', appearanceDesc: '' }
}
const submitCharacter = async () => {
    charDialog.loading = true
    try {
        const isEdit = charDialog.mode === 'edit'
        const api = isEdit ? updateCharacters : createCharacters
        const payload = { ...charFormData.value }
        if (isEdit) await api(charFormData.value.id, payload); else await api(payload)
        MessagePlugin.success(isEdit ? '更新成功' : '创建成功'); charDialog.visible = false; loadCharacters()
    } catch { MessagePlugin.error('操作失败') } finally { charDialog.loading = false }
}
const handleDeleteCharacter = async (id: number) => { await deleteCharacters(id); MessagePlugin.success('删除成功'); loadCharacters() }

// === 场景逻辑 ===
const loadScenes = async () => {
    const res = await getScenesList({
        projectId, page: scenePagination.current, pageSize: scenePagination.pageSize, name: sceneSearch.value || undefined
    })
    if (res.code === 0) {
        sceneList.value = res.data.list || []
        scenePagination.total = res.data.total || 0
    }
}
// 场景选择
const checkAllScenes = computed(() => sceneList.value.length > 0 && selectedSceneIds.value.length === sceneList.value.length)
const isSceneIndeterminate = computed(() => selectedSceneIds.value.length > 0 && selectedSceneIds.value.length < sceneList.value.length)
const handleSelectAllScenes = (checked: boolean) => { selectedSceneIds.value = checked ? sceneList.value.map((s: any) => s.id) : [] }
const toggleSceneSelection = (id: number) => { const idx = selectedSceneIds.value.indexOf(id); idx > -1 ? selectedSceneIds.value.splice(idx, 1) : selectedSceneIds.value.push(id) }

// ... 场景增删改 ...
const openSceneDialog = (mode: string, row?: any) => {
    sceneDialog.mode = mode; sceneDialog.visible = true; sceneFileList.value = []
    sceneFormData.value = mode === 'edit' && row ? { ...row } : { projectId: Number(projectId), name: '', location: '', time: '', atmosphere: '' }
}
const submitScene = async () => {
    sceneDialog.loading = true
    try {
        const isEdit = sceneDialog.mode === 'edit'
        const api = isEdit ? updateScenes : createScenes
        const payload = { ...sceneFormData.value }
        if (isEdit) await api(sceneFormData.value.id, payload); else await api(payload)
        MessagePlugin.success(isEdit ? '更新成功' : '创建成功'); sceneDialog.visible = false; loadScenes()
    } catch { MessagePlugin.error('操作失败') } finally { sceneDialog.loading = false }
}
const handleDeleteScene = async (id: number) => { await deleteScenes(id); MessagePlugin.success('删除成功'); loadScenes() }

// === 道具逻辑 ===
const loadProps = async () => {
    try {
        const res = await getPropsList({ projectId, page: propPagination.current, pageSize: propPagination.pageSize, name: propSearch.value || undefined })
        if (res.code === 0) { propList.value = res.data.list || []; propPagination.total = res.data.total || 0 }
    } catch { MessagePlugin.error('加载道具失败') }
}
const openPropDialog = (mode: string, row?: any) => {
    propDialog.mode = mode; propDialog.visible = true; propFileList.value = []
    propFormData.value = mode === 'edit' && row ? { ...row } : { projectId: Number(projectId), name: '', type: '', description: '', imagePrompt: '', imageUrl: '' }
}
const submitProp = async () => {
    propDialog.loading = true
    try {
        const isEdit = propDialog.mode === 'edit'
        const api = isEdit ? updateProps : createProps
        const payload = { ...propFormData.value }
        if (isEdit) await api(propFormData.value.id, payload); else await api(payload)
        MessagePlugin.success(isEdit ? '更新成功' : '创建成功'); propDialog.visible = false; loadProps()
    } catch { MessagePlugin.error('操作失败') } finally { propDialog.loading = false }
}
const handleDeleteProp = async (id: number) => { await deleteProps(id); MessagePlugin.success('删除成功'); loadProps() }

// === 通用 ===
const handleUploadSuccess = (ctx: any, type: string) => {
    if (ctx.response?.code === 0 || ctx.response?.code === 200) {
        const url = ctx.response.data?.url || ctx.response.data?.file_url
        const fullUrl = url.startsWith('http') ? url : import.meta.env.VITE_API_URL.replace(/\/admin\/v1$/, '') + url
        if (type === 'char') charFormData.value.avatarUrl = fullUrl
        else if (type === 'scene') sceneFormData.value.visualPrompt = fullUrl
        else if (type === 'prop') propFormData.value.imageUrl = fullUrl
    }
}

const goBack = () => router.push('/admin/projects/list')
const getStatusText = (s: number) => ['草稿', '连载中', '已完结'][s] || '未知'
const getStatusTheme = (s: number) => ['default', 'primary', 'success'][s] || 'default'
const formatDate = (val: string) => dayjs(val).format('YYYY-MM-DD HH:mm')
const formatDuration = (sec: number) => { if (!sec) return '0s'; const m = Math.floor(sec / 60); return m > 0 ? `${m}分${sec % 60}秒` : `${sec}秒` }
const getAspectRatio = (settingsStr: any) => { try { const s = typeof settingsStr === 'string' ? JSON.parse(settingsStr) : settingsStr; return s?.ratio === '16:9' ? '16/9' : '9/16' } catch { return '9/16' } }
const getRatioLabel = (settingsStr: any) => { try { const s = typeof settingsStr === 'string' ? JSON.parse(settingsStr) : settingsStr; return s?.ratio === '16:9' ? '16:9 (横屏)' : '9:16 (竖屏)' } catch { return '未知' } }

onMounted(init)
</script>

<style scoped lang="less">
.project-detail-container {
    min-height: 100vh;
    background: var(--td-bg-color-container);
    display: flex;
    flex-direction: column;
}

.detail-header {
    background: #fff;
    padding: 16px 24px;
    border-bottom: 1px solid var(--td-component-stroke);
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-left {
        display: flex;
        align-items: center;
        gap: 16px;

        .project-info {
            .title-row {
                display: flex;
                align-items: center;
                gap: 8px;

                .title {
                    font-size: 18px;
                    font-weight: 700;
                    color: var(--td-text-color-primary);
                }
            }

            .desc {
                font-size: 13px;
                color: var(--td-text-color-secondary);
                margin-top: 4px;
            }
        }
    }
}

.detail-content {
    flex: 1;
    padding: 24px;

    .detail-tabs {
        background: #fff;
        border-radius: 8px;

        .tab-panel-content {
            padding: 20px;
            min-height: 400px;
        }

        .overview-panel {
            .cover-section {
                border-radius: 8px;
                overflow: hidden;
                box-shadow: var(--td-shadow-1);

                .project-cover-large {
                    width: 100%;
                    display: block;
                }
            }

            .info-card {
                background: var(--td-bg-color-container);

                .stat-item {
                    text-align: center;

                    .label {
                        color: var(--td-text-color-secondary);
                        font-size: 12px;
                    }

                    .num {
                        font-size: 24px;
                        font-weight: 700;
                        color: var(--td-brand-color);
                        margin-top: 4px;
                    }
                }
            }
        }
    }
}

.action-bar {
    display: flex;
    justify-content: space-between;
    margin-bottom: 20px;

    .left-actions {
        display: flex;
        gap: 12px;
    }
}

.selection-bar {
    margin-bottom: 16px;
    padding: 8px 12px;
    background: var(--td-bg-color-secondarycontainer);
    border-radius: 6px;
}

.episode-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;

    .episode-card {
        cursor: pointer;
        transition: all 0.2s;

        .ep-cover {
            height: 100px;
            background: var(--td-brand-color-light);
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            color: var(--td-brand-color);
            border-radius: 4px;
            margin-bottom: 12px;
            position: relative;

            .ep-no {
                font-weight: 700;
                font-size: 18px;
            }

            .ep-status {
                font-size: 12px;
                margin-top: 4px;
                display: flex;
                align-items: center;
                gap: 4px;
                color: var(--td-success-color);
            }
        }

        .ep-title {
            font-weight: 600;
            margin-bottom: 4px;
        }

        .ep-meta {
            font-size: 12px;
            color: var(--td-text-color-secondary);
        }

        &:hover {
            transform: translateY(-2px);
            border-color: var(--td-brand-color);
        }
    }
}

.resource-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 20px;

    .resource-card {
        background: var(--td-bg-color-secondarycontainer);
        border: 1px solid transparent;
        transition: all 0.2s;
        position: relative;

        &:hover {
            border-color: var(--td-brand-color);
            transform: translateY(-2px);

            .res-overlay {
                opacity: 1;
            }
        }

        &.is-selected {
            border-color: var(--td-brand-color);
            background: var(--td-brand-color-light);
        }

        .card-select {
            position: absolute;
            top: 8px;
            left: 8px;
            z-index: 2;
        }

        .res-cover {
            height: 180px;
            position: relative;
            /* 必须是 relative */
            border-radius: 4px;
            overflow: hidden;
            background: #fff;

            &.scene-cover {
                height: 120px;
            }

            .res-img {
                width: 100%;
                height: 100%;
                display: block;
                /* 防止图片底部留白 */
            }

            /* 修复 Loading 遮罩层级 */
            .loading-mask {
                position: absolute;
                inset: 0;
                background: rgba(255, 255, 255, 0.8);
                display: flex;
                align-items: center;
                justify-content: center;
                z-index: 10;
                /* 确保 loading 在最上层 */
            }

            /* 修复操作遮罩层级 */
            .res-overlay {
                position: absolute;
                inset: 0;
                /* 铺满父容器 */
                background: rgba(0, 0, 0, 0.6);
                /* 稍微加深一点颜色 */
                display: flex;
                align-items: center;
                justify-content: center;
                gap: 12px;
                /* 增加按钮间距 */

                /* 默认隐藏，hover显示 */
                opacity: 0;
                z-index: 5;
                /* 必须比图片层级高 */
                transition: opacity 0.3s ease;
            }
        }

        /* 鼠标悬停在卡片上时显示遮罩 */
        &:hover {
            .res-overlay {
                opacity: 1;
            }
        }

        .res-info {
            margin-top: 12px;
            display: flex;
            justify-content: space-between;
            align-items: center;

            .res-name {
                font-weight: 600;
                font-size: 14px;
            }

            .scene-meta {
                font-size: 12px;
                color: var(--td-text-color-secondary);
            }
        }

        .res-desc {
            font-size: 12px;
            color: var(--td-text-color-placeholder);
            margin-top: 4px;
            height: 36px;
        }
    }
}

.upload-box {
    width: 100px;
    height: 100px;
    border: 1px dashed var(--td-component-stroke);
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;

    &:hover {
        border-color: var(--td-brand-color);
        color: var(--td-brand-color);
    }

    &.scene {
        width: 160px;
        height: 100px;
    }
}

.upload-preview {
    width: 100px;
    height: 100px;
    border-radius: 4px;

    &.scene {
        width: 160px;
        height: 100px;
    }
}

.text-ellipsis-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
</style>