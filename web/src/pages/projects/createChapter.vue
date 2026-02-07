<template>
    <div class="workflow-container">
        <div class="workflow-header">
            <div class="header-left">
                <t-button variant="text" shape="circle" @click="goBack">
                    <template #icon><t-icon name="arrow-left" /></template>
                </t-button>
                <div class="header-title">
                    <span class="title">{{ project?.title || '加载中...' }}</span>
                    <t-tag :theme="getStatusTheme(project?.status)" variant="light">{{ getStatusText(project?.status)
                        }}</t-tag>
                </div>
            </div>

            <div class="header-center">
                <t-steps :current="currentStep" readonly theme="dot" class="workflow-steps">
                    <t-step-item :title="`第 ${currentScriptNumber} 集剧本`" content="撰写剧本" />
                    <t-step-item title="分镜拆解" content="AI 拆分镜头" />
                    <t-step-item title="角色定妆" content="生成角色形象" />
                </t-steps>
            </div>

            <div class="header-right">
                <t-button theme="default" variant="outline" size="small" @click="loadProjectData">
                    <template #icon><t-icon name="refresh" /></template>
                </t-button>
            </div>
        </div>

        <div class="stage-area">

            <div v-show="currentStep === 0" class="stage-wrapper">
                <t-card :bordered="false" class="full-height-card">
                    <div v-if="!hasScriptContent && !showScriptInput" class="empty-state-wrapper">
                        <t-empty :description="`尚未创建第 ${currentScriptNumber} 集内容`">
                            <template #action>
                                <t-button theme="primary" size="large" @click="startCreateChapter">
                                    <template #icon><t-icon name="file-add" /></template>
                                    开始创作
                                </t-button>
                            </template>
                        </t-empty>
                    </div>

                    <div v-if="!hasScriptContent && showScriptInput" class="script-editor-container">
                        <div class="editor-toolbar">
                            <div class="toolbar-title">剧本编辑器</div>
                            <t-button theme="primary" variant="outline" @click="generateScriptByAI"
                                :loading="generatingScript">
                                <template #icon><t-icon name="magic" /></template>
                                AI 灵感生成
                            </t-button>
                        </div>

                        <t-textarea v-model="scriptContent" placeholder="请输入剧本内容，建议包含场景描述、人物对话等..."
                            class="script-textarea" :autosize="{ minRows: 15 }" :disabled="generatingScript" />

                        <div class="editor-footer">
                            <t-button theme="primary" size="large" @click="saveChapterScript"
                                :disabled="!scriptContent.trim() || generatingScript">
                                <template #icon><t-icon name="check" /></template>
                                保存章节
                            </t-button>
                        </div>
                    </div>

                    <div v-if="hasScriptContent" class="script-preview-container">
                        <div class="preview-header">
                            <div class="ph-left">
                                <h3>第 {{ currentScriptNumber }} 集剧本</h3>
                                <t-tag theme="success" variant="light">已保存</t-tag>
                            </div>
                            <t-button theme="default" variant="text" @click="editCurrentScript">
                                <template #icon><t-icon name="edit" /></template>修改
                            </t-button>
                        </div>

                        <div class="preview-content">
                            <t-textarea v-model="currentScript.content" readonly
                                :autosize="{ minRows: 10, maxRows: 20 }" class="readonly-textarea" />
                        </div>

                        <div class="step-actions">
                            <t-button theme="primary" size="large" @click="nextStep">
                                下一步：分镜拆解 <template #suffix><t-icon name="chevron-right" /></template>
                            </t-button>
                        </div>
                    </div>
                </t-card>
            </div>

            <div v-show="currentStep === 1" class="stage-wrapper">
                <t-card :bordered="false" class="full-height-card">
                    <template #header>
                        <div class="card-header-flex">
                            <div class="header-info">
                                <t-icon name="film" size="24px" style="color: var(--td-brand-color)" />
                                <span class="title">分镜列表</span>
                                <span class="subtitle" v-if="currentScript?.shots?.length">
                                    (共 {{ currentScript.shots.length }} 个镜头)
                                </span>
                            </div>
                            <div class="header-actions">
                                <t-button theme="default" @click="regenerateShots"
                                    :disabled="!currentScript?.shots?.length">
                                    <template #icon><t-icon name="refresh" /></template> 重新拆分
                                </t-button>
                                <t-button theme="primary" @click="parseShotsToCharacters" :loading="parsingCharacters">
                                    <template #icon><t-icon name="user-search" /></template> 解析角色
                                </t-button>
                            </div>
                        </div>
                    </template>

                    <div v-if="currentScript?.shots && currentScript.shots.length > 0" class="table-wrapper">
                        <t-table :data="currentScript.shots" :columns="shotColumns" row-key="id" stripe hover
                            :max-height="600">
                            <template #duration="{ row }">
                                {{ row.duration_ms ? row.duration_ms / 1000 : 0 }} 秒
                            </template>
                            <template #operation="{ row, rowIndex }">
                                <t-link theme="primary" @click="editShot(row, rowIndex)">编辑</t-link>
                            </template>
                        </t-table>

                        <div class="step-actions mt-4">
                            <t-button theme="default" @click="prevStep">上一步</t-button>
                            <t-button theme="primary" @click="nextStep" :disabled="!hasCharacters">
                                下一步：角色定妆 <template #suffix><t-icon name="chevron-right" /></template>
                            </t-button>
                        </div>
                    </div>

                    <div v-else class="empty-state-wrapper">
                        <t-empty description="剧本尚未拆分为分镜镜头">
                            <template #action>
                                <t-button theme="primary" size="large" @click="generateShots"
                                    :loading="generatingShots">
                                    <template #icon><t-icon name="magic" /></template>
                                    AI 智能拆分
                                </t-button>
                            </template>
                        </t-empty>
                    </div>
                </t-card>
            </div>

            <div v-show="currentStep === 2" class="stage-wrapper">
                <t-card :bordered="false" class="full-height-card">
                    <div class="toolbar-section">
                        <div class="toolbar-left">
                            <t-checkbox :checked="checkAll" :indeterminate="isIndeterminate" @change="handleSelectAll">
                                全选 ({{ selectedCharacterIds.length }}/{{ project?.characters?.length || 0 }})
                            </t-checkbox>
                        </div>
                        <div class="toolbar-right">
                            <span class="stat-text">已生成: {{ characterImagesCount }}</span>
                            <t-button theme="primary" variant="outline" :disabled="selectedCharacterIds.length === 0"
                                :loading="batchGenerating" @click="batchGenerateCharacterImages">
                                <template #icon><t-icon name="magic" /></template>
                                批量生成选中
                            </t-button>
                            <t-button theme="success" @click="finishWorkflow" :disabled="!allCharactersHaveImages">
                                完成创作 <template #icon><t-icon name="check-circle" /></template>
                            </t-button>
                        </div>
                    </div>

                    <div class="character-grid">
                        <div class="char-card add-card" @click="openAddCharacterDialog">
                            <div class="add-content">
                                <t-icon name="add" size="32px" />
                                <span>手动添加角色</span>
                            </div>
                        </div>

                        <div v-for="char in project?.characters" :key="char.id" class="char-card"
                            :class="{ 'is-selected': selectedCharacterIds.includes(char.id) }">
                            <div class="card-select">
                                <t-checkbox :checked="selectedCharacterIds.includes(char.id)"
                                    @change="() => toggleSelection(char.id)" />
                            </div>

                            <div class="char-image">
                                <t-image v-if="hasImage(char)" :src="getImageUrl(char)" fit="cover" class="img-box" />
                                <div v-else class="img-placeholder">
                                    <t-avatar size="large">{{ char.name[0] }}</t-avatar>
                                </div>
                            </div>

                            <div class="char-info">
                                <div class="info-head">
                                    <span class="name">{{ char.name }}</span>
                                    <t-tag size="small" :theme="getRoleTheme(char.role_type)">{{
                                        getRoleText(char.role_type)
                                        }}</t-tag>
                                </div>
                                <div class="desc text-ellipsis-2" :title="char.visual_prompt || char.appearance_desc">
                                    {{ char.visual_prompt || char.appearance_desc || '暂无描述' }}
                                </div>
                                <t-link theme="primary" size="small"
                                    @click="editCharacterDescription(char)">编辑描述</t-link>
                            </div>

                            <div class="char-actions">
                                <t-tooltip content="AI生成">
                                    <t-button shape="circle" size="small" theme="primary"
                                        :loading="generatingCharacterIds.includes(char.id)"
                                        @click="generateCharacterImage(char)">
                                        <t-icon name="magic" />
                                    </t-button>
                                </t-tooltip>
                                <t-tooltip content="上传图片">
                                    <t-button shape="circle" size="small" variant="outline"
                                        @click="openUploadDialog(char)">
                                        <t-icon name="upload" />
                                    </t-button>
                                </t-tooltip>
                                <t-tooltip content="从库选择">
                                    <t-button shape="circle" size="small" variant="outline"
                                        @click="openCharacterLibrary(char)">
                                        <t-icon name="folder-open" />
                                    </t-button>
                                </t-tooltip>
                                <t-popconfirm content="确认删除?" @confirm="deleteCharacter(char)">
                                    <t-button shape="circle" size="small" theme="danger" variant="text">
                                        <t-icon name="delete" />
                                    </t-button>
                                </t-popconfirm>
                            </div>
                        </div>
                    </div>
                </t-card>
            </div>

        </div>

        <t-dialog v-model:visible="addCharacterDialogVisible" header="添加新角色" width="500px" @confirm="addCharacter">
            <t-form :data="newCharacter" label-align="top">
                <t-form-item label="角色名称" name="name" required><t-input v-model="newCharacter.name" /></t-form-item>
                <t-form-item label="类型" name="role_type"><t-select v-model="newCharacter.role_type"
                        :options="roleOptions" /></t-form-item>
                <t-form-item label="外貌描述" name="visual_prompt"><t-textarea
                        v-model="newCharacter.visual_prompt" /></t-form-item>
            </t-form>
        </t-dialog>

    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import {
    ArrowLeftIcon, MagicIcon, FileAddIcon, EditIcon, CheckIcon,
    ChevronRightIcon, RefreshIcon, UserSearchIcon, AddIcon, DeleteIcon,
    UploadIcon, FolderOpenIcon, FilmIcon
} from 'tdesign-icons-vue-next'

