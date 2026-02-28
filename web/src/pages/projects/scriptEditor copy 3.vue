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
                            :class="{ active: String(currentStoryboardId) === String(shot.id) }"
                            @click="selectStoryboard(shot.id)" draggable="true"
                            @dragstart="handleDragStart($event, shot, 'storyboard')">
                            <div class="shot-thumb">
                                <t-image v-if="shot.image || shot.imageUrl"
                                    :src="getImageUrl(shot.image || shot.imageUrl)" fit="cover"
                                    style="width: 100%; height: 100%;" />
                                <div v-else class="drag-hint"><t-icon name="move" /></div>
                            </div>
                            <div class="shot-content">
                                <div class="shot-header">
                                    <div class="shot-title-row">
                                        <span class="shot-number">#{{ shot.sequenceNo || index + 1 }}</span>
                                        <span class="shot-title" :title="shot.title">{{ shot.title || '未命名镜头' }}</span>
                                    </div>
                                    <div class="shot-actions">
                                        <span class="shot-duration">{{ (shot.durationMs || 3000) / 1000 }}s</span>
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
                                        <t-image-viewer v-if="currentScene.visualPrompt" :close-on-overlay="true"
                                            :images="[getImageUrl(currentScene.visualPrompt)]">
                                            <template #trigger="{ open }">
                                                <t-image :src="getImageUrl(currentScene.visualPrompt)" fit="cover"
                                                    class="scene-cover" @click.stop="open" style="cursor: zoom-in;" lazy
                                                    error="图片加载失败" />
                                            </template>
                                        </t-image-viewer>

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
                                            <t-image-viewer v-if="getCharacterById(charId)?.avatarUrl"
                                                :close-on-overlay="true"
                                                :images="[getImageUrl(getCharacterById(charId)?.avatarUrl)]">
                                                <template #trigger="{ open }">
                                                    <t-avatar :image="getImageUrl(getCharacterById(charId)?.avatarUrl)"
                                                        size="medium" shape="circle" @click.stop="open"
                                                        style="cursor: zoom-in;" />
                                                </template>
                                            </t-image-viewer>
                                            <t-avatar v-else size="medium" shape="circle">{{
                                                getCharacterById(charId)?.name?.[0] || '?' }}</t-avatar>

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
                                            <t-image-viewer v-if="getPropById(propId)?.imageUrl"
                                                :close-on-overlay="true"
                                                :images="[getImageUrl(getPropById(propId)?.imageUrl)]">
                                                <template #trigger="{ open }">
                                                    <t-image :src="getImageUrl(getPropById(propId)?.imageUrl)"
                                                        fit="contain"
                                                        style="width: 40px; height: 40px; border-radius: 4px; background: #eee; cursor: zoom-in;"
                                                        @click.stop="open" lazy error="加载失败" />
                                                </template>
                                            </t-image-viewer>
                                            <t-icon v-else name="image" size="24px" style="color: #ccc" />

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
                                                    @change="saveStoryboardField" />
                                            </t-form-item>
                                        </t-col>
                                        <t-col :span="6">
                                            <t-form-item label="视角">
                                                <t-select v-model="currentStoryboard.angle" size="small"
                                                    :options="['平视', '俯视', '仰视', '侧视', '航拍'].map(v => ({ label: v, value: v }))"
                                                    @change="saveStoryboardField" />
                                            </t-form-item>
                                        </t-col>
                                    </t-row>
                                    <t-form-item label="运镜" style="margin-top: 10px;">
                                        <t-select v-model="currentStoryboard.cameraMovement" size="small"
                                            :options="['固定镜头', '推镜', '拉镜', '摇镜', '移镜', '跟镜', '环绕'].map(v => ({ label: v, value: v }))"
                                            @change="saveStoryboardField" />
                                    </t-form-item>
                                    <t-form-item label="时长 (秒)" style="margin-top: 10px;">
                                        <t-slider :value="(currentStoryboard.durationMs || 3000) / 1000" :min="1"
                                            :max="60" @change="updateShotDurationMs" />
                                    </t-form-item>
                                </div>

                                <div class="section-group">
                                    <div class="section-header"><span>叙事内容</span></div>
                                    <t-form-item label="动作 (Action)">
                                        <t-textarea v-model="currentStoryboard.action" :autosize="{ minRows: 2 }"
                                            placeholder="角色做的动作..." @blur="saveStoryboardField" />
                                    </t-form-item>
                                    <t-form-item label="结果 (Result)">
                                        <t-textarea v-model="currentStoryboard.result" :autosize="{ minRows: 2 }"
                                            placeholder="动作导致的结果..." @blur="saveStoryboardField" />
                                    </t-form-item>
                                    <t-form-item label="对白 (Dialogue)">
                                        <t-textarea v-model="currentStoryboard.dialogue" :autosize="{ minRows: 2 }"
                                            placeholder="角色台词..." @blur="saveStoryboardField" />
                                    </t-form-item>
                                    <t-form-item label="画面描述 (Visual)">
                                        <t-textarea v-model="currentStoryboard.visualDesc" :autosize="{ minRows: 3 }"
                                            placeholder="详细的画面描述..." @blur="saveStoryboardField" />
                                    </t-form-item>
                                    <t-form-item label="氛围 (Atmosphere)">
                                        <t-textarea v-model="currentStoryboard.atmosphere" :autosize="{ minRows: 2 }"
                                            placeholder="光影、色调、气氛..." @blur="saveStoryboardField" />
                                    </t-form-item>
                                </div>

                                <t-divider />

                                <div class="section-group">
                                    <div class="section-header"><span>音频设置</span></div>
                                    <t-form-item label="音效与配乐 (Audio Prompt)">
                                        <t-textarea v-model="currentStoryboard.audioPrompt" :autosize="{ minRows: 2 }"
                                            placeholder="例如：开门声、脚步声、悲伤的钢琴曲..." @blur="saveStoryboardField" />
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
                                <div class="section-header">
                                    <span>AI 绘画提示词</span>
                                    <t-button theme="primary" variant="text" size="small" :loading="extractingPrompt"
                                        @click="extractFramePrompt">
                                        提取提示词
                                    </t-button>
                                </div>
                                <t-textarea v-model="currentFramePromptText" :rows="4" placeholder="输入英文提示词..."
                                    @blur="saveStoryboardField" />
                            </div>

                            <div class="action-bar">
                                <t-button theme="primary" :loading="generatingImage" @click="generateFrameImage">
                                    <template #icon><t-icon name="magic" /></template> 生成画面
                                </t-button>
                                <t-upload theme="custom" :action="uploadConfig.action" :headers="uploadConfig.headers"
                                    :show-file-list="false" accept="image/*" :before-upload="beforeUpload"
                                    @success="handleUploadImageSuccess" @fail="handleUploadFail">
                                    <t-button variant="outline" :loading="uploadingImage">
                                        <template #icon><t-icon name="upload" /></template>上传
                                    </t-button>
                                </t-upload>
                            </div>

                            <div class="section-group" style="margin-top: 20px;">
                                <div class="section-header"><span>生成结果 ({{ currentFrameImages.length }})</span></div>

                                <div v-if="selectedFrameType === 'action'" class="grid-entry-card"
                                    @click="showGridEditor = true">
                                    <t-icon name="add" size="24px" />
                                    <span>创建动作序列 (宫格图)</span>
                                </div>

                                <div class="image-grid-list" v-if="currentFrameImages.length > 0">
                                    <div v-for="img in currentFrameImages" :key="img.id" class="image-grid-item">
                                        <t-image :src="getImageUrl(img.url || img.imageUrl)" fit="cover" class="img" />

                                        <div class="img-overlay">
                                            <div class="actions-wrapper">
                                                <t-image-viewer :close-on-overlay="true"
                                                    :images="[getImageUrl(img.url || img.imageUrl)]">
                                                    <template #trigger="{ open }">
                                                        <div class="icon-btn" @click.stop="open" title="预览大图">
                                                            <t-icon name="zoom-in" size="18px" style="color: #fff;" />
                                                        </div>
                                                    </template>
                                                </t-image-viewer>
                                                <div class="icon-btn danger" @click.stop="deleteImage(img)"
                                                    title="删除图片">
                                                    <t-icon name="delete" size="18px" style="color: #fff;" />
                                                </div>
                                            </div>
                                            <div class="crop-btn" v-if="selectedFrameType === 'action'"
                                                @click.stop="openCropDialog(img)" title="前往裁剪九宫格">
                                                <t-icon name="cut" size="14px" />
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div v-else-if="selectedFrameType !== 'action'" class="empty-text">暂无图片</div>
                            </div>
                        </div>
                        <t-empty v-else description="请选择一个镜头" style="margin-top: 40px" />
                    </t-tab-panel>

                    <t-tab-panel value="video" label="视频生成">
                        <div class="tab-content scrollable-content" v-if="currentStoryboard">

                            <div class="video-prompt-box">
                                {{ currentStoryboard.videoPrompt || currentStoryboard.imagePrompt ||
                                    currentStoryboard.visualDesc || '暂无提示词' }}
                            </div>

                            <div class="video-settings section-group">
                                <t-form-item label="视频模型">
                                    <t-select v-model="selectedVideoModel" placeholder="请选择视频模型">
                                        <t-option v-for="model in videoModelCapabilities" :key="model.id"
                                            :value="model.id" :label="model.name">
                                            <div
                                                style="display: flex; justify-content: space-between; align-items: center;">
                                                <span>{{ model.name }}</span>
                                                <div class="model-tags">
                                                    <t-tag v-if="model.supportMultipleImages" size="small"
                                                        theme="success" variant="light"
                                                        style="margin-left: 4px">多图</t-tag>
                                                    <t-tag v-if="model.supportFirstLastFrame" size="small"
                                                        theme="primary" variant="light"
                                                        style="margin-left: 4px">首尾帧</t-tag>
                                                    <t-tag size="small" theme="default" variant="light"
                                                        style="margin-left: 4px">最多{{ model.maxImages }}张</t-tag>
                                                </div>
                                            </div>
                                        </t-option>
                                    </t-select>
                                </t-form-item>

                                <t-form-item label="时长 (秒)">
                                    <t-slider v-model="videoDuration" :min="2" :max="10" />
                                </t-form-item>

                                <t-form-item label="参考图模式"
                                    v-if="selectedVideoModel && availableReferenceModes.length > 0">
                                    <t-select v-model="referenceMode" placeholder="请选择参考图模式">
                                        <t-option v-for="mode in availableReferenceModes" :key="mode.value"
                                            :value="mode.value" :label="mode.label">
                                            <div
                                                style="display: flex; justify-content: space-between; align-items: center;">
                                                <span>{{ mode.label }}</span>
                                                <span v-if="mode.description"
                                                    style="color: var(--td-text-color-placeholder); font-size: 12px;">{{
                                                        mode.description }}</span>
                                            </div>
                                        </t-option>
                                    </t-select>
                                </t-form-item>

                                <div class="reference-config-section" v-if="referenceMode && referenceMode !== 'none'">

                                    <div class="image-slots-container">
                                        <div v-if="referenceMode === 'single'" style="text-align: center">
                                            <div class="reference-mode-title">单图参考</div>
                                            <div style="display: inline-block">
                                                <t-upload theme="custom" :action="uploadConfig.action"
                                                    :headers="uploadConfig.headers" :show-file-list="false"
                                                    accept="image/*" :before-upload="beforeUpload"
                                                    @success="(ctx) => handleUploadRefSuccess(ctx, 'single')">
                                                    <div class="image-slot" :class="{ selected: !!singleRefImage }">
                                                        <t-image v-if="singleRefImage"
                                                            :src="getImageUrl(singleRefImage.url || singleRefImage.imageUrl)"
                                                            fit="cover" class="img" />
                                                        <div v-else class="image-slot-placeholder">
                                                            <t-icon name="add" size="24px" />
                                                            <div class="slot-hint">点击上传图片</div>
                                                        </div>
                                                    </div>
                                                </t-upload>
                                                <div class="image-slot-remove" v-if="singleRefImage"
                                                    @click.stop="removeSelectedImage(singleRefImage.id)">
                                                    <t-icon name="close" size="14px" />
                                                </div>
                                            </div>
                                        </div>

                                        <div v-else-if="referenceMode === 'first_last'" style="text-align: center">
                                            <div class="reference-mode-title">首尾帧</div>
                                            <div
                                                style="display: flex; gap: 20px; justify-content: center; align-items: center;">
                                                <div style="position: relative;">
                                                    <div class="frame-label">首帧</div>
                                                    <t-upload theme="custom" :action="uploadConfig.action"
                                                        :headers="uploadConfig.headers" :show-file-list="false"
                                                        accept="image/*" :before-upload="beforeUpload"
                                                        @success="(ctx) => handleUploadRefSuccess(ctx, 'first')">
                                                        <div class="image-slot" :class="{ selected: !!firstRefImage }">
                                                            <t-image v-if="firstRefImage"
                                                                :src="getImageUrl(firstRefImage.url || firstRefImage.imageUrl)"
                                                                fit="cover" class="img" />
                                                            <div v-else class="image-slot-placeholder">
                                                                <t-icon name="add" size="24px" />
                                                                <div class="slot-hint">点击上传首帧</div>
                                                            </div>
                                                        </div>
                                                    </t-upload>
                                                    <div class="image-slot-remove" v-if="firstRefImage"
                                                        @click.stop="removeSelectedImage(firstRefImage.id)">
                                                        <t-icon name="close" size="14px" />
                                                    </div>
                                                </div>

                                                <t-icon name="arrow-right" size="24px" style="color: #909399" />

                                                <div style="position: relative;">
                                                    <div class="frame-label">尾帧</div>
                                                    <t-upload theme="custom" :action="uploadConfig.action"
                                                        :headers="uploadConfig.headers" :show-file-list="false"
                                                        accept="image/*" :before-upload="beforeUpload"
                                                        @success="(ctx) => handleUploadRefSuccess(ctx, 'last')">
                                                        <div class="image-slot" :class="{ selected: !!lastRefImage }">
                                                            <t-image v-if="lastRefImage"
                                                                :src="getImageUrl(lastRefImage.url || lastRefImage.imageUrl)"
                                                                fit="cover" class="img" />
                                                            <div v-else class="image-slot-placeholder">
                                                                <t-icon name="add" size="24px" />
                                                                <div class="slot-hint">点击上传尾帧</div>
                                                            </div>
                                                        </div>
                                                    </t-upload>
                                                    <div class="image-slot-remove" v-if="lastRefImage"
                                                        @click.stop="removeSelectedImage(lastRefImage.id)">
                                                        <t-icon name="close" size="14px" />
                                                    </div>
                                                </div>
                                            </div>
                                        </div>

                                        <div v-else-if="referenceMode === 'multiple'" style="text-align: center">
                                            <div class="reference-mode-title">多图参考 ({{ selectedImagesForVideo.length
                                            }}/{{
                                                    currentModelCapability?.maxImages || 6 }})</div>
                                            <div
                                                style="display: flex; gap: 12px; justify-content: center; flex-wrap: wrap;">
                                                <div v-for="index in (currentModelCapability?.maxImages || 6)"
                                                    :key="index" style="position: relative;">
                                                    <t-upload theme="custom" :action="uploadConfig.action"
                                                        :headers="uploadConfig.headers" :show-file-list="false"
                                                        accept="image/*" :before-upload="beforeUpload"
                                                        @success="(ctx) => handleUploadRefSuccess(ctx, `multi_${index}`)">
                                                        <div class="image-slot image-slot-small"
                                                            :class="{ selected: !!getMultiRefImage(index) }">
                                                            <t-image v-if="getMultiRefImage(index)"
                                                                :src="getImageUrl(getMultiRefImage(index).url || getMultiRefImage(index).imageUrl)"
                                                                fit="cover" class="img" />
                                                            <div v-else class="image-slot-placeholder">
                                                                <t-icon name="add" size="16px" />
                                                                <div style="margin-top: 4px; font-size: 10px;">图 {{
                                                                    index }}
                                                                </div>
                                                            </div>
                                                        </div>
                                                    </t-upload>
                                                    <div class="image-slot-remove" v-if="getMultiRefImage(index)"
                                                        @click.stop="removeSelectedImage(getMultiRefImage(index).id)">
                                                        <t-icon name="close" size="12px" />
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <div class="reference-images-section">
                                        <div class="frame-type-buttons">
                                            <t-radio-group v-model="selectedVideoFrameType" variant="default-filled">
                                                <t-radio-button value="first">首帧</t-radio-button>
                                                <t-radio-button value="last">尾帧</t-radio-button>
                                                <t-radio-button value="action">动作序列</t-radio-button>
                                                <t-radio-button value="key">关键帧</t-radio-button>
                                            </t-radio-group>
                                        </div>

                                        <div class="frame-type-content">
                                            <div v-if="referenceMode === 'first_last' && selectedVideoFrameType === 'first' && previousStoryboardLastFrames.length > 0"
                                                class="previous-frame-section">
                                                <div
                                                    style="display: flex; align-items: center; gap: 6px; margin-bottom: 6px;">
                                                    <t-tag size="small" theme="primary" variant="light">上一镜头 #{{
                                                        previousStoryboard?.sequenceNo }} 尾帧</t-tag>
                                                    <span class="hint-text">点击添加为首帧参考</span>
                                                </div>
                                                <div class="reference-grid">
                                                    <div v-for="img in previousStoryboardLastFrames"
                                                        :key="'prev-' + img.id" class="reference-item"
                                                        :class="{ selected: isPreviousFrameSelected(img) }"
                                                        @click="handlePreviousImageSelect(img)">
                                                        <t-image :src="getImageUrl(img.url || img.imageUrl)" fit="cover"
                                                            class="img" />
                                                        <div v-if="isPreviousFrameSelected(img)" class="check-mark">✓
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>

                                            <div class="image-scroll-container">
                                                <div class="reference-grid"
                                                    v-if="currentFrameImagesFiltered(selectedVideoFrameType).length > 0">
                                                    <div v-for="img in currentFrameImagesFiltered(selectedVideoFrameType)"
                                                        :key="img.id" class="reference-item"
                                                        :class="{ selected: isImageSelected(img) }"
                                                        @click="handleImageSelect(img)">
                                                        <t-image :src="getImageUrl(img.url || img.imageUrl)" fit="cover"
                                                            class="img" />
                                                        <div v-if="isImageSelected(img)" class="check-mark">✓</div>
                                                        <t-image-viewer :close-on-overlay="true"
                                                            :images="[getImageUrl(img.url || img.imageUrl)]">
                                                            <template #trigger="{ open }">
                                                                <div class="preview-icon" @click.stop="open">
                                                                    <t-icon name="zoom-in" size="14px" />
                                                                </div>
                                                            </template>
                                                        </t-image-viewer>
                                                    </div>
                                                </div>
                                                <div v-else class="empty-text" style="padding: 20px 0;">
                                                    <t-empty
                                                        :description="`暂无${getFrameTypeName(selectedVideoFrameType)}图片，请在上方上传或在[镜头图片]面板生成`" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                <t-button theme="primary" block size="large" :loading="generatingVideo"
                                    @click="generateVideo" style="margin-top: 24px;">
                                    <template #icon><t-icon name="video" /></template> 生成视频
                                </t-button>
                            </div>

                            <div class="video-list-area section-group" v-if="generatedVideos.length > 0">
                                <div class="section-header"><span>生成结果 ({{ generatedVideos.length }})</span></div>
                                <div class="video-card-list">
                                    <div v-for="video in generatedVideos" :key="video.id" class="video-card">
                                        <video :src="getVideoUrl(video.url)" controls></video>
                                        <div class="video-actions">
                                            <t-tag theme="success" variant="light" size="small">已完成</t-tag>
                                            <div class="action-btns">
                                                <t-tooltip content="更新到时间线"><t-button size="small" variant="text"
                                                        @click="addVideoToAssets(video)"><t-icon
                                                            name="layers" /></t-button></t-tooltip>
                                                <t-button size="small" variant="text" theme="danger"
                                                    @click="deleteVideo(video)"><t-icon name="delete" /></t-button>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <t-empty v-else description="请选择一个镜头" style="margin-top: 40px" />
                    </t-tab-panel>

                    <t-tab-panel value="audio" label="音效配乐">
                        <div class="tab-content">
                            <t-empty description="音效与配乐生成功能开发中... 请在'镜头属性'中配置描述" />
                        </div>
                    </t-tab-panel>

                    <t-tab-panel value="merge" label="视频合成">
                        <div class="tab-content scrollable-content">
                            <div class="section-group">
                                <div class="section-header"><span>合成记录</span></div>
                                <div class="merge-list" v-if="videoMerges.length > 0">
                                    <div v-for="merge in videoMerges" :key="merge.id" class="merge-item">
                                        <div class="merge-info">
                                            <div class="title">{{ merge.title || '合成视频' }}</div>
                                            <div class="time">{{ merge.createTime }}</div>
                                        </div>
                                        <t-tag :theme="merge.status === 'completed' ? 'success' : 'warning'">{{
                                            merge.status ===
                                                'completed' ? '已完成' : '处理中' }}</t-tag>
                                        <t-button v-if="merge.url" size="small" variant="text"
                                            @click="previewImage(merge.url)">预览</t-button>
                                    </div>
                                </div>
                                <t-empty v-else description="暂无合成记录" size="small" />
                            </div>
                            <t-button theme="primary" block @click="exportVideo" size="large">
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
                            <t-image-viewer v-if="scene.visualPrompt" :close-on-overlay="true"
                                :images="[getImageUrl(scene.visualPrompt)]">
                                <template #trigger="{ open }">
                                    <t-image :src="getImageUrl(scene.visualPrompt)" fit="cover"
                                        style="width: 50px; height: 50px; border-radius: 4px; cursor: zoom-in;"
                                        @click.stop="open" lazy error="加载失败" />
                                </template>
                            </t-image-viewer>
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

        <t-dialog v-model:visible="showCharacterSelector" header="选择角色" width="500px"
            @confirm="showCharacterSelector = false">
            <div class="char-selector-grid">
                <div v-for="char in availableCharacters" :key="char.id" class="char-item"
                    :class="{ selected: selectedCharacters.includes(char.id) }" @click="toggleCharacterInShot(char.id)">
                    <t-image-viewer v-if="char.avatarUrl" :close-on-overlay="true"
                        :images="[getImageUrl(char.avatarUrl)]">
                        <template #trigger="{ open }">
                            <t-avatar :image="getImageUrl(char.avatarUrl)" size="large" @click.stop="open"
                                style="cursor: zoom-in;" />
                        </template>
                    </t-image-viewer>
                    <t-avatar v-else size="large">{{ char.name ? char.name[0] : '?' }}</t-avatar>
                    <span>{{ char.name }}</span>
                    <div class="check" v-if="selectedCharacters.includes(char.id)"><t-icon name="check" /></div>
                </div>
            </div>
            <t-empty v-if="availableCharacters.length === 0" description="暂无角色" />
        </t-dialog>

        <t-dialog v-model:visible="showPropSelector" header="选择道具" width="500px" @confirm="showPropSelector = false">
            <div class="char-selector-grid">
                <div v-for="prop in availableProps" :key="prop.id" class="char-item"
                    :class="{ selected: selectedProps.includes(prop.id) }" @click="togglePropInShot(prop.id)">
                    <t-image-viewer v-if="prop.imageUrl" :close-on-overlay="true"
                        :images="[getImageUrl(prop.imageUrl)]">
                        <template #trigger="{ open }">
                            <t-image :src="getImageUrl(prop.imageUrl)" fit="contain"
                                style="width: 50px; height: 50px; border-radius: 4px; background: #f9f9f9; cursor: zoom-in;"
                                @click.stop="open" lazy error="加载失败" />
                        </template>
                    </t-image-viewer>
                    <t-icon v-else name="image" size="24px" style="color: #ccc" />
                    <span>{{ prop.name }}</span>
                    <div class="check" v-if="selectedProps.includes(prop.id)"><t-icon name="check" /></div>
                </div>
            </div>
            <t-empty v-if="availableProps.length === 0" description="暂无道具" />
        </t-dialog>

        <GridImageEditor v-model="showGridEditor" :storyboard-id="currentStoryboardId" :drama-id="dramaId"
            :all-images="currentFrameImages" @success="handleGridSuccess" />

        <ImageCropDialog v-model="showCropDialog" :image-url="cropImageUrl" @save="handleCropSave" />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import {
    ArrowLeftIcon, RefreshIcon, AddIcon, DeleteIcon, MagicIcon,
    UploadIcon, ZoomInIcon, VideoIcon, LinkIcon, LayersIcon,
    MoveIcon, AddCircleIcon, FilmIcon, CheckIcon, DownloadIcon, CloseIcon, CutIcon, CloseCircleFilledIcon
} from 'tdesign-icons-vue-next'

