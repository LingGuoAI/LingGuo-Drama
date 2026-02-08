<template>
    <div class="professional-editor">
        <div class="editor-header">
            <div class="header-left">
                <t-button variant="text" shape="circle" @click="goBack">
                    <template #icon><t-icon name="arrow-left" /></template>
                </t-button>
                <div class="header-title">
                    <span class="title">{{ drama?.title || '加载中...' }}</span>
                    <span class="sub-title"> - 第 {{ episodeNumber }} 集</span>
                </div>
            </div>
            <div class="header-right">
                <t-button theme="default" variant="outline" size="small" @click="loadData">
                    <template #icon><t-icon name="refresh" /></template> 刷新
                </t-button>
            </div>
        </div>

        <div class="editor-main" v-loading="loading">

            <div class="storyboard-panel">
                <div class="panel-header">
                    <h3>分镜列表</h3>
                    <t-button theme="primary" variant="text" size="small" @click="handleAddStoryboard">
                        <template #icon><t-icon name="add" /></template>新增
                    </t-button>
                </div>

                <div class="storyboard-list">
                    <div v-for="shot in storyboards" :key="shot.id" class="storyboard-item"
                        :class="{ active: currentStoryboardId === shot.id }" @click="selectStoryboard(shot.id)">
                        <div class="shot-content">
                            <div class="shot-header">
                                <div class="shot-title-row">
                                    <span class="shot-number">#{{ shot.sequenceNo || shot.storyboard_number }}</span>
                                    <span class="shot-title" :title="shot.title">{{ shot.title || '未命名' }}</span>
                                </div>
                                <div class="shot-actions">
                                    <span class="shot-duration">{{ shot.duration || shot.durationMs / 1000 }}s</span>
                                    <t-popconfirm content="确认删除此镜头吗？" @confirm="handleDeleteStoryboard(shot)">
                                        <t-button shape="circle" size="small" theme="danger" variant="text" @click.stop>
                                            <template #icon><t-icon name="delete" /></template>
                                        </t-button>
                                    </t-popconfirm>
                                </div>
                            </div>
                            <div class="shot-desc">{{ shot.visualDesc || shot.action || '暂无描述' }}</div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="timeline-area">
                <div v-if="storyboards.length > 0" class="timeline-wrapper">
                    <div class="placeholder-timeline">
                        <t-icon name="film" size="48px" style="color: var(--td-brand-color); margin-bottom: 16px;" />
                        <div>时间线编辑器区域 (VideoTimelineEditor)</div>
                        <div style="font-size: 12px; color: #999; margin-top: 8px;">当前选中镜头 ID: {{ currentStoryboardId }}
                        </div>
                    </div>
                </div>
                <t-empty v-else description="暂无分镜数据" class="empty-timeline" />
            </div>

            <div class="edit-panel">
                <t-tabs v-model="activeTab" theme="normal" class="edit-tabs">

                    <t-tab-panel value="shot" label="镜头属性">
                        <div class="tab-content" v-if="currentStoryboard">
                            <t-form label-align="top">
                                <div class="info-block">
                                    <div class="block-title">场景信息</div>
                                    <div class="scene-preview-card" v-if="currentStoryboard.sceness">
                                        <t-image
                                            v-if="currentStoryboard.sceness.visualPrompt && currentStoryboard.sceness.visualPrompt.startsWith('http')"
                                            :src="getImageUrl(currentStoryboard.sceness.visualPrompt)" fit="cover"
                                            class="scene-img" />
                                        <div class="scene-text">
                                            <div class="loc">{{ currentStoryboard.sceness.location }}</div>
                                            <div class="time">{{ currentStoryboard.sceness.time }}</div>
                                        </div>
                                    </div>
                                    <t-button v-else block variant="dashed"
                                        @click="showSceneSelector = true">关联场景</t-button>
                                </div>

                                <t-row :gutter="16">
                                    <t-col :span="6">
                                        <t-form-item label="景别">
                                            <t-select v-model="currentStoryboard.shotType"
                                                :options="['大远景', '远景', '全景', '中景', '近景', '特写'].map(v => ({ label: v, value: v }))"
                                                @change="saveStoryboardField('shotType')" />
                                        </t-form-item>
                                    </t-col>
                                    <t-col :span="6">
                                        <t-form-item label="视角">
                                            <t-select v-model="currentStoryboard.angle"
                                                :options="['平视', '俯视', '仰视', '侧视'].map(v => ({ label: v, value: v }))"
                                                @change="saveStoryboardField('angle')" />
                                        </t-form-item>
                                    </t-col>
                                </t-row>
                                <t-form-item label="运镜">
                                    <t-select v-model="currentStoryboard.cameraMovement"
                                        :options="['固定', '推', '拉', '摇', '移', '跟'].map(v => ({ label: v, value: v }))"
                                        @change="saveStoryboardField('cameraMovement')" />
                                </t-form-item>

                                <t-form-item label="台词 (Dialogue)">
                                    <t-textarea v-model="currentStoryboard.dialogue" :autosize="{ minRows: 2 }"
                                        @blur="saveStoryboardField('dialogue')" />
                                </t-form-item>
                                <t-form-item label="画面描述 (Visual)">
                                    <t-textarea v-model="currentStoryboard.visualDesc" :autosize="{ minRows: 3 }"
                                        @blur="saveStoryboardField('visualDesc')" />
                                </t-form-item>
                                <t-form-item label="环境氛围 (Atmosphere)">
                                    <t-textarea v-model="currentStoryboard.atmosphere" :autosize="{ minRows: 2 }"
                                        @blur="saveStoryboardField('atmosphere')" />
                                </t-form-item>
                            </t-form>
                        </div>
                        <t-empty v-else description="请选择一个镜头" style="margin-top: 40px" />
                    </t-tab-panel>

                    <t-tab-panel value="image" label="画面生成">
                        <div class="tab-content" v-if="currentStoryboard">
                            <t-form-item label="绘画提示词 (Prompt)" style="margin-bottom: 12px;">
                                <t-textarea v-model="currentStoryboard.imagePrompt" :rows="5"
                                    placeholder="输入用于AI绘画的提示词..." />
                            </t-form-item>

                            <div class="action-bar">
                                <t-button theme="primary" :loading="generatingImage" @click="generateImage">
                                    <template #icon><t-icon name="magic" /></template> 生成画面
                                </t-button>
                                <t-upload theme="custom" :action="uploadConfig.action" :headers="uploadConfig.headers"
                                    :show-file-list="false" accept="image/*" @success="handleUploadImageSuccess">
                                    <t-button variant="outline"
                                        :icon="() => import('tdesign-icons-vue-next').UploadIcon">上传</t-button>
                                </t-upload>
                            </div>

                            <div class="image-result-area" v-if="currentStoryboard.imageUrl">
                                <div class="section-title">当前画面</div>
                                <div class="image-wrapper">
                                    <t-image :src="getImageUrl(currentStoryboard.imageUrl)" fit="contain" />
                                    <div class="img-actions">
                                        <t-button shape="circle" size="small" variant="text"
                                            @click="previewImage(currentStoryboard.imageUrl)"><t-icon
                                                name="zoom-in" /></t-button>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <t-empty v-else description="请选择一个镜头" style="margin-top: 40px" />
                    </t-tab-panel>

                    <t-tab-panel value="video" label="视频生成">
                        <div class="tab-content" v-if="currentStoryboard">
                            <div class="video-settings">
                                <t-form-item label="视频模型">
                                    <t-select v-model="selectedVideoModel" :options="videoModelOptions" />
                                </t-form-item>
                                <t-form-item label="时长 (秒)">
                                    <t-slider v-model="videoDuration" :min="2" :max="10" />
                                </t-form-item>

                                <t-button theme="primary" block size="large" :loading="generatingVideo"
                                    @click="generateVideo" :disabled="!currentStoryboard.imageUrl">
                                    <template #icon><t-icon name="video" /></template>
                                    {{ currentStoryboard.imageUrl ? '图生视频' : '请先生成图片' }}
                                </t-button>
                            </div>

                            <div class="video-list-area" v-if="currentStoryboard.videoUrl">
                                <div class="section-title">生成结果</div>
                                <div class="video-card">
                                    <video :src="getVideoUrl(currentStoryboard.videoUrl)" controls
                                        style="width: 100%; border-radius: 8px;"></video>
                                    <div class="video-actions">
                                        <t-tag theme="success" variant="light">已完成</t-tag>
                                        <t-button size="small" variant="text" @click="deleteVideo">删除</t-button>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <t-empty v-else description="请选择一个镜头" style="margin-top: 40px" />
                    </t-tab-panel>
                </t-tabs>
            </div>
        </div>

        <t-dialog v-model:visible="showSceneSelector" header="选择场景" width="600px">
            <div class="scene-grid">
                <t-empty description="开发中..." />
            </div>
        </t-dialog>

    </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import {
    ArrowLeftIcon, RefreshIcon, AddIcon, DeleteIcon, MagicIcon,
    UploadIcon, ZoomInIcon, VideoIcon, FilmIcon
} from 'tdesign-icons-vue-next'

