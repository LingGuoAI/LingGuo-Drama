<template>
    <div class="professional-editor">
        <div class="editor-header">
            <div class="header-left">
                <t-button variant="text" shape="circle" @click="goBack">
                    <template #icon><t-icon name="arrow-left" /></template>
                </t-button>
                <div class="header-title">
                    <span class="title">{{ project?.title || '加载中...' }}</span>
                    <t-tag theme="primary" variant="light" style="margin-left: 8px;">第 {{ episodeNumber }} 集</t-tag>
                </div>
            </div>
            <div class="header-right">
                <div class="status-text" v-if="saving"><t-loading size="small" /> 自动保存中...</div>
                <t-button theme="default" variant="outline" size="small" @click="loadData">
                    <template #icon><t-icon name="refresh" /></template> 刷新
                </t-button>
                <t-button theme="primary" size="small" @click="exportVideo">
                    <template #icon><t-icon name="download" /></template> 导出视频
                </t-button>
            </div>
        </div>

        <div class="editor-main" v-loading="loading">
            <div class="left-sidebar">
                <div class="storyboard-panel">
                    <div class="panel-header">
                        <h3>分镜列表 ({{ storyboards.length }})</h3>
                        <t-button theme="primary" variant="text" size="small" @click="handleAddStoryboard">
                            <template #icon><t-icon name="add" /></template>新增
                        </t-button>
                    </div>
                    <div class="storyboard-list" v-if="storyboards.length > 0">
                        <div v-for="(shot, index) in storyboards" :key="shot.id" class="storyboard-item"
                            :class="{ active: currentStoryboardId === shot.id }" @click="selectStoryboard(shot.id)"
                            draggable="true" @dragstart="handleDragStart($event, shot, 'storyboard')">
                            <div class="shot-thumb">
                                <t-image v-if="shot.image || shot.imageUrl"
                                    :src="getImageUrl(shot.image || shot.imageUrl)" fit="cover"
                                    style="width: 100%; height: 100%;" />
                                <div v-else class="drag-hint"><t-icon name="move" /></div>
                            </div>
                            <div class="shot-content">
                                <div class="shot-header">
                                    <div class="shot-title-row">
                                        <span class="shot-number">#{{ index + 1 }}</span>
                                        <span class="shot-title" :title="shot.title">{{ shot.title || '未命名镜头' }}</span>
                                    </div>
                                    <div class="shot-actions">
                                        <span class="shot-duration">{{ shot.duration || 0 }}s</span>
                                        <t-popconfirm content="确认删除此镜头吗？" @confirm="handleDeleteStoryboard(shot)">
                                            <t-button shape="circle" size="small" theme="danger" variant="text"
                                                @click.stop>
                                                <template #icon><t-icon name="delete" /></template>
                                            </t-button>
                                        </t-popconfirm>
                                    </div>
                                </div>
                                <div class="shot-desc">{{ shot.visualDesc || shot.action || '暂无画面描述' }}</div>
                            </div>
                        </div>
                    </div>
                    <t-empty v-else description="暂无分镜，点击新增" style="padding: 20px 0" />
                </div>

                <div class="assets-panel">
                    <div class="panel-header">
                        <h3>素材库 ({{ videoAssets.length }})</h3>
                        <t-button theme="default" variant="text" size="small" @click="loadVideoAssets">
                            <template #icon><t-icon name="refresh" /></template>
                        </t-button>
                    </div>
                    <div class="assets-grid" v-if="videoAssets.length > 0">
                        <div v-for="asset in videoAssets" :key="asset.id" class="asset-item" draggable="true"
                            @dragstart="handleDragStart($event, asset, 'asset')">
                            <div class="asset-thumb">
                                <video :src="getVideoUrl(asset.url)" muted @mouseenter="$event.target.play()"
                                    @mouseleave="$event.target.pause()"
                                    @loadedmetadata="$event.target.currentTime = 0"></video>
                                <span class="duration">{{ asset.duration }}s</span>
                                <div class="hover-overlay"><t-icon name="add-circle" /></div>
                            </div>
                            <div class="asset-name" :title="asset.name">{{ asset.name || '视频片段' }}</div>
                        </div>
                    </div>
                    <t-empty v-else description="暂无生成视频" size="small" class="empty-assets" />
                </div>
            </div>

            <div class="center-workspace">
                <div class="preview-stage">
                    <div class="player-container">
                        <video v-if="currentPreviewUrl" ref="mainPlayerRef" :src="currentPreviewUrl" controls
                            class="main-player"></video>
                        <div v-else class="player-placeholder">
                            <t-icon name="film" size="48px" />
                            <p>请在时间线上选择片段或点击播放</p>
                        </div>
                    </div>
                </div>
                <div class="timeline-stage">
                    <VideoTimelineEditor ref="timelineEditorRef" :clips="timelineClips" :audio-clips="audioClips"
                        :current-time="currentTime" :total-duration="totalDuration" :current-id="currentStoryboardId"
                        @update:time="updateCurrentTime" @drop-clip="handleTimelineDrop"
                        @select-clip="handleTimelineSelect" @delete-clip="removeClipFromTimeline" />
                </div>
            </div>

            <div class="edit-panel">
                <t-tabs v-model="activeTab" theme="normal" class="edit-tabs">

                    <t-tab-panel value="shot" label="镜头属性">
                        <div class="tab-content scrollable-content" v-if="currentStoryboard">
                            <t-form label-align="top" class="compact-form">
                                <div class="section-group">
                                    <div class="section-header">
                                        <span>场景 (Scene)</span>
                                        <t-button theme="primary" variant="text" size="small"
                                            @click="showSceneSelector = true">更换场景</t-button>
                                    </div>
                                    <div class="scene-card" v-if="currentScene">
                                        <t-image v-if="currentScene.visualPrompt"
                                            :src="getImageUrl(currentScene.visualPrompt)" fit="cover"
                                            class="scene-cover" @click="previewImage(currentScene.visualPrompt)" />
                                        <div class="scene-info">
                                            <div class="scene-loc">{{ currentScene.name }}</div>
                                            <div class="scene-meta">{{ currentScene.location }} · {{ currentScene.time
                                            }}</div>
                                        </div>
                                    </div>
                                    <div v-else class="empty-box" @click="showSceneSelector = true">
                                        <t-icon name="image" /> <span>点击关联场景</span>
                                    </div>
                                </div>

                                <div class="section-group">
                                    <div class="section-header">
                                        <span>登场角色 (Cast)</span>
                                        <t-button theme="primary" variant="text" size="small"
                                            @click="showCharacterSelector = true">
                                            <template #icon><t-icon name="add" /></template> 添加
                                        </t-button>
                                    </div>
                                    <div class="cast-list" v-if="selectedCharacters.length > 0">
                                        <div v-for="charId in selectedCharacters" :key="charId" class="cast-item">
                                            <t-avatar :image="getCharacterById(charId)?.avatarUrl" size="medium"
                                                shape="circle" />
                                            <span class="cast-name" :title="getCharacterById(charId)?.name">{{
                                                getCharacterById(charId)?.name }}</span>
                                            <div class="remove-btn" @click="toggleCharacterInShot(charId)"><t-icon
                                                    name="close" /></div>
                                        </div>
                                    </div>
                                    <div v-else class="empty-text">暂无角色</div>
                                </div>

                                <div class="section-group">
                                    <div class="section-header">
                                        <span>相关道具 (Props)</span>
                                        <t-button theme="primary" variant="text" size="small"
                                            @click="showPropSelector = true">
                                            <template #icon><t-icon name="add" /></template> 添加
                                        </t-button>
                                    </div>
                                    <div class="cast-list" v-if="selectedProps.length > 0">
                                        <div v-for="propId in selectedProps" :key="propId" class="cast-item">
                                            <t-image :src="getImageUrl(getPropById(propId)?.imageUrl)" fit="contain"
                                                style="width: 40px; height: 40px; border-radius: 4px; background: #eee;" />
                                            <span class="cast-name" :title="getPropById(propId)?.name">{{
                                                getPropById(propId)?.name }}</span>
                                            <div class="remove-btn" @click="togglePropInShot(propId)"><t-icon
                                                    name="close" />
                                            </div>
                                        </div>
                                    </div>
                                    <div v-else class="empty-text">暂无道具</div>
                                </div>

                                <t-divider />

                                <div class="section-group">
                                    <div class="section-header"><span>视效设置</span></div>
                                    <t-row :gutter="10">
                                        <t-col :span="6">
                                            <t-form-item label="景别">
                                                <t-select v-model="currentStoryboard.shotType" size="small"
                                                    :options="['大远景', '远景', '全景', '中景', '近景', '特写', '大特写'].map(v => ({ label: v, value: v }))"
                                                    @change="saveStoryboardField('shotType')" />
                                            </t-form-item>
                                        </t-col>
                                        <t-col :span="6">
                                            <t-form-item label="视角">
                                                <t-select v-model="currentStoryboard.cameraAngle" size="small"
                                                    :options="['平视', '俯视', '仰视', '侧视', '航拍'].map(v => ({ label: v, value: v }))"
                                                    @change="saveStoryboardField('cameraAngle')" />
                                            </t-form-item>
                                        </t-col>
                                    </t-row>
                                    <t-form-item label="运镜" style="margin-top: 10px;">
                                        <t-select v-model="currentStoryboard.cameraMove" size="small"
                                            :options="['固定', '推', '拉', '摇', '移', '跟', '环绕'].map(v => ({ label: v, value: v }))"
                                            @change="saveStoryboardField('cameraMove')" />
                                    </t-form-item>
                                    <t-form-item label="时长 (秒)" style="margin-top: 10px;">
                                        <t-slider v-model="currentStoryboard.duration" :min="2" :max="10"
                                            @change="saveStoryboardField('duration')" />
                                    </t-form-item>
                                </div>

                                <div class="section-group">
                                    <div class="section-header"><span>叙事内容</span></div>
                                    <t-form-item label="动作 (Action)">
                                        <t-textarea v-model="currentStoryboard.action" :autosize="{ minRows: 2 }"
                                            placeholder="角色做的动作..." @blur="saveStoryboardField('action')" />
                                    </t-form-item>
                                    <t-form-item label="对白 (Dialogue)">
                                        <t-textarea v-model="currentStoryboard.dialogue" :autosize="{ minRows: 2 }"
                                            placeholder="角色台词..." @blur="saveStoryboardField('dialogue')" />
                                    </t-form-item>
                                    <t-form-item label="画面描述 (Visual)">
                                        <t-textarea v-model="currentStoryboard.visualDesc" :autosize="{ minRows: 3 }"
                                            placeholder="详细的画面描述..." @blur="saveStoryboardField('visualDesc')" />
                                    </t-form-item>
                                </div>
                            </t-form>
                        </div>
                        <t-empty v-else description="请在左侧选择一个镜头" style="margin-top: 40px" />
                    </t-tab-panel>

                    <t-tab-panel value="image" label="镜头图片">
                        <div class="tab-content scrollable-content" v-if="currentStoryboard">

                            <div class="section-group">
                                <div class="section-header"><span>帧类型选择</span></div>
                                <t-radio-group variant="default-filled" v-model="selectedFrameType"
                                    style="width: 100%;">
                                    <t-radio-button value="first">首帧</t-radio-button>
                                    <t-radio-button value="last">尾帧</t-radio-button>
                                    <t-radio-button value="action">动作序列</t-radio-button>
                                    <t-radio-button value="key">关键帧</t-radio-button>
                                </t-radio-group>
                            </div>

                            <div class="section-group">
                                <div class="section-header"><span>AI 绘画提示词</span></div>
                                <t-textarea v-model="currentStoryboard.imagePrompt" :rows="4" placeholder="输入英文提示词..."
                                    @blur="saveStoryboardField('imagePrompt')" />
                            </div>

                            <div class="action-bar">
                                <t-button theme="primary" :loading="generatingImage" @click="generateImage">
                                    <template #icon><t-icon name="magic" /></template> 生成画面
                                </t-button>
                                <t-upload theme="custom" :action="uploadConfig.action" :headers="uploadConfig.headers"
                                    :show-file-list="false" accept="image/*" @success="handleUploadImageSuccess">
                                    <t-button variant="outline"><template #icon><t-icon
                                                name="upload" /></template>上传</t-button>
                                </t-upload>
                            </div>

                            <div class="section-group" style="margin-top: 20px;">
                                <div class="section-header"><span>生成结果</span></div>
                                <div class="image-wrapper" v-if="currentStoryboard.image || currentStoryboard.imageUrl">
                                    <t-image :src="getImageUrl(currentStoryboard.image || currentStoryboard.imageUrl)"
                                        fit="contain" style="max-height: 200px; width: 100%; display: block;" />
                                    <div class="img-actions">
                                        <t-button shape="circle" size="small" variant="text"
                                            @click="previewImage(currentStoryboard.image || currentStoryboard.imageUrl)"><t-icon
                                                name="zoom-in" /></t-button>
                                        <t-button shape="circle" size="small" variant="text" theme="danger"><t-icon
                                                name="delete" /></t-button>
                                    </div>
                                </div>
                                <div v-else class="empty-text">暂无生成图片</div>
                            </div>
                        </div>
                        <t-empty v-else description="请选择一个镜头" style="margin-top: 40px" />
                    </t-tab-panel>

                    <t-tab-panel value="video" label="视频生成">
                        <div class="tab-content scrollable-content" v-if="currentStoryboard">

                            <div class="video-prompt-box">
                                {{ currentStoryboard.videoPrompt || currentStoryboard.visualDesc || '暂无提示词' }}
                            </div>

                            <div class="video-settings section-group">
                                <t-form-item label="视频模型">
                                    <t-select v-model="selectedVideoModel" :options="videoModelOptions" />
                                </t-form-item>
                                <t-form-item label="时长 (秒)">
                                    <t-slider v-model="videoDuration" :min="2" :max="10" />
                                </t-form-item>
                                <t-form-item label="参考图模式">
                                    <t-select v-model="referenceMode" :options="referenceModeOptions" />
                                </t-form-item>

                                <div class="reference-selector" v-if="referenceMode !== 'none'">
                                    <div class="label">选择参考图</div>
                                    <div class="ref-image-slot" :class="{ selected: true }">
                                        <t-image v-if="currentStoryboard.image || currentStoryboard.imageUrl"
                                            :src="getImageUrl(currentStoryboard.image || currentStoryboard.imageUrl)"
                                            fit="cover" style="width: 100%; height: 100%;" />
                                        <div v-else class="placeholder"><t-icon name="image" /></div>
                                    </div>
                                </div>

                                <t-button theme="primary" block size="large" :loading="generatingVideo"
                                    @click="generateVideo" style="margin-top: 16px;">
                                    <template #icon><t-icon name="video" /></template> 生成视频
                                </t-button>
                            </div>

                            <div class="video-list-area section-group" v-if="currentStoryboard.videoUrl">
                                <div class="section-header"><span>生成结果</span></div>
                                <div class="video-card">
                                    <video :src="getVideoUrl(currentStoryboard.videoUrl)" controls
                                        style="width: 100%; border-radius: 8px;"></video>
                                    <div class="video-actions">
                                        <t-tag theme="success" variant="light">已完成</t-tag>
                                        <div class="action-btns">
                                            <t-tooltip content="更新到时间线">
                                                <t-button size="small" variant="text" @click="addVideoToAssets">
                                                    <template #icon><t-icon name="layers" /></template>
                                                </t-button>
                                            </t-tooltip>
                                            <t-button size="small" variant="text" theme="danger" @click="deleteVideo">
                                                <template #icon><t-icon name="delete" /></template>
                                            </t-button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <t-empty v-else description="请选择一个镜头" style="margin-top: 40px" />
                    </t-tab-panel>

                    <t-tab-panel value="audio" label="音效配乐">
                        <div class="tab-content">
                            <t-empty description="音效与配乐功能开发中..." />
                        </div>
                    </t-tab-panel>

                    <t-tab-panel value="merge" label="视频合成">
                        <div class="tab-content">
                            <div class="section-group">
                                <div class="section-header"><span>合成记录</span></div>
                                <t-empty description="暂无合成记录" size="small" />
                            </div>
                            <t-button theme="primary" block @click="exportVideo">
                                <template #icon><t-icon name="layers" /></template> 开始合成当前时间线
                            </t-button>
                        </div>
                    </t-tab-panel>

                </t-tabs>
            </div>
        </div>

        <t-dialog v-model:visible="showSceneSelector" header="关联场景" width="500px">
            <t-list :split="true" style="max-height: 400px; overflow-y: auto">
                <t-list-item v-for="scene in sceneList" :key="scene.id" @click="linkSceneToShot(scene)"
                    style="cursor: pointer">
                    <t-list-item-meta :title="scene.name" :description="`${scene.location} · ${scene.time}`">
                        <template #image>
                            <t-image v-if="scene.visualPrompt" :src="getImageUrl(scene.visualPrompt)" fit="cover"
                                style="width: 50px; height: 50px; border-radius: 4px;" />
                            <t-icon v-else name="image" size="24px" style="color: #ccc" />
                        </template>
                    </t-list-item-meta>
                    <template #action>
                        <t-icon v-if="currentStoryboard?.sceneId === scene.id" name="check"
                            style="color: var(--td-brand-color)" />
                    </template>
                </t-list-item>
                <t-empty v-if="sceneList.length === 0" description="暂无场景数据" />
            </t-list>
        </t-dialog>

        <t-dialog v-model:visible="showCharacterSelector" header="选择角色" width="500px">
            <div class="char-selector-grid">
                <div v-for="char in availableCharacters" :key="char.id" class="char-item"
                    :class="{ selected: selectedCharacters.includes(char.id) }" @click="toggleCharacterInShot(char.id)">
                    <t-avatar :image="getImageUrl(char.avatarUrl)" size="large" />
                    <span>{{ char.name }}</span>
                    <div class="check" v-if="selectedCharacters.includes(char.id)"><t-icon name="check" /></div>
                </div>
            </div>
            <t-empty v-if="availableCharacters.length === 0" description="暂无角色" />
        </t-dialog>

        <t-dialog v-model:visible="showPropSelector" header="选择道具" width="500px">
            <div class="char-selector-grid">
                <div v-for="prop in availableProps" :key="prop.id" class="char-item"
                    :class="{ selected: selectedProps.includes(prop.id) }" @click="togglePropInShot(prop.id)">
                    <t-image :src="getImageUrl(prop.imageUrl)" fit="contain"
                        style="width: 50px; height: 50px; border-radius: 4px; background: #f9f9f9;" />
                    <span>{{ prop.name }}</span>
                    <div class="check" v-if="selectedProps.includes(prop.id)"><t-icon name="check" /></div>
                </div>
            </div>
            <t-empty v-if="availableProps.length === 0" description="暂无道具" />
        </t-dialog>

    </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import {
    ArrowLeftIcon, RefreshIcon, AddIcon, DeleteIcon, MagicIcon,
    UploadIcon, ZoomInIcon, VideoIcon, LinkIcon, LayersIcon,
    MoveIcon, AddCircleIcon, FilmIcon, CheckIcon, DownloadIcon, CloseIcon
} from 'tdesign-icons-vue-next'