// API
import { findProjects } from '@/api/projects'
import { getScriptsList } from '@/api/scripts'
import { getScenesList } from '@/api/scenes'
import { getCharactersList } from '@/api/characters'
import { getPropsList } from '@/api/props'
import { getShotsList, createShots, updateShots, deleteShots } from '@/api/shots'
import { getAssetsList, createAsset, deleteAsset } from '@/api/assets'
import { extractFramePromptTask, findTasks, generateImageByPromptTask } from '@/api/tasks'
import { createShotFrameImages, deleteShotFrameImages } from '@/api/shot_frame_image'
import { getImageUrl } from '@/utils/format'

// 组件
import VideoTimelineEditor from '@/components/editor/VideoTimelineEditor.vue'
import GridImageEditor from '@/components/editor/GridImageEditor.vue'
import ImageCropDialog from '@/components/editor/ImageCropDialog.vue'

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

// === 右侧面板状态 ===
const activeTab = ref('shot')
const showSceneSelector = ref(false)
const showCharacterSelector = ref(false)
const showPropSelector = ref(false)

// 图片生成状态
const selectedFrameType = ref('first')
const generatingImage = ref(false)
const extractingPrompt = ref(false)
const showGridEditor = ref(false)
const showCropDialog = ref(false)
const cropImageUrl = ref('')
const uploadingImage = ref(false)

