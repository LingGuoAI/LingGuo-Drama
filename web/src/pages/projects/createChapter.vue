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
                <t-button theme="default" variant="outline" size="small" @click="initData">
                    <template #icon><t-icon name="refresh" /></template>
                </t-button>
            </div>
        </div>

        <div class="stage-area" v-loading="loading">

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

                    <div v-if="showScriptInput" class="script-editor-container">
                        <div class="editor-toolbar">
                            <div class="toolbar-title">剧本编辑器 (第 {{ currentScriptNumber }} 集)</div>
                            <t-button theme="primary" variant="outline" @click="generateScriptByAI"
                                :loading="generatingScript">
                                <template #icon><t-icon name="magic" /></template>
                                AI 灵感生成
                            </t-button>
                        </div>

                        <t-textarea v-model="scriptContent" placeholder="请输入剧本内容，建议包含场景描述、人物对话等..."
                            class="script-textarea" :autosize="{ minRows: 15 }" :disabled="generatingScript" />

                        <div class="editor-footer">
                            <t-button theme="default" style="margin-right: 12px" @click="cancelEdit">取消</t-button>
                            <t-button theme="primary" @click="handleSaveScript" :loading="saving"
                                :disabled="!scriptContent.trim()">
                                <template #icon><t-icon name="check" /></template>
                                保存章节
                            </t-button>
                        </div>
                    </div>

                    <div v-if="hasScriptContent && !showScriptInput" class="script-preview-container">
                        <div class="preview-header">
                            <div class="ph-left">
                                <h3>第 {{ currentScriptNumber }} 集剧本</h3>
                                <t-tag theme="success" variant="light">已保存</t-tag>
                            </div>
                            <t-button theme="primary" variant="text" @click="enterEditMode">
                                <template #icon><t-icon name="edit" /></template>修改剧本
                            </t-button>
                        </div>

                        <div class="preview-content">
                            <t-textarea :value="currentScriptData.content" readonly class="readonly-textarea"
                                :autosize="{ minRows: 10 }" />
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
                                <span class="subtitle" v-if="currentScriptData?.shots?.length">
                                    (共 {{ currentScriptData.shots.length }} 个镜头)
                                </span>
                            </div>
                            <div class="header-actions">
                                <t-button theme="default" @click="regenerateShots"
                                    :disabled="!currentScriptData?.shots?.length">
                                    <template #icon><t-icon name="refresh" /></template> 重新拆分
                                </t-button>
                                <t-button theme="primary" @click="parseShotsToCharacters" :loading="parsingCharacters">
                                    <template #icon><t-icon name="user-search" /></template> 解析角色
                                </t-button>
                            </div>
                        </div>
                    </template>

                    <div v-if="currentScriptData?.shots && currentScriptData.shots.length > 0" class="table-wrapper">
                        <t-table :data="currentScriptData.shots" :columns="shotColumns" row-key="id" stripe hover
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
                        <div class="toolbar-right">
                            <t-button theme="default" @click="prevStep">上一步</t-button>
                        </div>
                    </div>
                </t-card>
            </div>

        </div>

    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import {
    ArrowLeftIcon, MagicIcon, FileAddIcon, EditIcon, CheckIcon,
    ChevronRightIcon, RefreshIcon, UserSearchIcon, FilmIcon
} from 'tdesign-icons-vue-next'

// API
import { findProjects } from '@/api/projects'
import {
    createScripts,
    updateScripts,
    getScriptsList,
    findScripts // 新增引入
} from '@/api/scripts'
// import { generationAPI } from '@/api/generation'

const route = useRoute()
const router = useRouter()

// === 核心状态 ===
const loading = ref(false)
const currentStep = ref(0)
const project = ref<any>({})
const currentScriptData = ref<any>({}) // 存储当前章节的完整数据（含内容、分镜等）
const currentScriptNumber = ref(Number(route.params.episodeNumber) || 1)

// 编辑器状态
const showScriptInput = ref(false)
const scriptContent = ref("")
const generatingScript = ref(false)
const saving = ref(false)

// 分镜状态
const generatingShots = ref(false)
const parsingCharacters = ref(false)

// === Computed ===
const hasScriptContent = computed(() => {
    // 只有当 currentScriptData 有 ID 且有内容时，才算有剧本
    return !!(currentScriptData.value?.id && currentScriptData.value?.content)
})

const hasCharacters = computed(() => (project.value?.characters || []).length > 0)