// API
import { findProjects } from '@/api/projects'
import { getScriptsList } from '@/api/scripts'
import { getScenesList } from '@/api/scenes'
import { getCharactersList } from '@/api/characters'
import { getPropsList } from '@/api/props'
import { getShotsList, createShots, updateShots, deleteShots } from '@/api/shots'
import { getAssetsList, createAsset, deleteAsset } from '@/api/assets'
import { getImageUrl } from '@/utils/format'

// 组件
import VideoTimelineEditor from '@/components/editor/VideoTimelineEditor.vue'

const route = useRoute()
const router = useRouter()

// === 数据状态 ===
const loading = ref(false)
const saving = ref(false)
const dramaId = route.params.dramaId as string
const episodeNumber = Number(route.params.episodeNumber)

const project = ref<any>({})
const currentScriptId = ref<number | null>(null)
const storyboards = ref<any[]>([])
const currentStoryboardId = ref<string | number | null>(null)
const sceneList = ref<any[]>([])
const availableCharacters = ref<any[]>([])
const availableProps = ref<any[]>([])
const videoAssets = ref<any[]>([])
const timelineClips = ref<any[]>([])
const audioClips = ref<any[]>([])
const currentTime = ref(0)
const totalDuration = ref(60)

// === UI 状态 ===
const activeTab = ref('shot')
const showSceneSelector = ref(false)
const showCharacterSelector = ref(false)
const showPropSelector = ref(false)
const generatingImage = ref(false)
const generatingVideo = ref(false)
const selectedVideoModel = ref('kling')
const videoDuration = ref(5)
const currentPreviewUrl = ref('')
const selectedFrameType = ref('first') // 新增：帧类型
const referenceMode = ref('single')    // 新增：参考图模式