// 合成状态
const videoMerges = ref<any[]>([])

const currentPreviewUrl = ref('')
const timelineEditorRef = ref<any>(null)
const mainPlayerRef = ref<HTMLVideoElement | null>(null)

// ================= 🔴 视频生成相关逻辑重构 =================
const generatingVideo = ref(false)
const selectedVideoModel = ref('kling')
const videoDuration = ref(5)
const referenceMode = ref('single')
const selectedVideoFrameType = ref('first') // 视频参考图界面的子TAB
const generatedVideos = ref<any[]>([])

// 存放选中的图片ID (如果是单图或首尾帧的首帧，放[0]；多图放多个)
const selectedImagesForVideo = ref<number[]>([])
// 特供首尾帧的尾帧使用
const selectedLastImageForVideo = ref<number | null>(null)

interface VideoModelCapability {
    id: string;
    name: string;
    supportMultipleImages: boolean;
    supportFirstLastFrame: boolean;
    supportSingleImage: boolean;
    supportTextOnly: boolean;
    maxImages: number;
}

// 模拟配置 (实际可从后端接口拉取)
const videoModelCapabilities = ref<VideoModelCapability[]>([
    { id: 'kling', name: '可灵 (Kling)', supportSingleImage: true, supportMultipleImages: false, supportFirstLastFrame: false, supportTextOnly: true, maxImages: 1 },
    { id: 'runway', name: 'Runway Gen-3', supportSingleImage: true, supportMultipleImages: false, supportFirstLastFrame: true, supportTextOnly: true, maxImages: 2 },
    { id: 'luma', name: 'Luma DreamMachine', supportSingleImage: true, supportMultipleImages: false, supportFirstLastFrame: true, supportTextOnly: true, maxImages: 2 },
    { id: 'pika', name: 'Pika 1.0', supportSingleImage: true, supportMultipleImages: true, supportFirstLastFrame: false, supportTextOnly: true, maxImages: 4 }
]);