// API (复用之前的 API 定义)
import { findProjects } from '@/api/projects'
import { getScriptsList } from '@/api/scripts'
import { getScenesList } from '@/api/scenes'
import { generateSceneImageTask, findTasks } from '@/api/tasks' // 假设有单镜头生图/视频接口
// 注意：这里需要引入 updateShots 等接口，如果没有请在 api/shots.ts 中补充
import { updateShots, findShots, deleteShots } from '@/api/shots' // 假设您有这些接口
import { getImageUrl } from '@/utils/format'

const route = useRoute()
const router = useRouter()

// === 核心数据 ===
const loading = ref(false)
const dramaId = route.params.id as string
const episodeNumber = Number(route.params.episodeNumber)

const project = ref<any>({})
const storyboards = ref<any[]>([]) // 镜头列表
const currentStoryboardId = ref<string | number | null>(null)
const availableScenes = ref<any[]>([])

// === UI 状态 ===
const activeTab = ref('shot')
const showSceneSelector = ref(false)

// === 视频/图片生成状态 ===
const generatingImage = ref(false)
const generatingVideo = ref(false)
const selectedVideoModel = ref('kling')
const videoDuration = ref(5)

const videoModelOptions = [
    { label: '可灵 (Kling)', value: 'kling' },
    { label: 'Runway Gen-3', value: 'runway' },
    { label: 'Luma Dream Machine', value: 'luma' },
    { label: 'Sora (OpenAI)', value: 'sora' }
]