const timelineEditorRef = ref<any>(null)
const mainPlayerRef = ref<HTMLVideoElement | null>(null)

const videoModelOptions = [
    { label: '可灵 (Kling)', value: 'kling' },
    { label: 'Runway Gen-3', value: 'runway' },
    { label: 'Luma', value: 'luma' }
]

const referenceModeOptions = [
    { label: '单图参考', value: 'single' },
    { label: '首尾帧', value: 'first_last' },
    { label: '纯文本', value: 'none' }
]

// === Computed ===
const currentStoryboard = computed(() => storyboards.value.find(s => s.id === currentStoryboardId.value))
const currentScene = computed(() => {
    if (!currentStoryboard.value || !currentStoryboard.value.sceneId) return null
    return sceneList.value.find(s => s.id === currentStoryboard.value.sceneId)
})
const selectedCharacters = computed(() => currentStoryboard.value?.characterIds || [])
const selectedProps = computed(() => currentStoryboard.value?.propIds || [])

const getCharacterById = (id: number) => availableCharacters.value.find(c => c.id === id)
const getPropById = (id: number) => availableProps.value.find(p => p.id === id)

const getAuthToken = () => localStorage.getItem('token')
const uploadConfig = reactive({
    action: import.meta.env.VITE_API_URL + '/admin/v1/upload/singleUpload',
    headers: computed(() => ({ 'Authorization': `${getAuthToken()}` })),
})