const currentModelCapability = computed(() => {
    return videoModelCapabilities.value.find(m => m.id === selectedVideoModel.value);
});

const availableReferenceModes = computed(() => {
    const capability = currentModelCapability.value;
    if (!capability) return [];
    const modes: Array<{ value: string; label: string; description?: string }> = [];
    if (capability.supportTextOnly) modes.push({ value: "none", label: "纯文本", description: "不使用参考图" });
    if (capability.supportSingleImage) modes.push({ value: "single", label: "单图参考", description: "使用单张参考图" });
    if (capability.supportFirstLastFrame) modes.push({ value: "first_last", label: "首尾帧", description: "使用首帧和尾帧" });
    if (capability.supportMultipleImages) modes.push({ value: "multiple", label: "多图", description: `最多${capability.maxImages}张` });
    return modes;
});

// 监听视频模型切换，清空已选图片和参考图模式
watch(selectedVideoModel, () => {
    selectedImagesForVideo.value = [];
    selectedLastImageForVideo.value = null;
    referenceMode.value = availableReferenceModes.value[0]?.value || 'none';
});

// 监听参考图模式切换，清空已选图片
watch(referenceMode, () => {
    selectedImagesForVideo.value = [];
    selectedLastImageForVideo.value = null;
});

// === Computed ===
const currentStoryboard = computed(() => storyboards.value.find(s => String(s.id) === String(currentStoryboardId.value)))

// 获取上一个镜头信息
const previousStoryboard = computed(() => {
    if (!currentStoryboardId.value || storyboards.value.length < 2) return null;
    const currentIndex = storyboards.value.findIndex((s) => String(s.id) === String(currentStoryboardId.value));
    if (currentIndex <= 0) return null;
    return storyboards.value[currentIndex - 1];
});