// 1. API 导入 (使用 findProjects 和 updateProjects)
import { findProjects, updateProjects } from '@/api/projects'
// 假设有 generationAPI 用于AI操作
import { generationAPI } from '@/api/generation'
import { characterLibraryAPI } from '@/api/character-library'
import { getImageUrl } from '@/utils/format'

const route = useRoute()
const router = useRouter()

// === 核心状态 ===
const currentStep = ref(0)
const currentScriptNumber = ref(Number(route.params.episodeNumber) || 1)
const project = ref<any>({})
const scriptContent = ref("")
const showScriptInput = ref(false)

// === 加载状态 ===
const generatingScript = ref(false)
const generatingShots = ref(false)
const parsingCharacters = ref(false)
const batchGenerating = ref(false)
const saving = ref(false)

// === 角色选择 ===
const selectedCharacterIds = ref<number[]>([])
const generatingCharacterIds = ref<number[]>([])

// === Computed ===
// 获取当前正在编辑的剧本对象 (Script)
const currentScript = computed(() => {
    const scripts = project.value?.scripts || project.value?.episodes || []
    return scripts.find((ep: any) => (ep.episode_no || ep.episode_number) === currentScriptNumber.value) || {}
})

// 判断当前集是否有内容
const hasScriptContent = computed(() => {
    const content = currentScript.value.content || currentScript.value.script_content
    return !!(content && content.length > 0)
})