// === 初始化 ===
const initData = async () => {
    loading.value = true
    try {
        const res = await findProjects(dramaId)
        if (res.code === 0) project.value = res.data

        const sceneRes = await getScenesList({ projectId: dramaId, pageSize: 100 })
        if (sceneRes.code === 0) sceneList.value = sceneRes.data?.list || []

        const charRes = await getCharactersList({ projectId: dramaId, pageSize: 100 })
        if (charRes.code === 0) availableCharacters.value = charRes.data?.list || []

        const propRes = await getPropsList({ projectId: dramaId, pageSize: 100 })
        if (propRes.code === 0) availableProps.value = propRes.data?.list || []

        await loadShotsData()
        await loadVideoAssets()

    } catch (e) { console.error(e) } finally { loading.value = false }
}

const loadShotsData = async () => {
    const scriptRes = await getScriptsList({ projectId: dramaId, page: 1, pageSize: 100 })
    const list = scriptRes.data?.list || []
    const targetScript = list.find((s: any) => Number(s.episodeNo) === episodeNumber)

    if (targetScript) {
        currentScriptId.value = targetScript.id
        const shotRes = await getShotsList({ scriptId: targetScript.id, pageSize: 1000 })
        if (shotRes.code === 0) {
            storyboards.value = shotRes.data?.list || shotRes.data || []
            if (storyboards.value.length > 0 && !currentStoryboardId.value) {
                currentStoryboardId.value = storyboards.value[0].id
            }
        }
    }
}