// === Computed ===
const currentStoryboard = computed(() => {
    return storyboards.value.find(s => s.id === currentStoryboardId.value)
})

const getAuthToken = () => localStorage.getItem('token')
const uploadConfig = reactive({
    action: import.meta.env.VITE_API_URL + '/admin/v1/upload/singleUpload',
    headers: computed(() => ({ 'Authorization': `${getAuthToken()}` })),
})

// === 初始化 ===
const initData = async () => {
    loading.value = true
    try {
        // 1. 获取项目信息
        const res = await findProjects(dramaId)
        if (res.code === 0) project.value = res.data

        // 2. 获取分镜列表 (这里逻辑与 createChapter 类似，先获取脚本，再获取 shots)
        // 为了简化，假设后端有个直接按集数获取所有 shots 的接口，或者复用之前的逻辑
        await loadShots()

    } catch (e) { console.error(e) } finally { loading.value = false }
}

const loadShots = async () => {
    // 这里需要根据您的后端接口适配
    // 假设：先获取 scriptId，再获取 shots
    // ... 代码略，参考 createChapter.vue 的 loadScriptDetail ...
    // 模拟数据：
    // storyboards.value = [...] 
    // 真实逻辑：
    const listRes = await getScriptsList({ projectId: dramaId, page: 1, pageSize: 100 })
    const targetScript = listRes.data?.list?.find((s: any) => Number(s.episodeNo) === episodeNumber)
    if (targetScript) {
        // 假设有个接口获取分镜详情
        // const detailRes = await findScripts(targetScript.id) 
        // storyboards.value = detailRes.data.shots

        // 或者直接查 shots 表
        // const shotsRes = await getShotsList({ scriptId: targetScript.id })
    }
}