// === 初始化逻辑 ===
const initData = async () => {
    loading.value = true
    try {
        await Promise.all([
            loadProjectInfo(),
            loadScriptDetail()
        ])
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

// 1. 加载项目基础信息
const loadProjectInfo = async () => {
    const res = await findProjects(route.params.id as string)
    if (res.code === 0) {
        project.value = res.data
    }
}

// 2. 加载当前章节详情 (核心修复)
const loadScriptDetail = async () => {
    // 先尝试通过 getScriptsList 获取当前项目的脚本列表，找到对应集数的 ID
    // 这一步是因为路由参数只有 projectId 和 episodeNumber，没有 scriptId
    try {
        const listRes = await getScriptsList({
            projectId: route.params.id,
            page: 1,
            pageSize: 100
        })

        if (listRes.code === 0) {
            const list = Array.isArray(listRes.data) ? listRes.data : (listRes.data?.list || [])
            // 查找匹配当前集数的脚本
            const targetScript = list.find((s: any) => Number(s.episodeNo) === currentScriptNumber.value)

            if (targetScript) {
                // 如果找到了，调用 findScripts 获取完整详情（包含 content）
                const detailRes = await findScripts(targetScript.id)
                if (detailRes.code === 0) {
                    currentScriptData.value = detailRes.data
                    // 如果有内容，默认显示预览
                    if (detailRes.data.content) {
                        showScriptInput.value = false
                    }
                }
            } else {
                // 如果没找到，说明该集尚未创建，重置数据
                currentScriptData.value = {}
                showScriptInput.value = false // 显示 Empty State
            }
        }
    } catch (e) {
        MessagePlugin.error('加载章节数据失败')
    }
}

// === 剧本操作逻辑 ===

// 点击“开始创作”
const startCreateChapter = () => {
    scriptContent.value = ''
    showScriptInput.value = true
}

// 点击“修改”
const enterEditMode = () => {
    scriptContent.value = currentScriptData.value.content || ''
    showScriptInput.value = true
}

// 取消编辑
const cancelEdit = () => {
    showScriptInput.value = false
    // 如果是新建取消，清空内容
    if (!currentScriptData.value.id) {
        scriptContent.value = ''
    }
}

// AI 生成
const generateScriptByAI = async () => {
    if (!project.value?.title) return MessagePlugin.warning('项目信息缺失')
    generatingScript.value = true
    try {
        // Mock
        setTimeout(() => {
            scriptContent.value = `[第${currentScriptNumber.value}集]\n\n场景：${project.value.title}...\n\n(主角缓缓登场)`
            generatingScript.value = false
            MessagePlugin.success('生成成功')
        }, 1000)
    } catch { generatingScript.value = false }
}

// 保存/创建剧本 (核心修复)
const handleSaveScript = async () => {
    if (!scriptContent.value.trim()) return MessagePlugin.warning('内容不能为空')
    saving.value = true

    try {
        const payload = {
            projectId: project.value.id,
            episodeNo: currentScriptNumber.value,
            title: `第${currentScriptNumber.value}集`, // 默认标题，也可增加输入框让用户填
            content: scriptContent.value,
            isLocked: 0
        }

        // 判断是 更新 还是 创建
        if (currentScriptData.value?.id) {
            // Update
            await updateScripts(currentScriptData.value.id, payload)
            MessagePlugin.success('更新成功')
        } else {
            // Create
            await createScripts(payload)
            MessagePlugin.success('创建成功')
        }

        // 保存成功后，重新加载详情
        await loadScriptDetail()

        // 确保退出编辑模式
        showScriptInput.value = false

    } catch (e) {
        MessagePlugin.error('操作失败')
    } finally {
        saving.value = false
    }
}

// === 分镜逻辑 ===
const shotColumns = [
    { colKey: 'index', title: '序号', width: 60, cell: (h, { rowIndex }) => rowIndex + 1 },
    { colKey: 'visual_desc', title: '画面描述', ellipsis: true },
    { colKey: 'dialogue', title: '台词', ellipsis: true },
    { colKey: 'duration', title: '时长', width: 100, cell: 'duration' },
    { colKey: 'operation', title: '操作', width: 100, fixed: 'right', cell: 'operation' }
]

const generateShots = async () => {
    if (!currentScriptData.value?.id) return MessagePlugin.warning('请先保存剧本')

    generatingShots.value = true
    // try {
    //     // 调用 AI 拆分接口，传递 scriptId
    //     await generationAPI.generateShots({ script_id: currentScriptData.value.id })

    //     MessagePlugin.success('拆分成功')
    //     // 重新加载详情以获取最新的 shots
    //     await loadScriptDetail()
    // } catch (e) {
    //     MessagePlugin.error('拆分失败')
    // } finally {
    //     generatingShots.value = false
    // }
}

const regenerateShots = () => {
    const confirm = DialogPlugin.confirm({
        header: '重新拆分',
        body: '确定重新拆分吗？当前分镜数据将被覆盖',
        onConfirm: () => {
            confirm.hide()
            generateShots()
        }
    })
}

// === 通用 Helper ===
const goBack = () => router.back()
const getStatusTheme = (s: any) => s === 2 ? 'success' : 'primary'
const getStatusText = (s: any) => s === 2 ? '已完成' : '制作中'
const nextStep = () => { if (currentStep.value < 2) currentStep.value++ }
const prevStep = () => { if (currentStep.value > 0) currentStep.value-- }

// ... 其他占位函数 (parseShotsToCharacters, editShot 等) ...
const parseShotsToCharacters = () => { }
const editShot = () => { }

// Lifecycle
onMounted(() => {
    initData()
})
</script>

<style scoped lang="less">
/* 样式保持不变，复用之前的 CSS */
.workflow-container {
    min-height: 100vh;
    background: var(--td-bg-color-container-gray);
    display: flex;
    flex-direction: column;
}

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
}

/* ... 其他样式 ... */
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

.script-editor-container,
.script-preview-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.script-textarea,
.readonly-textarea {
    flex: 1;

    :deep(textarea) {
        height: 100% !important;
        resize: none;
        font-family: monospace;
        line-height: 1.8;
    }
}

.readonly-textarea :deep(textarea) {
    background: var(--td-bg-color-secondarycontainer);
    border: none;
}

.preview-header,
.editor-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .ph-left {
        display: flex;
        gap: 12px;
        align-items: center;

        h3 {
            margin: 0;
        }
    }
}

.editor-footer {
    display: flex;
    justify-content: flex-end;
}

.step-actions {
    display: flex;
    justify-content: center;
    gap: 16px;
    padding-top: 16px;
}

.empty-state-wrapper {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>