// const loadVideoAssets = async () => {
//     if (!currentScriptId.value) return
//     try {
//         const res = await getAssetsList({
//             projectId: dramaId,
//             episodeId: currentScriptId.value
//         })
//         if (res.code === 0) {
//             videoAssets.value = res.data.list || []
//         }
//     } catch {
//         if (videoAssets.value.length === 0) {
//             videoAssets.value = storyboards.value
//                 .filter(s => s.videoUrl)
//                 .map(s => ({
//                     id: `asset_${s.id}`,
//                     name: s.title || 'Video Clip',
//                     url: s.videoUrl,
//                     duration: s.duration || 5
//                 }))
//         }
//     }
// }

const loadVideoAssets = async () => {
    // 模拟加载过程
    // if(!currentScriptId.value) return 

    try {
        // 1. 尝试调用真实接口 (保留原有逻辑)
        // const res = await getAssetsList({ 
        //     projectId: dramaId, 
        //     episodeId: currentScriptId.value 
        // })
        // if (res.code === 0 && res.data.list.length > 0) {
        //     videoAssets.value = res.data.list
        //     return
        // }

        // 2. 👇👇👇 注入 Mock 测试数据 (如果没有接口数据) 👇👇👇
        console.log('注入素材库测试数据...')
        videoAssets.value = [
            {
                id: 'mock_1',
                name: '开场-街道空镜.mp4',
                // MDN 示例视频
                url: 'https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4',
                duration: 5,
                type: 'video'
            },
            {
                id: 'mock_2',
                name: '主角登场-特写.mp4',
                // W3Schools 示例视频 (Big Buck Bunny)
                url: 'https://www.w3schools.com/html/mov_bbb.mp4',
                duration: 10,
                type: 'video'
            },
            {
                id: 'mock_3',
                name: '动作场景-追逐.mp4',
                // Sintel Trailer
                url: 'https://media.w3.org/2010/05/sintel/trailer.mp4',
                duration: 52,
                type: 'video'
            },
            {
                id: 'mock_4',
                name: '海洋-氛围.mp4',
                // VideoJS 示例
                url: 'https://vjs.zencdn.net/v/oceans.mp4',
                duration: 46,
                type: 'video'
            }
        ]
        // 👆👆👆 注入结束 👆👆👆

        MessagePlugin.success('素材库已加载 (测试数据)')

    } catch (e) {
        console.error(e)
        MessagePlugin.error('加载素材失败')
    }
}

// === 关系管理 ===
const toggleCharacterInShot = async (charId: number) => {
    if (!currentStoryboard.value) return
    let ids = [...(currentStoryboard.value.characterIds || [])]
    const idx = ids.indexOf(charId)
    if (idx > -1) ids.splice(idx, 1)
    else ids.push(charId)
    currentStoryboard.value.characterIds = ids
    await saveStoryboardField('characterIds')
}

const togglePropInShot = async (propId: number) => {
    if (!currentStoryboard.value) return
    let ids = [...(currentStoryboard.value.propIds || [])]
    const idx = ids.indexOf(propId)
    if (idx > -1) ids.splice(idx, 1)
    else ids.push(propId)
    currentStoryboard.value.propIds = ids
    await saveStoryboardField('propIds')
}

// === 拖拽逻辑 ===
const handleDragStart = (e: DragEvent, item: any, type: 'storyboard' | 'asset') => {
    if (e.dataTransfer) {
        const videoUrl = type === 'asset' ? item.url : item.videoUrl
        const payload = {
            id: item.id,
            name: item.title || item.name,
            url: videoUrl,
            duration: item.duration || 5,
            type: 'video',
            storyboardId: type === 'storyboard' ? item.id : undefined
        }
        e.dataTransfer.setData('application/json', JSON.stringify(payload))
        e.dataTransfer.effectAllowed = 'copy'
    }
}