const hasCharacters = computed(() => (project.value?.characters || []).length > 0)

const checkAll = computed(() => {
    const total = project.value?.characters?.length || 0
    return total > 0 && selectedCharacterIds.value.length === total
})

const isIndeterminate = computed(() => {
    const total = project.value?.characters?.length || 0
    const selected = selectedCharacterIds.value.length
    return selected > 0 && selected < total
})

const characterImagesCount = computed(() => {
    return project.value?.characters?.filter((c: any) => c.avatar_url || c.image_url).length || 0
})

const allCharactersHaveImages = computed(() => {
    const chars = project.value?.characters || []
    return chars.length > 0 && chars.every((c: any) => !!(c.avatar_url || c.image_url))
})

// === Helpers ===
const getStatusTheme = (status: any) => status === 2 ? 'success' : 'primary'
const getStatusText = (status: any) => status === 2 ? '已完成' : '制作中'
const getRoleTheme = (role: string) => ({ main: 'danger', supporting: 'primary', minor: 'default' }[role] || 'default')
const getRoleText = (role: string) => ({ main: '主角', supporting: '配角', minor: '路人' }[role] || '未知')

// === Actions ===
const goBack = () => router.back()

const startCreateChapter = () => showScriptInput.value = true

// 加载项目数据
const loadProjectData = async () => {
    try {
        const id = route.params.id as string
        const res = await findProjects(id) // 使用 findProjects
        if (res.code === 0) {
            project.value = res.data
            // 回显内容
            if (currentScript.value) {
                scriptContent.value = currentScript.value.content || currentScript.value.script_content || ''
            }
        } else {
            MessagePlugin.error(res.message || '加载项目失败')
        }
    } catch (e) {
        MessagePlugin.error('网络异常')
    }
}