// 上一个镜头的尾帧图片 (必须是从生成记录里过滤)
const previousStoryboardLastFrames = computed(() => {
    if (!previousStoryboard.value || !previousStoryboard.value.frameImages) return [];
    return previousStoryboard.value.frameImages.filter((img: any) => img.frameType === 'last' && (!img.imageType || img.imageType === 'shot'));
});

// 🔴 供图库过滤展示当前镜头的对应类型的图片
const currentFrameImagesFiltered = (type: string) => {
    if (!currentStoryboard.value || !currentStoryboard.value.frameImages) return [];
    return currentStoryboard.value.frameImages.filter((img: any) =>
        img.frameType === type &&
        (!img.imageType || img.imageType === 'shot' || img.imageType === 'reference')
    );
};

// 获取所有的可用图片(包含当前分镜和上个分镜尾帧)
const getAllAvailableImages = () => {
    const curr = currentStoryboard.value?.frameImages || [];
    const prev = previousStoryboardLastFrames.value || [];
    return [...curr, ...prev];
};

// 🔴 判断图片是否被选中为参考图
const isImageSelected = (img: any) => {
    if (referenceMode.value === 'first_last' && selectedVideoFrameType.value === 'last') {
        return selectedLastImageForVideo.value === img.id;
    }
    return selectedImagesForVideo.value.includes(img.id);
}

// 🔴 判断上一镜头的尾帧是否被选为当前首帧
const isPreviousFrameSelected = (prevImg: any) => {
    if (referenceMode.value === 'single' || referenceMode.value === 'first_last') {
        return selectedImagesForVideo.value[0] === prevImg.id;
    } else if (referenceMode.value === 'multiple') {
        return selectedImagesForVideo.value.includes(prevImg.id);
    }
    return false;
}

// 🔴 处理上一镜头尾帧（当做首帧）的点击逻辑
const handlePreviousImageSelect = (img: any) => {
    if (!referenceMode.value || referenceMode.value === 'none') {
        MessagePlugin.warning("请先选择参考图模式");
        return;
    }

    if (referenceMode.value === 'single' || referenceMode.value === 'first_last') {
        if (selectedImagesForVideo.value[0] === img.id) {
            selectedImagesForVideo.value = [];
        } else {
            selectedImagesForVideo.value = [img.id];
        }
    } else if (referenceMode.value === 'multiple') {
        const index = selectedImagesForVideo.value.indexOf(img.id);
        if (index > -1) {
            selectedImagesForVideo.value.splice(index, 1);
        } else {
            if (selectedImagesForVideo.value.length >= (currentModelCapability.value?.maxImages || 6)) {
                MessagePlugin.warning(`最多只能选择${currentModelCapability.value?.maxImages}张图片`);
                return;
            }
            selectedImagesForVideo.value.push(img.id);
        }
    }
};

// 🔴 处理当前镜头常规图片的点击逻辑
const handleImageSelect = (img: any) => {
    if (!referenceMode.value || referenceMode.value === 'none') {
        MessagePlugin.warning("请先选择参考图模式");
        return;
    }

    const imageId = img.id;
    const isClickingLastFrame = selectedVideoFrameType.value === 'last';

    if (referenceMode.value === 'multiple') {
        const index = selectedImagesForVideo.value.indexOf(imageId);
        if (index > -1) {
            selectedImagesForVideo.value.splice(index, 1);
        } else {
            if (selectedImagesForVideo.value.length >= (currentModelCapability.value?.maxImages || 6)) {
                MessagePlugin.warning(`最多只能选择${currentModelCapability.value?.maxImages}张图片`);
                return;
            }
            selectedImagesForVideo.value.push(imageId);
        }
        return;
    }

    if (referenceMode.value === 'single') {
        if (selectedImagesForVideo.value[0] === imageId) {
            selectedImagesForVideo.value = [];
        } else {
            selectedImagesForVideo.value = [imageId];
        }
        return;
    }

    if (referenceMode.value === 'first_last') {
        if (isClickingLastFrame) {
            if (selectedLastImageForVideo.value === imageId) {
                selectedLastImageForVideo.value = null;
            } else {
                selectedLastImageForVideo.value = imageId;
            }
        } else {
            if (selectedImagesForVideo.value[0] === imageId) {
                selectedImagesForVideo.value = [];
            } else {
                selectedImagesForVideo.value = [imageId];
            }
        }
    }
};

// 🔴 辅助：获取某个类型的名称
const getFrameTypeName = (type: string) => {
    const map: Record<string, string> = { first: '首帧', last: '尾帧', action: '动作序列', key: '关键帧' };
    return map[type] || type;
}

// 🔴 从槽位上点击 X 号移除选择
const removeSelectedImage = (imageId: number) => {
    if (selectedLastImageForVideo.value === imageId) {
        selectedLastImageForVideo.value = null;
        return;
    }
    const index = selectedImagesForVideo.value.indexOf(imageId);
    if (index > -1) {
        selectedImagesForVideo.value.splice(index, 1);
    }
};

// 视频参考图相关 Computed 绑定 (找到对象以供槽位展示)
const singleRefImage = computed(() => {
    return getAllAvailableImages().find(i => i.id === selectedImagesForVideo.value[0]);
})
const firstRefImage = computed(() => {
    return getAllAvailableImages().find(i => i.id === selectedImagesForVideo.value[0]);
})
const lastRefImage = computed(() => {
    return getAllAvailableImages().find(i => i.id === selectedLastImageForVideo.value);
})
const getMultiRefImage = (index: number) => {
    return getAllAvailableImages().find(i => i.id === selectedImagesForVideo.value[index - 1]);
}


const currentScene = computed(() => {
    if (!currentStoryboard.value || !currentStoryboard.value.sceneId) return null
    return sceneList.value.find(s => s.id === currentStoryboard.value.sceneId)
})

const selectedCharacters = computed(() => {
    if (!currentStoryboard.value?.characters) return []
    return currentStoryboard.value.characters.map((c: any) => typeof c === 'object' ? c.id : c)
})

const selectedProps = computed(() => {
    if (!currentStoryboard.value?.props) return []
    return currentStoryboard.value.props.map((p: any) => typeof p === 'object' ? p.id : p)
})

// 动态计算当前选择帧类型的提示词
const currentFramePromptText = computed({
    get() {
        if (!currentStoryboard.value) return '';
        const prompts = currentStoryboard.value.framePrompts || [];
        const fp = prompts.find((p: any) => p.frameType === selectedFrameType.value);
        if (!fp && selectedFrameType.value === 'first' && currentStoryboard.value.imagePrompt) {
            return currentStoryboard.value.imagePrompt;
        }
        return fp ? fp.prompt : '';
    },
    set(val) {
        if (!currentStoryboard.value) return;
        if (!currentStoryboard.value.framePrompts) {
            currentStoryboard.value.framePrompts = [];
        }
        const fpIndex = currentStoryboard.value.framePrompts.findIndex((p: any) => p.frameType === selectedFrameType.value);
        if (fpIndex > -1) {
            currentStoryboard.value.framePrompts[fpIndex].prompt = val;
        } else {
            currentStoryboard.value.framePrompts.push({
                frameType: selectedFrameType.value,
                prompt: val
            });
        }
        if (selectedFrameType.value === 'first') {
            currentStoryboard.value.imagePrompt = val;
        }
    }
})

// 镜头图片列表：过滤 imageType 为 shot 或空(兼容老数据) 的图片
const currentFrameImages = computed(() => {
    if (!currentStoryboard.value || !currentStoryboard.value.frameImages) return [];
    return currentStoryboard.value.frameImages.filter((img: any) =>
        img.frameType === selectedFrameType.value &&
        (!img.imageType || img.imageType === 'shot')
    );
})