const handleTimelineDrop = (clipData: any) => {
    const newClip = {
        id: `clip_${Date.now()}`,
        assetId: clipData.id,
        storyboardId: clipData.storyboardId,
        name: clipData.name,
        url: clipData.url,
        start: clipData.startTime,
        duration: clipData.duration,
        type: 'video'
    }

    if (clipData.trackType === 'audio') audioClips.value.push(newClip)
    else timelineClips.value.push(newClip)

    const end = newClip.start + newClip.duration
    if (end > totalDuration.value) totalDuration.value = Math.ceil(end + 10)

    MessagePlugin.success('已添加到时间线')
}

// === 交互逻辑 ===
const goBack = () => router.back()
const loadData = () => { initData(); MessagePlugin.success('数据已刷新') }

const selectStoryboard = (id: number | string) => {
    currentStoryboardId.value = id
    if (timelineEditorRef.value?.scrollToClip) {
        timelineEditorRef.value.scrollToClip(id)
    }
}

const handleTimelineSelect = (clip: any) => {
    if (clip.storyboardId) currentStoryboardId.value = clip.storyboardId
    currentPreviewUrl.value = clip.url
    if (mainPlayerRef.value) {
        mainPlayerRef.value.currentTime = 0;
        mainPlayerRef.value.play();
    }
}

const removeClipFromTimeline = (clipId: string) => {
    const idx = timelineClips.value.findIndex(c => c.id === clipId)
    if (idx > -1) timelineClips.value.splice(idx, 1)
}

const updateCurrentTime = (time: number) => {
    currentTime.value = time
    const activeClip = timelineClips.value.find(c => time >= c.start && time < c.start + c.duration)
    if (activeClip && activeClip.url) {
        if (currentPreviewUrl.value !== activeClip.url) currentPreviewUrl.value = activeClip.url
        const offset = time - activeClip.start
        if (mainPlayerRef.value && Math.abs(mainPlayerRef.value.currentTime - offset) > 0.5) {
            mainPlayerRef.value.currentTime = offset
        }
    }
}

const updateTimelineClips = (clips: any[]) => timelineClips.value = clips

// === 分镜业务逻辑 ===
const handleAddStoryboard = async () => {
    if (!currentScriptId.value) return MessagePlugin.warning('未关联剧本')
    const newShot = {
        projectId: Number(dramaId),
        scriptId: currentScriptId.value,
        title: `新镜头 ${storyboards.value.length + 1}`,
        duration: 3,
        shotType: '中景',
        cameraAngle: '平视',
        cameraMove: '固定'
    }
    try {
        const res = await createShots(newShot)
        if (res.code === 0) {
            MessagePlugin.success('添加成功'); await loadShotsData()
            if (storyboards.value.length > 0) currentStoryboardId.value = storyboards.value[storyboards.value.length - 1].id
        }
    } catch { MessagePlugin.error('添加失败') }
}

const handleDeleteStoryboard = async (shot: any) => {
    try {
        await deleteShots(shot.id)
        MessagePlugin.success('删除成功')
        const idx = storyboards.value.findIndex(s => s.id === shot.id)
        if (idx > -1) storyboards.value.splice(idx, 1)
        if (currentStoryboardId.value === shot.id) currentStoryboardId.value = storyboards.value.length > 0 ? storyboards.value[0].id : null
    } catch { MessagePlugin.error('删除失败') }
}

const linkSceneToShot = async (scene: any) => {
    if (!currentStoryboard.value) return
    currentStoryboard.value.sceneId = scene.id
    await saveStoryboardField('sceneId')
    showSceneSelector.value = false; MessagePlugin.success('已关联场景')
}

const saveStoryboardField = async (field: string) => {
    if (!currentStoryboard.value) return
    saving.value = true
    try { await updateShots(currentStoryboard.value.id, { ...currentStoryboard.value }) } catch { MessagePlugin.error('保存失败') } finally { saving.value = false }
}

const generateImage = async () => {
    if (!currentStoryboard.value) return
    generatingImage.value = true
    setTimeout(() => { generatingImage.value = false; MessagePlugin.success('生图任务已提交(Mock)') }, 1000)
}

const handleUploadImageSuccess = (ctx: any) => {
    if (ctx.response?.code === 0) {
        if (currentStoryboard.value) {
            const url = ctx.response.data.url || ctx.response.data.file_url
            const fullUrl = url.startsWith('http') ? url : import.meta.env.VITE_API_URL.replace(/\/admin\/v1$/, '') + url
            currentStoryboard.value.image = fullUrl
            currentStoryboard.value.imageUrl = fullUrl
            saveStoryboardField('image')
            MessagePlugin.success('上传成功')
        }
    }
}

const generateVideo = async () => {
    if (!currentStoryboard.value) return
    generatingVideo.value = true
    setTimeout(() => { generatingVideo.value = false; MessagePlugin.success('视频生成任务已提交') }, 1500)
}

const addVideoToAssets = async () => {
    if (!currentStoryboard.value?.videoUrl) return
    try {
        await createAsset({
            projectId: dramaId,
            episodeId: currentScriptId.value,
            url: currentStoryboard.value.videoUrl,
            name: currentStoryboard.value.title,
            type: 'video',
            duration: currentStoryboard.value.duration
        })
        MessagePlugin.success('已添加到素材库')
        await loadVideoAssets()
    } catch { MessagePlugin.error('添加失败') }
}

const deleteVideo = async () => {
    if (currentStoryboard.value) {
        currentStoryboard.value.videoUrl = ''; await saveStoryboardField('videoUrl')
    }
}

const previewImage = (url: string) => window.open(url, '_blank')
const getVideoUrl = (url: string) => url ? (url.startsWith('http') ? url : import.meta.env.VITE_API_URL + url) : ''
const exportVideo = () => { MessagePlugin.info('导出合成视频功能开发中') }