// 1. 剧本操作
const generateScriptByAI = async () => {
    if (!project.value?.title) return MessagePlugin.warning('项目信息缺失')
    generatingScript.value = true
    scriptContent.value = ''
    try {
        // 模拟流式生成，请对接你的真实API
        scriptContent.value = `[场景：${project.value.title} 第${currentScriptNumber.value}集]\n\n(日，内景，办公室)\n主角坐在椅子上，若有所思...\n`
        MessagePlugin.success('AI生成完成')
    } finally {
        generatingScript.value = false
    }
}

const saveChapterScript = async () => {
    generatingScript.value = true
    try {
        // 构造更新数据
        const scripts = [...(project.value.scripts || project.value.episodes || [])]
        const idx = scripts.findIndex((s: any) => (s.episode_no || s.episode_number) === currentScriptNumber.value)

        const newScriptData = {
            episode_no: currentScriptNumber.value,
            title: `第${currentScriptNumber.value}集`,
            content: scriptContent.value,
            status: 1 // Draft/Generating
        }

        if (idx > -1) scripts[idx] = { ...scripts[idx], ...newScriptData }
        else scripts.push(newScriptData)

        // 调用 updateProjects 保存
        const updatePayload = {
            id: project.value.id,
            scripts: scripts
        }
        await updateProjects(project.value.id, updatePayload)

        MessagePlugin.success('保存成功')
        await loadProjectData()
        showScriptInput.value = false
    } catch (e) {
        MessagePlugin.error('保存失败')
    } finally {
        generatingScript.value = false
    }
}

const editCurrentScript = () => showScriptInput.value = true

// 2. 分镜
const shotColumns = [
    { colKey: 'index', title: '序号', width: 60, cell: (h, { rowIndex }) => rowIndex + 1 },
    { colKey: 'visual_desc', title: '画面描述', ellipsis: true },
    { colKey: 'dialogue', title: '台词', ellipsis: true },
    { colKey: 'duration', title: '时长', width: 100, cell: 'duration' },
    { colKey: 'operation', title: '操作', width: 100, fixed: 'right', cell: 'operation' }
]

const generateShots = async () => {
    generatingShots.value = true
    try {
        // 请替换为你的真实分镜拆分API
        setTimeout(async () => {
            MessagePlugin.success('分镜拆分成功')
            generatingShots.value = false
            await loadProjectData()
        }, 1500)
    } catch { generatingShots.value = false }
}

const regenerateShots = () => {
    const confirm = DialogPlugin.confirm({
        header: '重新拆分',
        body: '确定重新拆分吗？将覆盖现有数据',
        onConfirm: () => {
            confirm.hide()
            generateShots()
        }
    })
}

const parseShotsToCharacters = async () => {
    parsingCharacters.value = true
    // 请替换为你的真实角色解析API
    setTimeout(() => {
        parsingCharacters.value = false
        MessagePlugin.success('角色解析完成')
        loadProjectData()
    }, 1000)
}

// 3. 角色
const toggleSelection = (id: number) => {
    const idx = selectedCharacterIds.value.indexOf(id)
    if (idx > -1) selectedCharacterIds.value.splice(idx, 1)
    else selectedCharacterIds.value.push(id)
}