const getCharacterById = (id: number) => availableCharacters.value.find(c => c.id === id)
const getPropById = (id: number) => availableProps.value.find(p => p.id === id)

const getAuthToken = () => localStorage.getItem('token')

const uploadConfig = reactive({
    action: import.meta.env.VITE_API_URL + '/admin/v1/upload/singleUpload',
    headers: computed(() => ({ 'Authorization': `${getAuthToken()}` })),
    sizeLimit: 5 * 1024 * 1024,
    allowedFormats: ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
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
                selectStoryboard(storyboards.value[0].id)
            }
        }
    }
}

const loadVideoAssets = async () => { /* 忽略，已保留 */ }

// === 核心逻辑 ===
const toggleCharacterInShot = async (charId: number) => {
    if (!currentStoryboard.value) return
    let chars = currentStoryboard.value.characters || []
    const idx = chars.findIndex((c: any) => (typeof c === 'object' ? c.id === charId : c === charId))

    if (idx > -1) {
        chars.splice(idx, 1)
    } else {
        const fullChar = availableCharacters.value.find(c => c.id === charId)
        if (fullChar) chars.push(fullChar)
    }
    currentStoryboard.value.characters = chars
    currentStoryboard.value.characterIds = chars.map((c: any) => typeof c === 'object' ? c.id : c)
    await saveStoryboardField()
}

const togglePropInShot = async (propId: number) => {
    if (!currentStoryboard.value) return
    let propsArray = currentStoryboard.value.props || []
    const idx = propsArray.findIndex((p: any) => (typeof p === 'object' ? p.id === propId : p === propId))

    if (idx > -1) {
        propsArray.splice(idx, 1)
    } else {
        const fullProp = availableProps.value.find(p => p.id === propId)
        if (fullProp) propsArray.push(fullProp)
    }
    currentStoryboard.value.props = propsArray
    currentStoryboard.value.propIds = propsArray.map((p: any) => typeof p === 'object' ? p.id : p)
    await saveStoryboardField()
}

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

const handleTimelineDrop = (clipData: any) => { /* 忽略 */ }

const goBack = () => router.back()
const loadData = () => { initData(); MessagePlugin.success('数据已刷新') }

const selectStoryboard = (id: number | string) => {
    currentStoryboardId.value = id
}