onMounted(() => initData())
</script>

<style scoped lang="less">
.professional-editor {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--td-bg-color-container);
    color: var(--td-text-color-primary);

    .editor-header {
        height: 56px;
        background: #fff;
        border-bottom: 1px solid var(--td-component-stroke);
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0 16px;
        flex-shrink: 0;

        .header-title .title {
            font-weight: 700;
            color: var(--td-text-color-primary);
        }

        .status-text {
            font-size: 12px;
            color: var(--td-text-color-placeholder);
            display: flex;
            align-items: center;
            gap: 4px;
        }
    }

    .editor-main {
        flex: 1;
        display: flex;
        overflow: hidden;

        /* 左侧侧边栏 */
        .left-sidebar {
            width: 280px;
            background: #fff;
            border-right: 1px solid var(--td-component-stroke);
            display: flex;
            flex-direction: column;
            flex-shrink: 0;

            .storyboard-panel {
                flex: 1;
                min-height: 0;
                display: flex;
                flex-direction: column;
                border-bottom: 1px solid var(--td-component-stroke);

                .panel-header {
                    padding: 12px;
                    border-bottom: 1px solid var(--td-component-stroke);
                    display: flex;
                    justify-content: space-between;
                    align-items: center;

                    h3 {
                        margin: 0;
                        font-size: 14px;
                        font-weight: 600;
                        color: var(--td-text-color-primary);
                    }
                }

                .storyboard-list {
                    flex: 1;
                    overflow-y: auto;
                    padding: 10px;
                    display: flex;
                    flex-direction: column;
                    gap: 8px;

                    .storyboard-item {
                        display: flex;
                        gap: 10px;
                        padding: 8px;
                        border-radius: 4px;
                        background: var(--td-bg-color-container);
                        border: 1px solid var(--td-component-stroke);
                        cursor: pointer;
                        transition: all 0.2s;
                        position: relative;

                        &:hover {
                            border-color: var(--td-brand-color);

                            .drag-hint {
                                opacity: 1;
                            }
                        }

                        &.active {
                            border-color: var(--td-brand-color);
                            background: var(--td-brand-color-light);
                        }

                        .shot-thumb {
                            width: 70px;
                            height: 45px;
                            background: #eee;
                            border-radius: 2px;
                            flex-shrink: 0;
                            overflow: hidden;
                            position: relative;
                            display: flex;
                            align-items: center;
                            justify-content: center;

                            .placeholder {
                                font-size: 12px;
                                color: #999;
                            }

                            .drag-hint {
                                position: absolute;
                                inset: 0;
                                background: rgba(0, 0, 0, 0.3);
                                color: #fff;
                                display: flex;
                                align-items: center;
                                justify-content: center;
                                opacity: 0;
                                transition: opacity 0.2s;
                            }
                        }

                        .shot-content {
                            flex: 1;
                            min-width: 0;
                            display: flex;
                            flex-direction: column;
                            justify-content: center;

                            .shot-title {
                                font-size: 12px;
                                font-weight: 600;
                                color: var(--td-text-color-primary);
                                white-space: nowrap;
                                overflow: hidden;
                                text-overflow: ellipsis;
                            }

                            .shot-desc {
                                font-size: 10px;
                                color: var(--td-text-color-secondary);
                                margin-top: 2px;
                            }
                        }
                    }
                }
            }

            .assets-panel {
                height: 40%;
                display: flex;
                flex-direction: column;
                background: var(--td-bg-color-secondarycontainer);

                .panel-header {
                    padding: 8px 12px;
                    display: flex;
                    justify-content: space-between;
                    align-items: center;

                    h3 {
                        margin: 0;
                        font-size: 13px;
                        color: var(--td-text-color-primary);
                    }
                }

                .assets-grid {
                    flex: 1;
                    overflow-y: auto;
                    padding: 10px;
                    display: grid;
                    grid-template-columns: repeat(2, 1fr);
                    gap: 10px;
                    align-content: start;

                    .asset-item {
                        background: #fff;
                        border-radius: 4px;
                        overflow: hidden;
                        border: 1px solid transparent;
                        cursor: grab;
                        position: relative;

                        &:hover {
                            border-color: var(--td-brand-color);
                            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

                            .hover-overlay {
                                display: flex;
                            }
                        }

                        .asset-thumb {
                            height: 60px;
                            background: #000;
                            position: relative;

                            video {
                                width: 100%;
                                height: 100%;
                                object-fit: cover;
                            }

                            .duration {
                                position: absolute;
                                right: 2px;
                                bottom: 2px;
                                background: rgba(0, 0, 0, 0.6);
                                color: #fff;
                                font-size: 10px;
                                padding: 1px 3px;
                                border-radius: 2px;
                            }

                            .hover-overlay {
                                position: absolute;
                                inset: 0;
                                background: rgba(0, 0, 0, 0.3);
                                display: none;
                                align-items: center;
                                justify-content: center;
                                color: #fff;
                                font-size: 20px;
                            }
                        }

                        .asset-name {
                            font-size: 10px;
                            padding: 4px;
                            white-space: nowrap;
                            overflow: hidden;
                            text-overflow: ellipsis;
                            color: var(--td-text-color-primary);
                        }
                    }
                }

                .empty-assets {
                    margin-top: 20px;
                }
            }
        }

        /* 中间工作区 */
        .center-workspace {
            flex: 1;
            display: flex;
            flex-direction: column;
            background-color: #1e1e1e;
            overflow: hidden;

            .preview-stage {
                flex: 1;
                display: flex;
                justify-content: center;
                align-items: center;
                border-bottom: 1px solid #333;
                background-image: radial-gradient(#333 1px, transparent 1px);
                background-size: 20px 20px;

                .player-container {
                    width: 90%;
                    height: 90%;
                    background: #000;
                    display: flex;
                    align-items: center;
                    justify-content: center;

                    .main-player {
                        width: 100%;
                        height: 100%;
                        object-fit: contain;
                    }

                    .player-placeholder {
                        text-align: center;
                        color: #666;

                        p {
                            margin-top: 10px;
                            font-size: 12px;
                        }
                    }
                }
            }

            .timeline-stage {
                height: 320px;
                flex-shrink: 0;
                background: #252525;
                border-top: 1px solid #333;
            }
        }

        /* 右侧属性面板 */
        .edit-panel {
            width: 360px;
            background: #fff;
            border-left: 1px solid var(--td-component-stroke);
            display: flex;
            flex-direction: column;
            flex-shrink: 0;

            .tab-content {
                padding: 16px;
                flex: 1;
                overflow-y: auto;

                &::-webkit-scrollbar {
                    width: 4px;
                }

                &::-webkit-scrollbar-thumb {
                    background: #e0e0e0;
                    border-radius: 2px;
                }
            }

            .section-group {
                margin-bottom: 24px;

                .section-header {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    margin-bottom: 12px;
                    font-size: 13px;
                    font-weight: 600;
                    color: var(--td-text-color-primary);
                    padding-left: 8px;
                    border-left: 3px solid var(--td-brand-color);
                }
            }

            .scene-card {
                border: 1px solid var(--td-component-stroke);
                border-radius: 6px;
                overflow: hidden;

                .scene-cover {
                    height: 120px;
                    width: 100%;
                    cursor: zoom-in;
                }

                .scene-info {
                    padding: 8px 10px;
                    background: var(--td-bg-color-secondarycontainer);

                    .scene-loc {
                        font-weight: 600;
                        font-size: 13px;
                        color: var(--td-text-color-primary);
                    }

                    .scene-meta {
                        font-size: 11px;
                        color: var(--td-text-color-secondary);
                        margin-top: 2px;
                    }
                }
            }

            .empty-box {
                border: 1px dashed var(--td-component-stroke);
                border-radius: 6px;
                height: 80px;
                display: flex;
                align-items: center;
                justify-content: center;
                gap: 6px;
                cursor: pointer;
                color: var(--td-text-color-placeholder);

                &:hover {
                    border-color: var(--td-brand-color);
                    color: var(--td-brand-color);
                }
            }

            .cast-list {
                display: flex;
                flex-wrap: wrap;
                gap: 8px;

                .cast-item {
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    width: 60px;
                    position: relative;

                    .cast-name {
                        font-size: 11px;
                        margin-top: 4px;
                        color: var(--td-text-color-secondary);
                        text-align: center;
                        width: 100%;
                        white-space: nowrap;
                        overflow: hidden;
                        text-overflow: ellipsis;
                    }

                    .remove-btn {
                        position: absolute;
                        top: 0;
                        right: 0;
                        background: rgba(0, 0, 0, 0.5);
                        color: #fff;
                        border-radius: 50%;
                        width: 16px;
                        height: 16px;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        font-size: 10px;
                        cursor: pointer;
                        opacity: 0;
                        transition: opacity 0.2s;
                    }

                    &:hover .remove-btn {
                        opacity: 1;
                    }
                }
            }

            .empty-text {
                font-size: 12px;
                color: var(--td-text-color-placeholder);
                padding: 10px;
                text-align: center;
                background: var(--td-bg-color-secondarycontainer);
                border-radius: 4px;
            }

            .video-prompt-box {
                padding: 10px;
                background: var(--td-bg-color-secondarycontainer);
                border-radius: 4px;
                font-size: 12px;
                color: var(--td-text-color-secondary);
                margin-bottom: 16px;
                border: 1px solid var(--td-component-stroke);
            }

            .reference-selector {
                margin-top: 16px;

                .label {
                    font-size: 12px;
                    margin-bottom: 8px;
                    font-weight: 500;
                }

                .ref-image-slot {
                    width: 120px;
                    height: 70px;
                    border: 2px dashed var(--td-component-stroke);
                    border-radius: 4px;
                    overflow: hidden;
                    cursor: pointer;
                    display: flex;
                    align-items: center;
                    justify-content: center;

                    &.selected {
                        border-color: var(--td-brand-color);
                        border-style: solid;
                    }

                    .placeholder {
                        color: #ccc;
                        font-size: 24px;
                    }
                }
            }

            .image-wrapper,
            .video-card {
                background: #000;
                border-radius: 4px;
                overflow: hidden;
                border: 1px solid var(--td-component-stroke);
                margin-top: 8px;

                .img-actions,
                .video-actions {
                    padding: 8px;
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    background: #fff;
                    border-top: 1px solid var(--td-component-stroke);
                }
            }
        }
    }
}

/* 角色/道具选择器样式 */
.char-selector-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    padding: 10px;

    .char-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        cursor: pointer;
        border: 2px solid transparent;
        padding: 10px;
        border-radius: 8px;
        position: relative;

        &:hover {
            background: var(--td-bg-color-secondarycontainer);
        }

        &.selected {
            border-color: var(--td-brand-color);
            background: var(--td-brand-color-light);
        }

        span {
            margin-top: 8px;
            font-size: 12px;
            font-weight: 500;
        }

        .check {
            position: absolute;
            top: 8px;
            right: 8px;
            color: var(--td-brand-color);
        }
    }
}
</style>