const handleSelectAll = (checked: boolean) => {
    if (checked) selectedCharacterIds.value = (project.value?.characters || []).map((c: any) => c.id)
    else selectedCharacterIds.value = []
}

const generateCharacterImage = (char: any) => {
    generatingCharacterIds.value.push(char.id)
    // 请替换为你的真实生图API
    setTimeout(() => {
        const idx = generatingCharacterIds.value.indexOf(char.id)
        if (idx > -1) generatingCharacterIds.value.splice(idx, 1)
        MessagePlugin.success('生成任务已提交')
    }, 1000)
}

const batchGenerateCharacterImages = () => {
    batchGenerating.value = true
    setTimeout(() => {
        batchGenerating.value = false
        MessagePlugin.success('批量生成中')
    }, 1500)
}

// Dialogs
const addCharacterDialogVisible = ref(false)
const newCharacter = ref({ name: '', role_type: 'supporting', visual_prompt: '' })
const roleOptions = [{ label: '主角', value: 'main' }, { label: '配角', value: 'supporting' }]
const openAddCharacterDialog = () => addCharacterDialogVisible.value = true
const addCharacter = async () => {
    addCharacterDialogVisible.value = false
    MessagePlugin.success('添加成功')
    loadProjectData()
}

// Navigation
const nextStep = () => { if (currentStep.value < 2) currentStep.value++ }
const prevStep = () => { if (currentStep.value > 0) currentStep.value-- }
const finishWorkflow = () => {
    MessagePlugin.success('创作流程结束')
    router.push({ name: 'ProjectDetail', params: { id: project.value.id } })
}

// Mock Functions
const editShot = (row: any, index: number) => { MessagePlugin.info('编辑分镜') }
const openUploadDialog = (char: any) => { MessagePlugin.info('上传图片') }
const openCharacterLibrary = (char: any) => { MessagePlugin.info('从库选择') }
const editCharacterDescription = (char: any) => { MessagePlugin.info('编辑描述') }
const deleteCharacter = (char: any) => { MessagePlugin.success('已删除') }

// Initialize
onMounted(loadProjectData)
</script>

<style scoped lang="less">
/* 基础布局 */
.workflow-container {
    min-height: 100vh;
    background: var(--td-bg-color-container-gray);
    display: flex;
    flex-direction: column;
}

/* 顶部 Header */
.workflow-header {
    background: #fff;
    height: 64px;
    padding: 0 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    flex-shrink: 0;
    z-index: 10;

    .header-left {
        display: flex;
        align-items: center;
        gap: 16px;
        width: 280px;

        .header-title {
            display: flex;
            align-items: center;
            gap: 8px;

            .title {
                font-weight: 700;
                font-size: 16px;
            }
        }
    }

    .header-center {
        flex: 1;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
    }

    .header-right {
        width: 280px;
        display: flex;
        justify-content: flex-end;
    }
}

/* 美化后的步骤条样式 */
.workflow-steps {
    max-width: 600px;
    width: 100%;

    :deep(.t-steps-item) {
        .t-steps-item__title {
            font-size: 14px;
            color: var(--td-text-color-secondary);
            font-weight: 500;
        }

        .t-steps-item__description {
            font-size: 12px;
            color: var(--td-text-color-placeholder);
            margin-top: 4px;
        }

        .t-steps-item__icon--dot {
            width: 10px;
            height: 10px;
            border: 2px solid transparent;
        }
    }

    /* 当前进行中步骤的高亮样式 */
    :deep(.t-steps-item--process) {
        .t-steps-item__title {
            color: var(--td-brand-color);
            font-weight: 700;
            transform: scale(1.05);
            transition: all 0.3s;
        }

        .t-steps-item__description {
            color: var(--td-brand-color);
            opacity: 0.8;
        }

        .t-steps-item__icon--dot {
            background-color: var(--td-brand-color);
            box-shadow: 0 0 0 3px var(--td-brand-color-light);
        }
    }

    /* 已完成步骤的样式 */
    :deep(.t-steps-item--finish) {
        .t-steps-item__title {
            color: var(--td-text-color-primary);
        }

        .t-steps-item__icon--dot {
            background-color: var(--td-brand-color);
        }

        .t-steps-item__inner::after {
            background-color: var(--td-brand-color);
        }
    }
}

