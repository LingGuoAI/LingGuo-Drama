
<template>
    <div class="shots-list">
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
                            <!-- 所属项目ID搜索 -->
                    <t-col :span="4">
                        <t-form-item label="短剧项目" name="projectId">
                            <t-select
                                v-model="searchInfo.projectId"
                                placeholder="请选择短剧项目"
                                clearable
                                :loading="projectsSelectLoading"
                            >
                                <t-option
                                    v-for="item in projectsSelectData"
                                    :key="item.id"
                                    :label="item.title"
                                    :value="item.id"
                                ></t-option>
                            </t-select>
                        </t-form-item>
                    </t-col>
                            <!-- 所属剧本/分集ID搜索 -->
                    <t-col :span="4">
                        <t-form-item label="剧本" name="scriptId">
                            <t-select
                                v-model="searchInfo.scriptId"
                                placeholder="请选择剧本"
                                clearable
                                :loading="scriptsSelectLoading"
                            >
                                <t-option
                                    v-for="item in scriptsSelectData"
                                    :key="item.id"
                                    :label="item.title"
                                    :value="item.id"
                                ></t-option>
                            </t-select>
                        </t-form-item>
                    </t-col>
                            <!-- 镜头序号搜索 -->
                    <t-col :span="4">
                        <t-form-item label="镜头序号" name="sequenceNo">
                            <t-input
                                v-model="searchInfo.sequenceNo"
                                placeholder="请输入镜头序号"
                                clearable
                                    type="number"
                            >
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
            <t-form-item label="短剧项目" name="projectId">
                <t-select v-model="formData.projectId" placeholder="请选择短剧项目" clearable
                    :loading="projectsSelectLoading" filterable
                    :status="!formData.projectId ? 'error' : 'default'">
                    <t-option v-for="item in projectsSelectData"
                        :key="item.id"
                        :label="item.title"
                        :value="item.id"></t-option>
                </t-select>
            </t-form-item>
            <t-form-item label="剧本" name="scriptId">
                <t-select v-model="formData.scriptId" placeholder="请选择剧本" clearable
                    :loading="scriptsSelectLoading" filterable
                    :status="!formData.scriptId ? 'error' : 'default'">
                    <t-option v-for="item in scriptsSelectData"
                        :key="item.id"
                        :label="item.title"
                        :value="item.id"></t-option>
                </t-select>
            </t-form-item>
            <t-form-item label="镜头序号" name="sequenceNo"><t-input v-model="formData.sequenceNo" placeholder="请输入镜头序号" type="number" clearable
                    :status="formData.sequenceNo === null || formData.sequenceNo === undefined ? 'error' : 'default'" />
            </t-form-item>
            <t-form-item label="景别: 全景/特写/中景" name="shotType">
                <t-input v-model="formData.shotType" clearable placeholder="请输入景别: 全景/特写/中景" :maxlength="50" show-word-limit />
            </t-form-item>
            <t-form-item label="运镜: 推/拉/摇/移" name="cameraMovement">
                <!-- 视频上传 -->
                <t-upload v-model="formData.cameraMovement" :action="uploadConfig.action" :headers="uploadConfig.headers"
                    accept="video/*" :auto-upload="true" :max="1" :size-limit="uploadConfig.videoSizeLimit"
                    :before-upload="beforeVideoUpload" tips="支持 mp4、mov、avi 格式，大小不超过 100MB" @success="handleVideoUploadSuccess" @fail="handleUploadFail">
                    <template #trigger>
                        <t-button theme="primary" variant="outline">
                            <template #icon><t-icon name="upload"></t-icon></template>
                            选择视频
                        </t-button>
                    </template>
                </t-upload>
            </t-form-item>
            <t-form-item label="视角: 俯视/平视" name="angle">
                <t-input v-model="formData.angle" clearable placeholder="请输入视角: 俯视/平视" :maxlength="50" show-word-limit />
            </t-form-item>
            <t-form-item label="台词/旁白" name="dialogue">
                <t-input v-model="formData.dialogue" clearable placeholder="请输入台词/旁白" :maxlength="255" show-word-limit />
            </t-form-item>
            <t-form-item label="画面描述" name="visualDesc">
                <t-input v-model="formData.visualDesc" clearable placeholder="请输入画面描述" :maxlength="255" show-word-limit />
            </t-form-item>
            <t-form-item label="氛围/环境描述" name="atmosphere">
                <t-input v-model="formData.atmosphere" clearable placeholder="请输入氛围/环境描述" :maxlength="255" show-word-limit />
            </t-form-item>
            <t-form-item label="绘画Prompt" name="imagePrompt">
                <!-- 单张图片上传 -->
                <div class="image-upload-container">
                    <!-- 已上传图片显示 -->
                    <div v-if="formData.imagePrompt && formData.imagePrompt.length > 0" class="uploaded-images">
                        <div v-for="(file, index) in formData.imagePrompt" :key="index" class="uploaded-item">
                            <div class="image-preview-wrapper">
                                <t-image-viewer v-if="file.url" :close-on-overlay="true" :images="[getImageUrl(file.url)]">
                                    <template #trigger="{ open }">
                                        <t-image
                                            :src="getImageUrl(file.url)"
                                            @click="open"
                                            fit="cover"
                                            class="image-preview"
                                            lazy
                                            error="图片加载失败">
                                        </t-image>
                                    </template>
                                </t-image-viewer>

                                <!-- 图片操作覆盖层 -->
                                <div class="image-overlay">
                                    <t-space>
                                        <t-button theme="primary" variant="text" size="small"
                                            @click="previewImage(file.url)" class="overlay-btn">
                                            <t-icon name="view"></t-icon>
                                        </t-button>
                                        <t-button theme="danger" variant="text" size="small"
                                            @click="handleImageRemove(index, 'imagePrompt')" class="overlay-btn">
                                            <t-icon name="delete"></t-icon>
                                        </t-button>
                                    </t-space>
                                </div>
                            </div>

                            <div class="image-info">
                                <div class="image-name">{{ file.name || '图片文件' }}</div>
                                <div class="image-size">{{ formatFileSize(file.size) }}</div>
                            </div>

                            <!-- 上传进度 -->
                            <t-progress v-if="file.percent !== undefined && file.percent < 100" :percentage="file.percent"
                                size="small" class="upload-progress"></t-progress>
                        </div>
                    </div>

                    <!-- 上传区域 -->
                    <t-upload v-show="!formData.imagePrompt || formData.imagePrompt.length === 0"
                        v-model="tempimagePromptList" :action="uploadConfig.action" :headers="uploadConfig.headers"
                        :data="uploadConfig.data" accept="image/*" :show-image-filename="false" :auto-upload="true" :max="1"
                        :size-limit="uploadConfig.sizeLimit" :format="uploadConfig.allowedFormats"
                        :before-upload="beforeUpload"
                        @success="(response) => handleImageUploadSuccess(response, 'imagePrompt')"
                        @fail="handleUploadFail" @progress="handleUploadProgress" class="upload-area">
                        <template #trigger>
                            <div class="upload-trigger">
                                <t-icon name="upload" size="32px"></t-icon>
                                <div class="upload-text">
                                    <div class="upload-title">点击上传图片</div>
                                    <div class="upload-desc">支持 jpg、jpeg、png、gif、webp 格式，大小不超过 5MB</div>
                                </div>
                            </div>
                        </template>
                    </t-upload>

                    <!-- 重新上传按钮 -->
                    <div v-if="formData.imagePrompt && formData.imagePrompt.length > 0" class="reupload-section">
                        <t-upload v-model="tempimagePromptReuploadList" :action="uploadConfig.action"
                            :headers="uploadConfig.headers" :data="uploadConfig.data" accept="image/*"
                            :show-image-filename="false" :auto-upload="true" :max="1" :size-limit="uploadConfig.sizeLimit"
                            :format="uploadConfig.allowedFormats" :before-upload="beforeReupload"
                            @success="(response) => handleReuploadSuccess(response, 'imagePrompt')"
                            @fail="handleUploadFail" class="reupload-component">
                            <template #trigger>
                                <t-button theme="default" variant="outline" size="small" :loading="uploading">
                                    <template #icon>
                                        <t-icon :name="uploading ? 'loading' : 'refresh'"></t-icon>
                                    </template>
                            {{ uploading ? '上传中...' : '重新上传' }}
                                </t-button>
                            </template>
                        </t-upload>
                    </div>
                </div>
            </t-form-item>
            <t-form-item label="视频生成Prompt" name="videoPrompt">
                <!-- 视频上传 -->
                <t-upload v-model="formData.videoPrompt" :action="uploadConfig.action" :headers="uploadConfig.headers"
                    accept="video/*" :auto-upload="true" :max="1" :size-limit="uploadConfig.videoSizeLimit"
                    :before-upload="beforeVideoUpload" tips="支持 mp4、mov、avi 格式，大小不超过 100MB" @success="handleVideoUploadSuccess" @fail="handleUploadFail">
                    <template #trigger>
                        <t-button theme="primary" variant="outline">
                            <template #icon><t-icon name="upload"></t-icon></template>
                            选择视频
                        </t-button>
                    </template>
                </t-upload>
            </t-form-item>
            <t-form-item label="音效/BGM提示词" name="audioPrompt">
                <t-input v-model="formData.audioPrompt" clearable placeholder="请输入音效/BGM提示词" :maxlength="255" show-word-limit />
            </t-form-item>
            <t-form-item label="分镜图" name="imageUrl">
                <!-- 单张图片上传 -->
                <div class="image-upload-container">
                    <!-- 已上传图片显示 -->
                    <div v-if="formData.imageUrl && formData.imageUrl.length > 0" class="uploaded-images">
                        <div v-for="(file, index) in formData.imageUrl" :key="index" class="uploaded-item">
                            <div class="image-preview-wrapper">
                                <t-image-viewer v-if="file.url" :close-on-overlay="true" :images="[getImageUrl(file.url)]">
                                    <template #trigger="{ open }">
                                        <t-image
                                            :src="getImageUrl(file.url)"
                                            @click="open"
                                            fit="cover"
                                            class="image-preview"
                                            lazy
                                            error="图片加载失败">
                                        </t-image>
                                    </template>
                                </t-image-viewer>

                                <!-- 图片操作覆盖层 -->
                                <div class="image-overlay">
                                    <t-space>
                                        <t-button theme="primary" variant="text" size="small"
                                            @click="previewImage(file.url)" class="overlay-btn">
                                            <t-icon name="view"></t-icon>
                                        </t-button>
                                        <t-button theme="danger" variant="text" size="small"
                                            @click="handleImageRemove(index, 'imageUrl')" class="overlay-btn">
                                            <t-icon name="delete"></t-icon>
                                        </t-button>
                                    </t-space>
                                </div>
                            </div>

                            <div class="image-info">
                                <div class="image-name">{{ file.name || '图片文件' }}</div>
                                <div class="image-size">{{ formatFileSize(file.size) }}</div>
                            </div>

                            <!-- 上传进度 -->
                            <t-progress v-if="file.percent !== undefined && file.percent < 100" :percentage="file.percent"
                                size="small" class="upload-progress"></t-progress>
                        </div>
                    </div>

                    <!-- 上传区域 -->
                    <t-upload v-show="!formData.imageUrl || formData.imageUrl.length === 0"
                        v-model="tempimageUrlList" :action="uploadConfig.action" :headers="uploadConfig.headers"
                        :data="uploadConfig.data" accept="image/*" :show-image-filename="false" :auto-upload="true" :max="1"
                        :size-limit="uploadConfig.sizeLimit" :format="uploadConfig.allowedFormats"
                        :before-upload="beforeUpload"
                        @success="(response) => handleImageUploadSuccess(response, 'imageUrl')"
                        @fail="handleUploadFail" @progress="handleUploadProgress" class="upload-area">
                        <template #trigger>
                            <div class="upload-trigger">
                                <t-icon name="upload" size="32px"></t-icon>
                                <div class="upload-text">
                                    <div class="upload-title">点击上传图片</div>
                                    <div class="upload-desc">支持 jpg、jpeg、png、gif、webp 格式，大小不超过 5MB</div>
                                </div>
                            </div>
                        </template>
                    </t-upload>

                    <!-- 重新上传按钮 -->
                    <div v-if="formData.imageUrl && formData.imageUrl.length > 0" class="reupload-section">
                        <t-upload v-model="tempimageUrlReuploadList" :action="uploadConfig.action"
                            :headers="uploadConfig.headers" :data="uploadConfig.data" accept="image/*"
                            :show-image-filename="false" :auto-upload="true" :max="1" :size-limit="uploadConfig.sizeLimit"
                            :format="uploadConfig.allowedFormats" :before-upload="beforeReupload"
                            @success="(response) => handleReuploadSuccess(response, 'imageUrl')"
                            @fail="handleUploadFail" class="reupload-component">
                            <template #trigger>
                                <t-button theme="default" variant="outline" size="small" :loading="uploading">
                                    <template #icon>
                                        <t-icon :name="uploading ? 'loading' : 'refresh'"></t-icon>
                                    </template>
                            {{ uploading ? '上传中...' : '重新上传' }}
                                </t-button>
                            </template>
                        </t-upload>
                    </div>
                </div>
            </t-form-item>
            <t-form-item label="最终视频片段" name="videoUrl">
                <!-- 视频上传 -->
                <t-upload v-model="formData.videoUrl" :action="uploadConfig.action" :headers="uploadConfig.headers"
                    accept="video/*" :auto-upload="true" :max="1" :size-limit="uploadConfig.videoSizeLimit"
                    :before-upload="beforeVideoUpload" tips="支持 mp4、mov、avi 格式，大小不超过 100MB" @success="handleVideoUploadSuccess" @fail="handleUploadFail">
                    <template #trigger>
                        <t-button theme="primary" variant="outline">
                            <template #icon><t-icon name="upload"></t-icon></template>
                            选择视频
                        </t-button>
                    </template>
                </t-upload>
            </t-form-item>
            <t-form-item label="配音/音效" name="audioUrl">
                <t-input v-model="formData.audioUrl" clearable placeholder="请输入配音/音效" :maxlength="1024" show-word-limit />
            </t-form-item>
            <t-form-item label="时长(毫秒, 原duration*1000)" name="durationMs"><t-input v-model="formData.durationMs" placeholder="请输入时长(毫秒, 原duration*1000)" type="number" clearable />
            </t-form-item>
            <t-form-item label="状态" name="status">
                <t-select v-model="formData.status" placeholder="请选择状态">
                    <t-option v-for="item in statusStatusOptions" :key="item.value" :label="item.label"
                        :value="item.value"></t-option>
                </t-select>
            </t-form-item>
        </t-form>
    </t-dialog>

        