// === 交互逻辑 ===

const goBack = () => {
    router.back()
}

const loadData = () => {
    initData()
    MessagePlugin.success('数据已刷新')
}

const selectStoryboard = (id: number | string) => {
    currentStoryboardId.value = id
}

const handleAddStoryboard = () => {
    MessagePlugin.info('新增分镜功能待实现')
}

const handleDeleteStoryboard = async (shot: any) => {
    try {
        await deleteShots(shot.id)
        MessagePlugin.success('删除成功')
        loadShots()
        if (currentStoryboardId.value === shot.id) currentStoryboardId.value = null
    } catch {
        MessagePlugin.error('删除失败')
    }
}

// 字段保存 (防抖或失焦保存)
const saveStoryboardField = async (field: string) => {
    if (!currentStoryboard.value) return
    const id = currentStoryboard.value.id
    const payload = { [field]: currentStoryboard.value[field] }
    try {
        await updateShots(id, payload)
        // MessagePlugin.success('保存成功') // 可选，避免太打扰
    } catch {
        MessagePlugin.error('自动保存失败')
    }
}

// 生成图片
const generateImage = async () => {
    if (!currentStoryboard.value) return
    generatingImage.value = true
    // 调用生图接口...
    setTimeout(() => {
        generatingImage.value = false
        MessagePlugin.info('生图请求已发送，请稍后刷新')
    }, 1000)
}

// 上传成功
const handleUploadImageSuccess = (context: any) => {
    if (context.response?.code === 0) {
        if (currentStoryboard.value) {
            currentStoryboard.value.imageUrl = context.response.data.url // 需加上域名处理
            saveStoryboardField('imageUrl')
            MessagePlugin.success('上传成功')
        }
    }
}

// 预览图片
const previewImage = (url: string) => {
    // TDesign 的 ImageViewer 使用方式需要根据文档调整，这里简单处理
    window.open(url, '_blank')
}

// 生成视频
const generateVideo = async () => {
    if (!currentStoryboard.value?.imageUrl) return
    generatingVideo.value = true
    // 调用图生视频接口...
    setTimeout(() => {
        generatingVideo.value = false
        MessagePlugin.info('视频生成请求已发送')
    }, 1000)
}

const deleteVideo = () => {
    if (currentStoryboard.value) {
        currentStoryboard.value.videoUrl = ''
        saveStoryboardField('videoUrl')
    }
}

// Helper
const getVideoUrl = (url: string) => {
    if (!url) return ''
    return url.startsWith('http') ? url : import.meta.env.VITE_API_URL + url
}

onMounted(() => {
    initData()
})
</script>