/* 主内容区 */
.stage-area {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
}

.stage-wrapper {
    height: 100%;
}

.full-height-card {
    height: 100%;
    display: flex;
    flex-direction: column;

    :deep(.t-card__body) {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }
}

/* 剧本编辑 */
.script-editor-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 16px;

    .editor-toolbar {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .toolbar-title {
            font-weight: 600;
            font-size: 16px;
        }
    }

    .script-textarea {
        flex: 1;

        :deep(textarea) {
            height: 100% !important;
            font-family: monospace;
            line-height: 1.8;
        }
    }

    .editor-footer {
        display: flex;
        justify-content: flex-end;
    }
}

/* 剧本预览 */
.script-preview-container {
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 16px;

    .preview-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-bottom: 12px;
        border-bottom: 1px solid var(--td-component-stroke);

        .ph-left {
            display: flex;
            gap: 12px;
            align-items: center;

            h3 {
                margin: 0;
            }
        }
    }

    .preview-content {
        flex: 1;
        background: var(--td-bg-color-secondarycontainer);
        border-radius: 8px;

        .readonly-textarea :deep(textarea) {
            background: transparent;
            border: none;
            height: 100% !important;
            resize: none;
        }
    }
}

.step-actions {
    display: flex;
    justify-content: center;
    gap: 16px;
    padding-top: 16px;
}

/* 分镜 Header */
.card-header-flex {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-info {
        display: flex;
        align-items: center;
        gap: 8px;

        .title {
            font-weight: 700;
            font-size: 16px;
        }

        .subtitle {
            color: var(--td-text-color-secondary);
            font-size: 12px;
        }
    }

    .header-actions {
        display: flex;
        gap: 12px;
    }
}

/* 角色定妆 */
.toolbar-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px;
    background: var(--td-bg-color-secondarycontainer);
    border-radius: 6px;
    margin-bottom: 20px;

    .stat-text {
        margin-right: 16px;
        color: var(--td-text-color-secondary);
        font-size: 12px;
    }

    .toolbar-right {
        display: flex;
        align-items: center;
        gap: 12px;
    }
}

.character-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 20px;
    overflow-y: auto;
    padding: 4px;
}

.char-card {
    background: #fff;
    border: 1px solid var(--td-border-level-1-color);
    border-radius: 8px;
    overflow: hidden;
    position: relative;
    transition: all 0.2s;

    &:hover {
        transform: translateY(-2px);
        box-shadow: var(--td-shadow-2);

        .char-actions {
            opacity: 1;
        }
    }

    &.is-selected {
        border-color: var(--td-brand-color);
        background: var(--td-brand-color-light);
    }

    &.add-card {
        border: 2px dashed var(--td-component-stroke);
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        min-height: 320px;

        &:hover {
            border-color: var(--td-brand-color);
            color: var(--td-brand-color);
        }

        .add-content {
            text-align: center;
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 8px;
        }
    }

    .card-select {
        position: absolute;
        top: 8px;
        right: 8px;
        z-index: 2;
    }

    .char-image {
        height: 200px;
        background: var(--td-bg-color-secondarycontainer);
        display: flex;
        align-items: center;
        justify-content: center;
        overflow: hidden;

        .img-box {
            width: 100%;
            height: 100%;
        }

        .img-placeholder {
            color: var(--td-text-color-disabled);
        }
    }

    .char-info {
        padding: 12px;

        .info-head {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 6px;

            .name {
                font-weight: 700;
                font-size: 14px;
            }
        }

        .desc {
            font-size: 12px;
            color: var(--td-text-color-secondary);
            height: 36px;
            margin-bottom: 8px;
        }
    }

    .char-actions {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        padding: 8px;
        display: flex;
        gap: 8px;
        opacity: 0;
        transition: opacity 0.2s;
        background: linear-gradient(to bottom, rgba(0, 0, 0, 0.3), transparent);
    }
}

.empty-state-wrapper {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
}

.mt-4 {
    margin-top: 16px;
}

.text-ellipsis-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
</style>