<t-dialog v-model:visible="detailVisible" header="查看详情" width="600px" size="large" :footer="false" :close-btn="true"
                  :show-overlay="true" @close="detailVisible = false">
    <t-descriptions :column="1" layout="vertical" bordered
        :content-style="{ overflowWrap: 'break-word',whiteSpace:'normal' }">
        <t-descriptions-item label="短剧项目">
            <span v-if="detailData.projects">
                {{ detailData.projects.title }}
            </span>
            <span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="剧本">
            <span v-if="detailData.scripts">
                {{ detailData.scripts.title }}
            </span>
            <span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="镜头序号">
            <span v-if="detailData.sequenceNo !== null && detailData.sequenceNo !== undefined && detailData.sequenceNo !== ''">
    {{ detailData.sequenceNo }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="景别: 全景/特写/中景">
            <span v-if="detailData.shotType !== null && detailData.shotType !== undefined && detailData.shotType !== ''">
    {{ detailData.shotType }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="运镜: 推/拉/摇/移">
            <video
                    v-if="detailData.cameraMovement"
                    :src="getImageUrl(detailData.cameraMovement)"
                    style="width: 300px; height: 200px;"
                    controls
                    preload="metadata"
            />
            <span v-else>--</span>
        </t-descriptions-item>
        <t-descriptions-item label="视角: 俯视/平视">
            <span v-if="detailData.angle !== null && detailData.angle !== undefined && detailData.angle !== ''">
    {{ detailData.angle }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="台词/旁白">
            <span v-if="detailData.dialogue !== null && detailData.dialogue !== undefined && detailData.dialogue !== ''">
    {{ detailData.dialogue }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="画面描述">
            <span v-if="detailData.visualDesc !== null && detailData.visualDesc !== undefined && detailData.visualDesc !== ''">
    {{ detailData.visualDesc }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="氛围/环境描述">
            <span v-if="detailData.atmosphere !== null && detailData.atmosphere !== undefined && detailData.atmosphere !== ''">
    {{ detailData.atmosphere }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="绘画Prompt">
            <t-image-viewer
                    v-if="detailData.imagePrompt"
                    :close-on-overlay="true"
                    :images="[getImageUrl(detailData.imagePrompt)]"
            >
              <template #trigger="{ open }">
                <t-image
                        :src="getImageUrl(detailData.imagePrompt)"
                        @click="open"
                        fit="cover"
                        style="width: 200px; height: 200px; border-radius: 8px; cursor: pointer;"
                        lazy
                        error="图片加载失败"
                />
              </template>
            </t-image-viewer>
            <span v-else>--</span>
        </t-descriptions-item>
        <t-descriptions-item label="视频生成Prompt">
            <video
                    v-if="detailData.videoPrompt"
                    :src="getImageUrl(detailData.videoPrompt)"
                    style="width: 300px; height: 200px;"
                    controls
                    preload="metadata"
            />
            <span v-else>--</span>
        </t-descriptions-item>
        <t-descriptions-item label="音效/BGM提示词">
            <span v-if="detailData.audioPrompt !== null && detailData.audioPrompt !== undefined && detailData.audioPrompt !== ''">
    {{ detailData.audioPrompt }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="分镜图">
            <t-image-viewer
                    v-if="detailData.imageUrl"
                    :close-on-overlay="true"
                    :images="[getImageUrl(detailData.imageUrl)]"
            >
              <template #trigger="{ open }">
                <t-image
                        :src="getImageUrl(detailData.imageUrl)"
                        @click="open"
                        fit="cover"
                        style="width: 200px; height: 200px; border-radius: 8px; cursor: pointer;"
                        lazy
                        error="图片加载失败"
                />
              </template>
            </t-image-viewer>
            <span v-else>--</span>
        </t-descriptions-item>
        <t-descriptions-item label="最终视频片段">
            <video
                    v-if="detailData.videoUrl"
                    :src="getImageUrl(detailData.videoUrl)"
                    style="width: 300px; height: 200px;"
                    controls
                    preload="metadata"
            />
            <span v-else>--</span>
        </t-descriptions-item>
        <t-descriptions-item label="配音/音效">
            <span v-if="detailData.audioUrl !== null && detailData.audioUrl !== undefined && detailData.audioUrl !== ''">
    {{ detailData.audioUrl }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="时长(毫秒, 原duration*1000)">
            <span v-if="detailData.durationMs !== null && detailData.durationMs !== undefined && detailData.durationMs !== ''">
    {{ detailData.durationMs }}
</span>
<span v-else style="color: var(--td-text-color-placeholder);">--</span>
        </t-descriptions-item>
        <t-descriptions-item label="状态">
            <t-tag
                    :theme="getStatusTagTheme(detailData.status, statusStatusOptions)"
                    variant="light"
            >
              {{ getStatusLabel(detailData.status, statusStatusOptions) }}
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
    createShots,
    deleteShots,
    updateShots,
    findShots,
    getShotsList,
    getProjectsSelectList,
    getScriptsSelectList
} from '@/api/shots'
import { formatDate, getImageUrl } from '@/utils/format'

    defineOptions({
        name: 'ShotsList'
    })

    const router = useRouter()

    // ========== 状态选项定义 ==========
    // 状态状态选项
    const statusStatusOptions = ref([{"value":0,"label":"Pending"},{"value":1,"label":"Done"},{"value":2,"label":"Fail"}])
    // 获取token的方法
    const getAuthToken = () => {
        return localStorage.getItem('token')
    }

    // 上传配置
    const uploadConfig = reactive({
        action: import.meta.env.VITE_API_URL + '/admin/v1/upload/singleUpload',
        headers: computed(() => ({
            'Authorization': `${getAuthToken()}`,
        })),
        data: {},
        sizeLimit: 5 * 1024 * 1024, // 5MB
        videoSizeLimit: 100 * 1024 * 1024, // 100MB
        fileSizeLimit: 50 * 1024 * 1024, // 50MB
        allowedFormats: ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
    })

    // 上传相关状态
    const uploading = ref(false)
    const tempimagePromptList = ref([]) // 绘画Prompt临时上传列表
    const tempimagePromptReuploadList = ref([]) // 绘画Prompt重新上传临时列表
    const tempimageUrlList = ref([]) // 分镜图临时上传列表
    const tempimageUrlReuploadList = ref([]) // 分镜图重新上传临时列表

    // === 文件大小格式化 ===
    const formatFileSize = (size) => {
        if (!size) return '0 B'
        const units = ['B', 'KB', 'MB', 'GB']
        let index = 0
        while (size >= 1024 && index < units.length - 1) {
            size /= 1024
            index++
        }
        return `${size.toFixed(1)} ${units[index]}`
    }

    // === 预览图片 ===
    const previewImage = (url) => {
        // 这里可以使用图片查看器组件
        console.log('预览图片:', getImageUrl(url))
    }

    // === 上传前验证 ===
    const beforeUpload = (file) => {
        return beforeUploadValidation(file)
    }

    // === 上传进度处理 ===
    const handleUploadProgress = (progress) => {
        console.log('上传进度:', progress)
    }

    // === 重新上传前验证 ===
    const beforeReupload = (file) => {
        return beforeUploadValidation(file)
    }

    // === 通用上传验证 ===
    const beforeUploadValidation = (file) => {
        // 验证文件类型
        if (!uploadConfig.allowedFormats.includes(file.type)) {
            MessagePlugin.error('不支持的文件格式，请上传 jpg、jpeg、png、gif、webp 格式的图片')
            return false
        }

        // 验证文件大小
        if (file.size > uploadConfig.sizeLimit) {
            MessagePlugin.error(`文件大小不能超过 ${formatFileSize(uploadConfig.sizeLimit)}`)
            return false
        }

        // 验证token
        const token = getAuthToken()
        if (!token) {
            MessagePlugin.error('用户未登录，请重新登录后上传')
            return false
        }

        uploading.value = true
        return true
    }
    const handleUploadFail = (response) => {
        uploading.value = false
        console.error('上传失败:', response)

        let errorMessage = '上传失败，请重试'

        // 根据不同的错误状态码给出不同的提示
        if (response.status === 401) {
            errorMessage = '认证失败，请重新登录'
        } else if (response.status === 413) {
            errorMessage = '文件过大，请选择较小的文件'
        } else if (response.status === 415) {
            errorMessage = '不支持的文件格式'
        } else if (response.response && response.response.message) {
            errorMessage = response.response.message
        }

        MessagePlugin.error(errorMessage)
    }

    // === 图片删除处理 ===
    const handleImageRemove = (index, fieldName) => {
        if (Array.isArray(formData.value[fieldName])) {
            formData.value[fieldName].splice(index, 1)
        }
        MessagePlugin.success('图片已删除')
    }
    // === 重新上传成功处理 ===
    const handleReuploadSuccess = (response, fieldName) => {
        handleUploadSuccess(response, fieldName, true)
    }

    // === 通用上传成功处理 ===
    const handleUploadSuccess = (response, fieldName, isReupload = false, isMultiple = false) => {
        uploading.value = false
        console.log('上传成功响应完整数据:', response)

        // 判断响应结构和成功状态
        const isSuccess = response.response?.code === 200 || response.code === 200 || response.response?.code === 0 || response.code === 0

        if (isSuccess) {
            MessagePlugin.success(isReupload ? '重新上传成功' : '图片上传成功')

            // 获取响应数据，兼容不同的响应结构
            const responseData = response.response?.data || response.data
            console.log('解析后的响应数据:', responseData)

            if (responseData && responseData.file_url) {
                // 处理文件URL - 如果是相对路径则拼接完整地址
                let fileUrl = responseData.file_url

                if (fileUrl.startsWith('/')) {
                    // 拼接API基础地址
                    const baseUrl = import.meta.env.VITE_API_URL.replace(/\/admin\/v1$/, '').replace(/\/v1$/, '')
                    fileUrl = baseUrl + fileUrl
                    console.log('拼接后的完整URL:', fileUrl)
                }

                const fileInfo = {
                    url: fileUrl,
                    name: responseData.file_name || '图片',
                    size: responseData.file_size || 0
                }

                if (isReupload) {
                    // 重新上传：替换现有图片
                    formData.value[fieldName] = [fileInfo]
                    // 清空对应的临时列表
                    eval(`temp${fieldName}ReuploadList.value = []`)
                } else if (isMultiple) {
                    // 多张图片上传：添加到现有列表
                    if (!Array.isArray(formData.value[fieldName])) {
                        formData.value[fieldName] = []
                    }
                    formData.value[fieldName].push(fileInfo)
                    // 清空对应的临时列表
                    eval(`temp${fieldName}List.value = []`)
                } else {
                    // 单张图片上传：设置为新的图片
                    formData.value[fieldName] = [fileInfo]
                    // 清空对应的临时列表
                    eval(`temp${fieldName}List.value = []`)
                }
            } else {
                console.error('响应数据中缺少file_url字段:', responseData)
                MessagePlugin.error('上传成功但获取图片地址失败')
            }
        } else {
            const errorMsg = response.response?.message || response.response?.msg || response.message || response.msg || '图片上传失败'
            console.error('上传失败:', errorMsg)
            MessagePlugin.error(errorMsg)
        }
    }
    // === 单张图片上传成功处理 ===
    const handleImageUploadSuccess = (response, fieldName) => {
        handleUploadSuccess(response, fieldName, false)
    }

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
        scriptId: undefined,
        sequenceNo: undefined,
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
    // 剧本选择数据相关
    const scriptsSelectData = ref([])
    const scriptsSelectLoading = ref(false)

    // 获取剧本选择数据
    const getScriptsSelectData = async () => {
        scriptsSelectLoading.value = true
        try {
            const res = await getScriptsSelectList()
            if (res.code === 0) {
                // 不需要再进行数据格式化，因为后端已经返回了正确格式
                scriptsSelectData.value = res.data || []
            } else {
                scriptsSelectData.value = []
                console.error('获取剧本数据失败:', res.message)
            }
        } catch (error) {
            console.error('获取剧本数据失败:', error)
            scriptsSelectData.value = []
        } finally {
            scriptsSelectLoading.value = false
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
            title: '剧本',
            colKey: 'scriptId',
            sorter: false,
            cell: (h, { row }) => {
                const relationObj = row.scripts
                if (relationObj && relationObj.title) {
                    return relationObj.title
                }
                return '--'
            }
        },
        {
            title: '镜头序号',
            colKey: 'sequenceNo',
            sorter: false,
            cell: (h, { row }) => row.sequenceNo || '--'
        },
        {
            title: '景别: 全景/特写/中景',
            colKey: 'shotType',
            cell: (h, { row }) => row.shotType || '--'
        },
        {
            title: '运镜: 推/拉/摇/移',
            colKey: 'cameraMovement',
            width: 120,
            cell: (h, { row }) => {
                if (!row.cameraMovement) return '--'
                return h('video', {
                    src: getImageUrl(row.cameraMovement),
                    style: 'width: 80px; height: 80px; border-radius: 4px;',
                    muted: true,
                    preload: 'metadata'
                })
            }
        },
        {
            title: '视角: 俯视/平视',
            colKey: 'angle',
            cell: (h, { row }) => row.angle || '--'
        },
        {
            title: '台词/旁白',
            colKey: 'dialogue',
            cell: (h, { row }) => row.dialogue || '--'
        },
        {
            title: '画面描述',
            colKey: 'visualDesc',
            cell: (h, { row }) => row.visualDesc || '--'
        },
        {
            title: '氛围/环境描述',
            colKey: 'atmosphere',
            cell: (h, { row }) => row.atmosphere || '--'
        },
        {
            title: '绘画Prompt',
            colKey: 'imagePrompt',
            width: 120,
            cell: (h, { row }) => {
                if (!row.imagePrompt) return '--'
                return (
                    <t-image-viewer
                        closeOnOverlay
                        images={[getImageUrl(row.imagePrompt)]}
                        trigger={(h,{open}: {open: () => void}) =>
                            <t-image
                                src={getImageUrl(row.imagePrompt)}
                                onClick={open}
                                fit="cover"
                                style="width: 80px; height: 80px; border-radius: 4px; cursor: pointer;"
                                lazy
                                error="加载失败"
                            />
                        }>
                    </t-image-viewer>
                )
            }
        },
        {
            title: '视频生成Prompt',
            colKey: 'videoPrompt',
            width: 120,
            cell: (h, { row }) => {
                if (!row.videoPrompt) return '--'
                return h('video', {
                    src: getImageUrl(row.videoPrompt),
                    style: 'width: 80px; height: 80px; border-radius: 4px;',
                    muted: true,
                    preload: 'metadata'
                })
            }
        },
        {
            title: '音效/BGM提示词',
            colKey: 'audioPrompt',
            cell: (h, { row }) => row.audioPrompt || '--'
        },
        {
            title: '分镜图',
            colKey: 'imageUrl',
            width: 120,
            cell: (h, { row }) => {
                if (!row.imageUrl) return '--'
                return (
                    <t-image-viewer
                        closeOnOverlay
                        images={[getImageUrl(row.imageUrl)]}
                        trigger={(h,{open}: {open: () => void}) =>
                            <t-image
                                src={getImageUrl(row.imageUrl)}
                                onClick={open}
                                fit="cover"
                                style="width: 80px; height: 80px; border-radius: 4px; cursor: pointer;"
                                lazy
                                error="加载失败"
                            />
                        }>
                    </t-image-viewer>
                )
            }
        },
        {
            title: '最终视频片段',
            colKey: 'videoUrl',
            width: 120,
            cell: (h, { row }) => {
                if (!row.videoUrl) return '--'
                return h('video', {
                    src: getImageUrl(row.videoUrl),
                    style: 'width: 80px; height: 80px; border-radius: 4px;',
                    muted: true,
                    preload: 'metadata'
                })
            }
        },
        {
            title: '配音/音效',
            colKey: 'audioUrl',
            cell: (h, { row }) => row.audioUrl || '--'
        },
        {
            title: '时长(毫秒, 原duration*1000)',
            colKey: 'durationMs',
            sorter: false,
            cell: (h, { row }) => row.durationMs || '--'
        },
        {
            title: '状态',
            colKey: 'status',
            sorter: false,
            cell: (h, { row }) => {
                const option = statusStatusOptions.value.find(item => item.value === row.status)
                if (option) {
                    return (
                        <t-tag shape="round" theme={getStatusTagTheme(row.status, statusStatusOptions.value)} variant="light">
                            {option.label}
                        </t-tag>
                    );
                }
                return row.status
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
        scriptId: null,
        sequenceNo: null,
        shotType: '',
        cameraMovement: '',
        angle: '',
        dialogue: '',
        visualDesc: '',
        atmosphere: '',
        imagePrompt: [],
        videoPrompt: [],
        audioPrompt: '',
        imageUrl: [],
        videoUrl: [],
        audioUrl: '',
        durationMs: null,
        status: 0, // 状态字段默认第一个选项
    })

    // === 验证规则 ===
    const rules = reactive({
        projectId: [
            { required: true, message: '请输入所属项目ID', trigger: ['blur', 'change'] },
            { type: 'number', message: '所属项目ID必须是数字', trigger: ['blur', 'change'] }

        ],
        scriptId: [
            { required: true, message: '请输入所属剧本/分集ID', trigger: ['blur', 'change'] },
            { type: 'number', message: '所属剧本/分集ID必须是数字', trigger: ['blur', 'change'] }

        ],
        sequenceNo: [
            { required: true, message: '请输入镜头序号', trigger: ['blur', 'change'] },
            { type: 'number', message: '镜头序号必须是数字', trigger: ['blur', 'change'] }

        ],
        shotType: [
            { whitespace: true, message: '景别: 全景/特写/中景不能只包含空格', trigger: 'blur' },
            { max: 50, message: '景别: 全景/特写/中景长度不能超过50个字符', trigger: ['blur', 'change'] }

        ],
        cameraMovement: [
            { whitespace: true, message: '运镜: 推/拉/摇/移不能只包含空格', trigger: 'blur' },
            { max: 50, message: '运镜: 推/拉/摇/移长度不能超过50个字符', trigger: ['blur', 'change'] }

        ],
        angle: [
            { whitespace: true, message: '视角: 俯视/平视不能只包含空格', trigger: 'blur' },
            { max: 50, message: '视角: 俯视/平视长度不能超过50个字符', trigger: ['blur', 'change'] }

        ],
        dialogue: [
            { whitespace: true, message: '台词/旁白不能只包含空格', trigger: 'blur' },
            { max: 255, message: '台词/旁白长度不能超过255个字符', trigger: ['blur', 'change'] }

        ],
        visualDesc: [
            { whitespace: true, message: '画面描述不能只包含空格', trigger: 'blur' },
            { max: 255, message: '画面描述长度不能超过255个字符', trigger: ['blur', 'change'] }

        ],
        atmosphere: [
            { whitespace: true, message: '氛围/环境描述不能只包含空格', trigger: 'blur' },
            { max: 255, message: '氛围/环境描述长度不能超过255个字符', trigger: ['blur', 'change'] }

        ],
        imagePrompt: [

        ],
        videoPrompt: [
            { whitespace: true, message: '视频生成Prompt不能只包含空格', trigger: 'blur' },
            { max: 255, message: '视频生成Prompt长度不能超过255个字符', trigger: ['blur', 'change'] }

        ],
        audioPrompt: [
            { whitespace: true, message: '音效/BGM提示词不能只包含空格', trigger: 'blur' },
            { max: 255, message: '音效/BGM提示词长度不能超过255个字符', trigger: ['blur', 'change'] }

        ],
        imageUrl: [
            { type: 'url', message: '请输入正确的URL格式', trigger: ['blur', 'change'] }

        ],
        videoUrl: [
            { whitespace: true, message: '最终视频片段不能只包含空格', trigger: 'blur' },
            { max: 1024, message: '最终视频片段长度不能超过1024个字符', trigger: ['blur', 'change'] },
            { type: 'url', message: '请输入正确的URL格式', trigger: ['blur', 'change'] }

        ],
        audioUrl: [
            { whitespace: true, message: '配音/音效不能只包含空格', trigger: 'blur' },
            { max: 1024, message: '配音/音效长度不能超过1024个字符', trigger: ['blur', 'change'] },
            { type: 'url', message: '请输入正确的URL格式', trigger: ['blur', 'change'] }

        ],
        durationMs: [
            { type: 'number', message: '时长(毫秒, 原duration*1000)必须是数字', trigger: ['blur', 'change'] }

        ],
        status: [

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
            const res = await getShotsList(params)
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
        // 所属剧本/分集ID
        if (searchInfo.value.scriptId !== undefined && searchInfo.value.scriptId !== '') {
            params.scriptId = searchInfo.value.scriptId
        }
        // 镜头序号
        if (searchInfo.value.sequenceNo !== undefined && searchInfo.value.sequenceNo !== '') {
            params.sequenceNo = searchInfo.value.sequenceNo
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
            scriptId: undefined,
            sequenceNo: undefined,
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
            const res = await findShots(row.id)
            if (res.code === 0) {
                let data = res.data
                // 只对图片字段进行特殊处理
                if (data.imagePrompt && typeof data.imagePrompt === 'string' && data.imagePrompt.includes(',')) {
                    data.imagePrompt = data.imagePrompt.split(',').filter(url => url.trim()).map(url => url.trim())
                }
                // 只对图片字段进行特殊处理
                if (data.imageUrl && typeof data.imageUrl === 'string' && data.imageUrl.includes(',')) {
                    data.imageUrl = data.imageUrl.split(',').filter(url => url.trim()).map(url => url.trim())
                }

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
        drawerTitle.value = '新增镜头表'
        resetForm()
        getProjectsSelectData()
        getScriptsSelectData()
        drawerVisible.value = true
    }

    // 编辑
    const onEdit = async (row) => {
        try {
            const res = await findShots(row.id)
            if (res.code === 0) {
                formType.value = 'update'
                drawerTitle.value = '编辑镜头表'

                // 处理返回的数据，确保上传字段格式正确
                const data = res.data
                // 处理单张图片回显
                if (data.imagePrompt) {
                    if (typeof data.imagePrompt === 'string') {
                        data.imagePrompt = [{ url: data.imagePrompt, name: '图片' }]
                    } else if (Array.isArray(data.imagePrompt)) {
                        data.imagePrompt = data.imagePrompt.map(item => ({
                            ...item,
                            url: item.url || item
                        }))
                    } else {
                        data.imagePrompt = []
                    }
                } else {
                    data.imagePrompt = []
                }
                // 处理文件回显
                if (data.videoPrompt) {
                    if (typeof data.videoPrompt === 'string') {
                        data.videoPrompt = [{ url: data.videoPrompt, name: '文件' }]
                    } else if (!Array.isArray(data.videoPrompt)) {
                        data.videoPrompt = []
                    }
                } else {
                    data.videoPrompt = []
                }
                // 处理单张图片回显
                if (data.imageUrl) {
                    if (typeof data.imageUrl === 'string') {
                        data.imageUrl = [{ url: data.imageUrl, name: '图片' }]
                    } else if (Array.isArray(data.imageUrl)) {
                        data.imageUrl = data.imageUrl.map(item => ({
                            ...item,
                            url: item.url || item
                        }))
                    } else {
                        data.imageUrl = []
                    }
                } else {
                    data.imageUrl = []
                }
                // 处理文件回显
                if (data.videoUrl) {
                    if (typeof data.videoUrl === 'string') {
                        data.videoUrl = [{ url: data.videoUrl, name: '文件' }]
                    } else if (!Array.isArray(data.videoUrl)) {
                        data.videoUrl = []
                    }
                } else {
                    data.videoUrl = []
                }

                // 确保所有字符串字段都有默认值，避免null或undefined导致的trim()错误
                if (data.shotType === null || data.shotType === undefined) {
                    data.shotType = ''
                }
                if (data.cameraMovement === null || data.cameraMovement === undefined) {
                    data.cameraMovement = ''
                }
                if (data.angle === null || data.angle === undefined) {
                    data.angle = ''
                }
                if (data.dialogue === null || data.dialogue === undefined) {
                    data.dialogue = ''
                }
                if (data.visualDesc === null || data.visualDesc === undefined) {
                    data.visualDesc = ''
                }
                if (data.atmosphere === null || data.atmosphere === undefined) {
                    data.atmosphere = ''
                }
                if (data.videoPrompt === null || data.videoPrompt === undefined) {
                    data.videoPrompt = ''
                }
                if (data.audioPrompt === null || data.audioPrompt === undefined) {
                    data.audioPrompt = ''
                }
                if (data.videoUrl === null || data.videoUrl === undefined) {
                    data.videoUrl = ''
                }
                if (data.audioUrl === null || data.audioUrl === undefined) {
                    data.audioUrl = ''
                }

                formData.value = data
                getProjectsSelectData()
                getScriptsSelectData()
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
                    const res = await deleteShots(row.id)
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
            // 处理单张图片字段
            if (submitData.imagePrompt && Array.isArray(submitData.imagePrompt)) {
                if (submitData.imagePrompt.length > 0) {
                    const imageItem = submitData.imagePrompt[0]
                    submitData.imagePrompt = typeof imageItem === 'object' ? imageItem.url : imageItem
                } else {
                    submitData.imagePrompt = ''
                }
            }
            // 处理文件字段
            if (submitData.videoPrompt && Array.isArray(submitData.videoPrompt)) {
                submitData.videoPrompt = submitData.videoPrompt.map(item =>
                    typeof item === 'object' ? item.url : item
                ).filter(url => url)
            }
            // 处理单张图片字段
            if (submitData.imageUrl && Array.isArray(submitData.imageUrl)) {
                if (submitData.imageUrl.length > 0) {
                    const imageItem = submitData.imageUrl[0]
                    submitData.imageUrl = typeof imageItem === 'object' ? imageItem.url : imageItem
                } else {
                    submitData.imageUrl = ''
                }
            }
            // 处理文件字段
            if (submitData.videoUrl && Array.isArray(submitData.videoUrl)) {
                submitData.videoUrl = submitData.videoUrl.map(item =>
                    typeof item === 'object' ? item.url : item
                ).filter(url => url)
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
                res = await createShots(submitData)
            } else {
                res = await updateShots(submitData.id, submitData)
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
            scriptId: null,
            sequenceNo: null,
            shotType: '',
            cameraMovement: '',
            angle: '',
            dialogue: '',
            visualDesc: '',
            atmosphere: '',
            imagePrompt: [],
            videoPrompt: [],
            audioPrompt: '',
            imageUrl: [],
            videoUrl: [],
            audioUrl: '',
            durationMs: null,
            status: 0, // 状态字段默认第一个选项
        }

        // 重置所有临时上传列表
        tempimagePromptList.value = []
        tempimagePromptReuploadList.value = []
        tempimageUrlList.value = []
        tempimageUrlReuploadList.value = []

        // 清除验证状态
        nextTick(() => {
            formRef.value?.clearValidate()
        })
    }

    // 初始化
    const init = async () => {
        getProjectsSelectData()
        getScriptsSelectData()
        getTableData()
    }

    onMounted(() => {
        init()
    })
</script>

<style scoped>
    .shots-list {
        padding: 20px;
    }

    .search-form {
        margin-bottom: 20px;
    }
    /* 图片上传容器样式 */
    .image-upload-container {
        width: 100%;
    }

    /* 多图上传容器 */
    .multi-upload {
        border: 1px dashed var(--td-border-level-2-color);
        border-radius: 8px;
        padding: 12px;
        background: var(--td-bg-color-container);
    }

    /* 已上传图片显示区域 */
    .uploaded-images {
        margin-bottom: 16px;
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
    }

    .uploaded-item {
        display: flex;
        align-items: flex-start;
        padding: 12px;
        border: 1px solid var(--td-border-level-1-color);
        border-radius: 8px;
        background: var(--td-bg-color-container);
        margin-bottom: 8px;
        position: relative;
        width: 100%;
    }

    .uploaded-item-multi {
        position: relative;
        width: 100px;
        height: 100px;
    }

    .image-preview-wrapper {
        position: relative;
        margin-right: 12px;
        flex-shrink: 0;
    }

    .image-preview {
        width: 80px;
        height: 80px;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.3s ease;
    }

    .image-preview-small {
        width: 100px;
        height: 100px;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.3s ease;
    }

    .image-preview:hover,
    .image-preview-small:hover {
        transform: scale(1.05);
    }

    .image-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.6);
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        opacity: 0;
        transition: opacity 0.3s ease;
    }

    .image-delete-btn {
        position: absolute;
        top: -8px;
        right: -8px;
        z-index: 10;
    }

    .delete-btn {
        width: 24px;
        height: 24px;
        border-radius: 50%;
        background: var(--td-error-color);
        color: white;
        border: none;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 0;
        min-width: 24px;
    }

    .image-preview-wrapper:hover .image-overlay {
        opacity: 1;
    }

    .overlay-btn {
        color: white !important;
        border-color: white !important;
    }

    .overlay-btn:hover {
        background-color: rgba(255, 255, 255, 0.2) !important;
    }

    .image-info {
        flex: 1;
        min-width: 0;
    }

    .image-name {
        font-weight: 500;
        color: var(--td-text-color-primary);
        margin-bottom: 4px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .image-size {
        font-size: 12px;
        color: var(--td-text-color-placeholder);
    }

    .upload-progress {
        position: absolute;
        bottom: 4px;
        left: 12px;
        right: 12px;
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

    .upload-trigger {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 24px;
        text-align: center;
        width: 100%;
        height: 100%;
    }

    .upload-trigger .t-icon {
        color: var(--td-text-color-placeholder);
        margin-bottom: 12px;
        display: block;
    }

    .upload-text {
        line-height: 1.5;
        text-align: center;
        width: 100%;
    }

    .upload-title {
        font-size: 16px;
        font-weight: 500;
        color: var(--td-text-color-primary);
        margin-bottom: 4px;
        text-align: center;
        display: block;
        width: 100%;
    }

    .upload-desc {
        font-size: 12px;
        color: var(--td-text-color-placeholder);
        text-align: center;
        display: block;
        width: 100%;
        margin: 0 auto;
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

    .upload-trigger-small {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 100%;
        color: var(--td-text-color-placeholder);
    }

    /* 重新上传区域 */
    .reupload-section {
        margin-top: 8px;
        text-align: center;
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
        .shots-list {
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