<style scoped lang="less">
.professional-editor {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--td-bg-color-container);

    .editor-header {
        height: 56px;
        background: #fff;
        border-bottom: 1px solid var(--td-component-stroke);
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0 16px;
        flex-shrink: 0;

        .header-left {
            display: flex;
            align-items: center;
            gap: 12px;

            .header-title {
                .title {
                    font-weight: 700;
                    font-size: 16px;
                }

                .sub-title {
                    color: var(--td-text-color-secondary);
                    font-size: 14px;
                }
            }
        }
    }

    .editor-main {
        flex: 1;
        display: flex;
        overflow: hidden;

        // 1. 左侧分镜列表
        .storyboard-panel {
            width: 280px;
            background: #fff;
            border-right: 1px solid var(--td-component-stroke);
            display: flex;
            flex-direction: column;
            flex-shrink: 0;

            .panel-header {
                padding: 12px 16px;
                border-bottom: 1px solid var(--td-component-stroke);
                display: flex;
                justify-content: space-between;
                align-items: center;

                h3 {
                    margin: 0;
                    font-size: 14px;
                    font-weight: 600;
                }
            }

            .storyboard-list {
                flex: 1;
                overflow-y: auto;
                padding: 12px;
                display: flex;
                flex-direction: column;
                gap: 8px;

                .storyboard-item {
                    border: 1px solid var(--td-component-stroke);
                    border-radius: 6px;
                    padding: 10px;
                    cursor: pointer;
                    transition: all 0.2s;
                    background: var(--td-bg-color-container);

                    &:hover {
                        border-color: var(--td-brand-color);
                        box-shadow: var(--td-shadow-1);
                    }

                    &.active {
                        border-color: var(--td-brand-color);
                        background: var(--td-brand-color-light);
                    }

                    .shot-header {
                        display: flex;
                        justify-content: space-between;
                        margin-bottom: 4px;
                        font-size: 12px;

                        .shot-title-row {
                            font-weight: 600;
                            color: var(--td-text-color-primary);
                        }

                        .shot-actions {
                            display: flex;
                            align-items: center;
                            gap: 4px;
                            color: var(--td-text-color-secondary);
                        }
                    }

                    .shot-desc {
                        font-size: 12px;
                        color: var(--td-text-color-secondary);
                        display: -webkit-box;
                        -webkit-line-clamp: 2;
                        -webkit-box-orient: vertical;
                        overflow: hidden;
                    }
                }
            }
        }

        // 2. 中间时间线
        .timeline-area {
            flex: 1;
            background: var(--td-bg-color-secondarycontainer);
            position: relative;
            display: flex;
            flex-direction: column;

            .timeline-wrapper {
                flex: 1;
                display: flex;
                align-items: center;
                justify-content: center;
            }

            .placeholder-timeline {
                text-align: center;
                color: var(--td-text-color-disabled);
            }
        }

        // 3. 右侧编辑面板
        .edit-panel {
            width: 360px;
            background: #fff;
            border-left: 1px solid var(--td-component-stroke);
            flex-shrink: 0;
            display: flex;
            flex-direction: column;

            .edit-tabs {
                flex: 1;
                display: flex;
                flex-direction: column;

                :deep(.t-tabs__content) {
                    flex: 1;
                    overflow-y: auto;
                    padding: 0; // 移除默认padding
                }
            }

            .tab-content {
                padding: 16px;

                .info-block {
                    margin-bottom: 24px;

                    .block-title {
                        font-weight: 600;
                        margin-bottom: 8px;
                        font-size: 13px;
                    }

                    .scene-preview-card {
                        border: 1px solid var(--td-component-stroke);
                        border-radius: 6px;
                        overflow: hidden;

                        .scene-img {
                            height: 100px;
                            width: 100%;
                        }

                        .scene-text {
                            padding: 8px;
                            font-size: 12px;

                            .loc {
                                font-weight: 600;
                            }

                            .time {
                                color: var(--td-text-color-secondary);
                            }
                        }
                    }
                }

                .action-bar {
                    display: flex;
                    gap: 10px;
                    margin-bottom: 16px;

                    :deep(.t-upload) {
                        display: inline-block;
                    }
                }

                .image-result-area {
                    margin-top: 16px;

                    .section-title {
                        font-weight: 600;
                        margin-bottom: 8px;
                    }

                    .image-wrapper {
                        position: relative;
                        border-radius: 8px;
                        overflow: hidden;
                        border: 1px solid var(--td-component-stroke);

                        .img-actions {
                            position: absolute;
                            top: 8px;
                            right: 8px;
                            opacity: 0;
                            transition: opacity 0.2s;
                        }

                        &:hover .img-actions {
                            opacity: 1;
                        }
                    }
                }

                .video-settings {
                    display: flex;
                    flex-direction: column;
                    gap: 16px;
                }
            }
        }
    }
}
</style>