const handleTimelineSelect = (clip: any) => {
    if (clip.storyboardId) selectStoryboard(clip.storyboardId)
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

const handleAddStoryboard = async () => {
    const newShot = {
        projectId: Number(dramaId),
        scriptId: currentScriptId.value,
        title: `新镜头 ${storyboards.value.length + 1}`,
        durationMs: 3000,
        shotType: '中景',
        angle: '平视',
        cameraMovement: '固定'
    }
    try {
        MessagePlugin.success('添加成功 (Mock)')
        storyboards.value.push({ id: Date.now(), ...newShot })
    } catch { MessagePlugin.error('添加失败') }
}

const handleDeleteStoryboard = async (shot: any) => { MessagePlugin.success('删除成功') }

const linkSceneToShot = async (scene: any) => {
    if (!currentStoryboard.value) return
    currentStoryboard.value.sceneId = scene.id
    await saveStoryboardField()
    showSceneSelector.value = false; MessagePlugin.success('已关联场景')
}

const saveStoryboardField = async () => {
    if (!currentStoryboard.value) return
    saving.value = true
    try {
        const payload = {
            ...currentStoryboard.value,
            characterIds: currentStoryboard.value.characters?.map((c: any) => typeof c === 'object' ? c.id : c) || [],
            propIds: currentStoryboard.value.props?.map((p: any) => typeof p === 'object' ? p.id : p) || []
        }
        delete payload.characters
        delete payload.props
        delete payload.scenes
        delete payload.frameImages
        delete payload.framePrompts

        await updateShots(payload.id, payload)
    } catch {
        MessagePlugin.error('保存失败')
    } finally {
        saving.value = false
    }
}

const updateShotDurationMs = (secVal: number) => {
    if (!currentStoryboard.value) return
    currentStoryboard.value.durationMs = secVal * 1000
    saveStoryboardField()
}

// === 图片提取与生成相关 ===
const extractFramePrompt = async () => {
    if (!currentStoryboard.value) return;
    extractingPrompt.value = true;
    try {
        const res = await extractFramePromptTask({
            shotId: currentStoryboard.value.id,
            frameType: selectedFrameType.value
        });
        const taskId = res.data?.task_id || res.data?.taskId || res.data?.data?.task_id;

        if (taskId) {
            MessagePlugin.loading('AI 正在提取提示词...');
            const timer = setInterval(async () => {
                try {
                    const taskRes = await findTasks(taskId);
                    const taskData = taskRes.data?.data || taskRes.data;
                    const status = taskData?.status;
                    if (status === 'completed' || status === 2) {
                        clearInterval(timer);
                        extractingPrompt.value = false;
                        MessagePlugin.success('提示词提取成功');
                        await loadShotsData();
                    } else if (status === 'failed' || status === 3) {
                        clearInterval(timer);
                        extractingPrompt.value = false;
                        MessagePlugin.error(taskData?.error || '提取失败');
                    }
                } catch (e) {
                    clearInterval(timer);
                    extractingPrompt.value = false;
                }
            }, 2000);
        } else {
            extractingPrompt.value = false;
            MessagePlugin.error('任务提交失败');
        }
    } catch (e) {
        extractingPrompt.value = false;
        MessagePlugin.error('请求异常');
    }
}

// 对接生成图片
const generateFrameImage = async () => {
    if (!currentStoryboard.value) return;
    if (!currentFramePromptText.value) {
        MessagePlugin.warning('请先输入或提取提示词');
        return;
    }

    generatingImage.value = true;
    try {
        const res = await generateImageByPromptTask({
            shotId: currentStoryboard.value.id,
            frameType: selectedFrameType.value,
            prompt: currentFramePromptText.value
        });

        const taskId = res.data?.task_id || res.data?.taskId || res.data?.data?.task_id;

        if (taskId) {
            MessagePlugin.loading('AI 正在生成画面，请耐心等待...');
            const timer = setInterval(async () => {
                try {
                    const taskRes = await findTasks(taskId);
                    const taskData = taskRes.data?.data || taskRes.data;
                    const status = taskData?.status;

                    if (status === 'completed' || status === 2) {
                        clearInterval(timer);
                        generatingImage.value = false;
                        MessagePlugin.success('画面生成成功！');
                        await loadShotsData();
                    } else if (status === 'failed' || status === 3) {
                        clearInterval(timer);
                        generatingImage.value = false;
                        MessagePlugin.error(taskData?.error || '生成失败');
                    }
                } catch (e) {
                    clearInterval(timer);
                    generatingImage.value = false;
                }
            }, 3000);
        } else {
            generatingImage.value = false;
            MessagePlugin.error('生图任务提交失败');
        }
    } catch (e) {
        generatingImage.value = false;
        MessagePlugin.error('生图请求异常');
    }
}

const beforeUpload = (file: any) => {
    if (!uploadConfig.allowedFormats.includes(file.type)) {
        MessagePlugin.error('不支持的文件格式')
        return false
    }
    if (file.size > uploadConfig.sizeLimit) {
        MessagePlugin.error('图片大小不能超过 5MB')
        return false
    }
    uploadingImage.value = true
    return true
}

const handleUploadFail = () => {
    uploadingImage.value = false
    MessagePlugin.error('上传失败')
}

// 镜头图片上传成功
const handleUploadImageSuccess = async (ctx: any) => {
    uploadingImage.value = false
    const response = ctx.response

    if (response?.code === 0 || response?.code === 200) {
        let fileUrl = response.data.file_url || response.data.url
        if (fileUrl && fileUrl.startsWith('/')) {
            fileUrl = import.meta.env.VITE_API_URL.replace(/\/admin\/v1$/, '').replace(/\/v1$/, '') + fileUrl
        }

        if (currentStoryboard.value) {
            try {
                const res = await createShotFrameImages({
                    projectId: Number(dramaId),
                    shotId: currentStoryboard.value.id,
                    frameType: selectedFrameType.value,
                    imageType: 'shot',
                    imageUrl: fileUrl
                });

                if (res.code === 0) {
                    MessagePlugin.success('图片添加成功');
                    if (!currentStoryboard.value.frameImages) {
                        currentStoryboard.value.frameImages = [];
                    }
                    currentStoryboard.value.frameImages.unshift(res.data);

                    if (!currentStoryboard.value.imageUrl || selectedFrameType.value === 'first') {
                        currentStoryboard.value.imageUrl = fileUrl;
                        saveStoryboardField();
                    }
                } else {
                    MessagePlugin.error(res.message || '图片数据保存失败');
                }
            } catch (err) {
                console.error(err);
                MessagePlugin.error('图片数据请求异常');
            }
        }
    } else {
        MessagePlugin.error(response?.msg || '上传失败')
    }
}

// 🔴 参考图上传成功，自动记录并填入对应槽位
const handleUploadRefSuccess = async (ctx: any, targetSlot: string) => {
    const response = ctx.response;
    if (response?.code === 0 || response?.code === 200) {
        let fileUrl = response.data.file_url || response.data.url;
        if (fileUrl && fileUrl.startsWith('/')) {
            fileUrl = import.meta.env.VITE_API_URL.replace(/\/admin\/v1$/, '').replace(/\/v1$/, '') + fileUrl;
        }
        if (currentStoryboard.value) {
            try {
                // 统一传为 reference 类型的图片
                const res = await createShotFrameImages({
                    projectId: Number(dramaId),
                    shotId: currentStoryboard.value.id,
                    frameType: targetSlot, // 作为标记
                    imageType: 'reference',
                    imageUrl: fileUrl
                });

                if (res.code === 0) {
                    MessagePlugin.success('上传并设置成功');
                    if (!currentStoryboard.value.frameImages) {
                        currentStoryboard.value.frameImages = [];
                    }
                    currentStoryboard.value.frameImages.unshift(res.data);

                    // 自动选择到对应的槽位
                    if (targetSlot === 'single' || targetSlot === 'first') {
                        selectedImagesForVideo.value = [res.data.id];
                    } else if (targetSlot === 'last') {
                        selectedLastImageForVideo.value = res.data.id;
                    } else if (targetSlot.startsWith('multi')) {
                        if (selectedImagesForVideo.value.length < (currentModelCapability.value?.maxImages || 6)) {
                            selectedImagesForVideo.value.push(res.data.id);
                        }
                    }
                }
            } catch (err) {
                MessagePlugin.error('请求异常');
            }
        }
    } else {
        MessagePlugin.error(response?.msg || '上传失败');
    }
};

// 统一删除图片 (支持镜头图和参考图)
const deleteImage = async (img: any) => {
    if (!img.id || !currentStoryboard.value) return;

    const confirmDialog = DialogPlugin.confirm({
        header: '确认删除',
        body: '确定要删除这张图片吗？',
        onConfirm: async () => {
            try {
                const res = await deleteShotFrameImages(img.id);
                if (res.code === 0 || res.code === 200) {
                    const idx = currentStoryboard.value.frameImages.findIndex((i: any) => i.id === img.id);
                    if (idx > -1) {
                        currentStoryboard.value.frameImages.splice(idx, 1);
                    }
                    // 如果删的图正好在选中列表里，清除选中
                    removeSelectedImage(img.id);
                    MessagePlugin.success('删除成功');
                } else {
                    MessagePlugin.error(res.message || '删除失败');
                }
            } catch (e) {
                MessagePlugin.error('删除请求异常');
            } finally {
                confirmDialog.destroy();
            }
        },
        onCancel: () => {
            confirmDialog.destroy();
        }
    });
}

const openCropDialog = (img: any) => {
    let fullUrl = img.url || img.imageUrl;
    if (!fullUrl.startsWith('http')) {
        fullUrl = getImageUrl(fullUrl);
    }
    cropImageUrl.value = fullUrl;
    showCropDialog.value = true;
}

const handleCropSave = (newUrl: string) => { showCropDialog.value = false }

// 宫格图生成成功回调处理
const handleGridSuccess = async (data: { url: string, frameType: string }) => {
    if (data && data.url && currentStoryboard.value) {
        try {
            const res = await createShotFrameImages({
                projectId: Number(dramaId),
                shotId: currentStoryboard.value.id,
                frameType: data.frameType,
                imageType: 'shot', // 🔴 宫格图属于镜头图
                imageUrl: data.url
            });

            if (res.code === 0) {
                if (!currentStoryboard.value.frameImages) {
                    currentStoryboard.value.frameImages = [];
                }
                currentStoryboard.value.frameImages.unshift(res.data);
                MessagePlugin.success('宫格图保存成功');
            } else {
                MessagePlugin.error(res.message || '宫格图保存记录失败');
            }
        } catch (e) {
            console.error(e);
            MessagePlugin.error('宫格图保存异常');
        }
    }
}

// === 视频生成相关 ===
const generateVideo = async () => {
    if (!selectedVideoModel.value) {
        MessagePlugin.warning("请先选择视频生成模型");
        return;
    }
    if (!currentStoryboard.value) return;

    if (referenceMode.value === 'single' && selectedImagesForVideo.value.length === 0) {
        MessagePlugin.warning("请上传或选择单图参考");
        return;
    }
    if (referenceMode.value === 'first_last' && (selectedImagesForVideo.value.length === 0 || !selectedLastImageForVideo.value)) {
        MessagePlugin.warning("请选择完整的首尾帧参考图");
        return;
    }

    generatingVideo.value = true;
    try {
        const requestPayload: any = {
            projectId: dramaId,
            shotId: currentStoryboard.value.id,
            model: selectedVideoModel.value,
            duration: videoDuration.value,
            prompt: currentStoryboard.value.videoPrompt || currentStoryboard.value.imagePrompt || currentStoryboard.value.visualDesc || '',
            referenceMode: referenceMode.value
        }

        const allImages = getAllAvailableImages();

        if (referenceMode.value === 'single') {
            const img = allImages.find(i => i.id === selectedImagesForVideo.value[0]);
            requestPayload.imageUrl = img?.url || img?.imageUrl;
        } else if (referenceMode.value === 'first_last') {
            const firstImg = allImages.find(i => i.id === selectedImagesForVideo.value[0]);
            const lastImg = allImages.find(i => i.id === selectedLastImageForVideo.value);
            requestPayload.firstFrameUrl = firstImg?.url || firstImg?.imageUrl;
            requestPayload.lastFrameUrl = lastImg?.url || lastImg?.imageUrl;
        } else if (referenceMode.value === 'multiple') {
            requestPayload.imageUrls = selectedImagesForVideo.value.map(id => {
                const img = allImages.find(i => i.id === id);
                return img?.url || img?.imageUrl;
            });
        }

        // 假设有个提交生成视频的接口
        // const res = await createVideoTask(requestPayload);
        setTimeout(() => {
            generatingVideo.value = false;
            MessagePlugin.success('视频生成任务已提交');
        }, 1500)

    } catch (e) {
        generatingVideo.value = false;
        MessagePlugin.error("任务提交失败");
    }
}

const addVideoToAssets = async (video: any) => { MessagePlugin.success('已添加到素材库') }
const deleteVideo = async (video: any) => { /* 忽略 */ }

const previewImage = (url: string) => window.open(getImageUrl(url), '_blank')
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

            .edit-tabs {
                height: 100%;
                display: flex;
                flex-direction: column;

                :deep(.t-tabs__nav) {
                    flex-shrink: 0;
                }

                :deep(.t-tabs__content) {
                    flex: 1;
                    overflow: hidden;
                    display: flex;
                    flex-direction: column;
                }

                :deep(.t-tab-panel) {
                    flex: 1;
                    overflow: hidden;
                    display: flex;
                    flex-direction: column;
                }
            }

            .tab-content {
                padding: 16px;
                flex: 1;
                overflow-y: auto;
                padding-bottom: 60px;

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

            /* 图片生成区域 */
            .grid-entry-card {
                margin-bottom: 12px;
                height: 50px;
                border: 1px dashed var(--td-brand-color);
                border-radius: 4px;
                display: flex;
                align-items: center;
                justify-content: center;
                gap: 8px;
                cursor: pointer;
                color: var(--td-brand-color);
                font-size: 13px;

                &:hover {
                    background: var(--td-brand-color-light);
                }
            }

            .image-grid-list {
                display: grid;
                grid-template-columns: repeat(2, 1fr);
                gap: 10px;

                .image-grid-item {
                    position: relative;
                    height: 100px;
                    border-radius: 4px;
                    overflow: hidden;

                    .img {
                        width: 100%;
                        height: 100%;
                    }

                    .img-overlay {
                        position: absolute;
                        inset: 0;
                        background: rgba(0, 0, 0, 0.5);
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        opacity: 0;
                        transition: opacity 0.2s;
                        z-index: 5;

                        .actions-wrapper {
                            display: flex;
                            gap: 16px;
                        }

                        .icon-btn {
                            width: 32px;
                            height: 32px;
                            border-radius: 50%;
                            background: rgba(255, 255, 255, 0.2);
                            color: #fff;
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            cursor: pointer;
                            transition: all 0.2s;

                            &:hover {
                                background: rgba(255, 255, 255, 0.4);
                                transform: scale(1.1);
                            }

                            &.danger:hover {
                                background: var(--td-error-color);
                            }
                        }
                    }

                    &:hover .img-overlay {
                        opacity: 1;
                    }

                    /* 动作序列特有的裁剪按钮 */
                    .crop-btn {
                        position: absolute;
                        top: 4px;
                        right: 4px;
                        background: rgba(0, 0, 0, 0.6);
                        color: #fff;
                        width: 24px;
                        height: 24px;
                        border-radius: 4px;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        cursor: pointer;
                        z-index: 10;
                        transition: all 0.2s;
                        display: none;

                        &:hover {
                            background: var(--td-brand-color);
                        }
                    }

                    &:hover .crop-btn {
                        display: flex;
                    }
                }
            }

            /* 🔴 参考图选择器样式 */
            .reference-images-section {
                margin-top: 16px;

                .frame-type-buttons {
                    margin-bottom: 12px;
                    text-align: center;
                }

                .frame-type-content {
                    background: var(--td-bg-color-secondarycontainer);
                    padding: 12px;
                    border-radius: 6px;
                }

                .previous-frame-section {
                    margin-bottom: 12px;
                }

                .reference-grid {
                    display: grid;
                    grid-template-columns: repeat(3, 1fr);
                    gap: 8px;

                    .reference-item-mini {
                        position: relative;
                        height: 70px;
                        border-radius: 4px;
                        overflow: hidden;
                        cursor: pointer;
                        border: 2px solid transparent;
                        transition: all 0.2s;

                        &.selected {
                            border-color: var(--td-brand-color);
                            box-shadow: 0 0 0 1px var(--td-brand-color);
                        }

                        &:hover {
                            border-color: var(--td-brand-color);
                        }

                        .img {
                            width: 100%;
                            height: 100%;
                        }

                        .check-mark {
                            position: absolute;
                            top: 0;
                            right: 0;
                            background: var(--td-brand-color);
                            color: #fff;
                            border-bottom-left-radius: 4px;
                            padding: 2px 4px;
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            font-size: 12px;
                        }

                        .preview-icon {
                            position: absolute;
                            bottom: 4px;
                            right: 4px;
                            background: rgba(0, 0, 0, 0.5);
                            color: #fff;
                            border-radius: 4px;
                            padding: 4px;
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            opacity: 0;
                            transition: opacity 0.2s;
                        }

                        &:hover .preview-icon {
                            opacity: 1;
                        }
                    }
                }

                .ref-container {
                    display: flex;
                    gap: 10px;

                    &.center {
                        justify-content: center;
                    }

                    &.row {
                        justify-content: space-between;
                        align-items: center;
                    }

                    .ref-container-column {
                        display: flex;
                        flex-direction: column;
                        gap: 10px;
                    }

                    .slot-wrapper {
                        display: flex;
                        flex-direction: column;
                        align-items: center;
                        gap: 4px;
                    }

                    /* 🔴 新增：用于包裹上传插槽和删除按钮 */
                    .ref-image-wrapper {
                        position: relative;
                        display: inline-block;

                        .ref-delete-btn {
                            position: absolute;
                            top: -6px;
                            right: -6px;
                            color: var(--td-error-color);
                            background: #fff;
                            border-radius: 50%;
                            cursor: pointer;
                            z-index: 10;
                            display: flex;
                            transition: transform 0.2s;
                            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

                            &:hover {
                                transform: scale(1.1);
                            }
                        }
                    }

                    .ref-image-slot {
                        width: 120px;
                        height: 70px;
                        border: 1px dashed var(--td-component-stroke);
                        border-radius: 4px;
                        overflow: hidden;
                        cursor: pointer;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        background: var(--td-bg-color-container);

                        &.slot-small {
                            width: 60px;
                            height: 60px;
                        }

                        &.selected {
                            border-color: var(--td-brand-color);
                            border-style: solid;
                        }

                        .placeholder {
                            color: #ccc;
                            transition: color 0.2s;
                            display: flex;
                            flex-direction: column;
                            align-items: center;
                        }

                        &:hover .placeholder {
                            color: var(--td-brand-color);
                        }
                    }
                }
            }

            /* 视频列表 */
            .video-card-list {
                display: flex;
                flex-direction: column;
                gap: 12px;

                .video-card {
                    background: #000;
                    border-radius: 4px;
                    overflow: hidden;
                    border: 1px solid var(--td-component-stroke);

                    video {
                        width: 100%;
                        max-height: 150px;
                    }

                    .video-actions {
                        padding: 8px;
                        display: flex;
                        justify-content: space-between;
                        align-items: center;
                        background: #fff;
                        border-top: 1px solid var(--td-component-stroke);

                        .action-btns {
                            display: flex;
                            gap: 4px;
                        }
                    }
                }
            }

            /* 合成记录 */
            .merge-list {
                display: flex;
                flex-direction: column;
                gap: 8px;

                .merge-item {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding: 10px;
                    background: var(--td-bg-color-container);
                    border-radius: 4px;

                    .merge-info {
                        .title {
                            font-size: 13px;
                            font-weight: 500;
                        }

                        .time {
                            font-size: 11px;
                            color: #999;
                        }
                    }
                }
            }
        }
    }
}

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