<template>
    <t-dialog v-model:visible="visible" :header="$t('宫格图片编辑器')" width="900px"
        :confirm-btn="{ content: creating ? '生成中...' : '生成宫格图', loading: creating, disabled: !isGridComplete }"
        @confirm="createGridImage" @close="handleClose">
        <div class="grid-editor-container">
            <div class="grid-type-selector">
                <div class="label">选择宫格类型</div>
                <t-radio-group v-model="gridType" variant="default-filled" @change="initGridImages">
                    <t-radio-button :value="4">四宫格 (2x2)</t-radio-button>
                    <t-radio-button :value="6">六宫格 (3x2)</t-radio-button>
                    <t-radio-button :value="9">九宫格 (3x3)</t-radio-button>
                </t-radio-group>
            </div>

            <div class="grid-preview-area">
                <div class="grid-container" :class="`grid-${gridType}`">
                    <div v-for="(item, index) in gridImages" :key="index" class="grid-cell"
                        @click="handleGridCellClick(index)">
                        <img v-if="item.url" :src="item.url" class="cell-img" />
                        <div v-else class="cell-placeholder">
                            <t-icon name="add" size="32px" />
                            <span>点击添加</span>
                        </div>

                        <div v-if="item.url" class="cell-actions">
                            <t-button shape="circle" size="small" variant="text" @click.stop="removeGridCell(index)">
                                <t-icon name="delete" style="color: #fff" />
                            </t-button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <t-dialog v-model:visible="showImageSelector" header="选择图片" width="800px" :footer="false">
            <t-tabs defaultValue="existing">
                <t-tab-panel value="existing" label="已有图片">
                    <div class="image-grid-selector">
                        <div v-for="img in allImages" :key="img.id" class="selector-item"
                            @click="selectImageForGrid(img)">
                            <t-image :src="getImageUrl(img)" fit="cover" style="width: 100%; height: 120px;" />
                        </div>
                        <t-empty v-if="allImages.length === 0" description="暂无可用图片" />
                    </div>
                </t-tab-panel>
                <t-tab-panel value="upload" label="上传新图">
                    <t-upload theme="custom" accept="image/*" :auto-upload="false" @change="handleUploadForGrid">
                        <div class="upload-area-custom">
                            <t-icon name="cloud-upload" size="40px" />
                            <p>点击或拖拽上传</p>
                        </div>
                    </t-upload>
                </t-tab-panel>
            </t-tabs>
        </t-dialog>

    </t-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { getImageUrl } from "@/utils/format";
// 假设您有 imageAPI
// import { imageAPI } from "@/api/image"; 

const props = defineProps<{
    modelValue: boolean;
    storyboardId: number;
    dramaId: number;
    allImages: any[]; // ImageGeneration[]
}>();

const emit = defineEmits(["update:modelValue", "success"]);

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit("update:modelValue", val),
});

const gridType = ref<4 | 6 | 9>(4);
const gridImages = ref<any[]>([]);
const showImageSelector = ref(false);
const currentGridIndex = ref(-1);
const creating = ref(false);

const initGridImages = () => {
    gridImages.value = Array.from({ length: gridType.value }, () => ({}));
};

const isGridComplete = computed(() => {
    return gridImages.value.every((item) => item.url);
});

const handleGridCellClick = (index: number) => {
    currentGridIndex.value = index;
    showImageSelector.value = true;
};

const selectImageForGrid = (img: any) => {
    gridImages.value[currentGridIndex.value] = {
        id: img.id,
        url: getImageUrl(img), // 确保是完整URL
        source: 'existing'
    }
    showImageSelector.value = false;
}

const handleUploadForGrid = (file: any) => {
    // 读取本地文件预览
    const reader = new FileReader();
    reader.onload = (e) => {
        gridImages.value[currentGridIndex.value] = {
            url: e.target?.result,
            file: file.raw, // 保存原始文件用于上传
            source: 'upload'
        }
        showImageSelector.value = false;
    }
    reader.readAsDataURL(file.raw);
}

const removeGridCell = (index: number) => {
    gridImages.value[index] = {}
}

const handleClose = () => {
    visible.value = false
}

// 核心逻辑：使用 Canvas 合成图片
const createGridImage = async () => {
    if (!isGridComplete.value) return;
    creating.value = true;

    try {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        if (!ctx) throw new Error('Canvas context failed');

        const cellSize = 512;
        const gap = 10;

        // 计算行列
        let cols = 2, rows = 2;
        if (gridType.value === 6) { cols = 3; rows = 2; }
        if (gridType.value === 9) { cols = 3; rows = 3; }

        canvas.width = cellSize * cols + gap * (cols - 1);
        canvas.height = cellSize * rows + gap * (rows - 1);

        ctx.fillStyle = '#ffffff';
        ctx.fillRect(0, 0, canvas.width, canvas.height);

        // 绘制图片
        const loadImage = (url: string): Promise<HTMLImageElement> => {
            return new Promise((resolve, reject) => {
                const img = new Image();
                img.crossOrigin = "Anonymous";
                img.onload = () => resolve(img);
                img.onerror = reject;
                img.src = url;
            })
        }

        for (let i = 0; i < gridImages.value.length; i++) {
            const item = gridImages.value[i];
            const img = await loadImage(item.url);

            const col = i % cols;
            const row = Math.floor(i / cols);
            const x = col * (cellSize + gap);
            const y = row * (cellSize + gap);

            ctx.drawImage(img, x, y, cellSize, cellSize);
        }

        // 导出 Blob
        canvas.toBlob(async (blob) => {
            if (!blob) return;

            // TODO: 调用您的上传接口
            // const formData = new FormData();
            // formData.append('file', blob, 'grid.jpg');
            // const res = await uploadAPI(formData);

            // TODO: 调用保存 ImageGeneration 接口

            MessagePlugin.success('宫格图生成成功');
            emit('success');
            handleClose();

        }, 'image/jpeg', 0.9);

    } catch (e) {
        console.error(e);
        MessagePlugin.error('生成失败');
    } finally {
        creating.value = false;
    }
}

watch(visible, (val) => {
    if (val && gridImages.value.length !== gridType.value) {
        initGridImages();
    }
})

</script>

<style scoped lang="less">
.grid-editor-container {
    padding: 10px;

    .grid-type-selector {
        margin-bottom: 20px;

        .label {
            margin-bottom: 8px;
            font-weight: 600;
            color: var(--td-text-color-primary);
        }
    }

    .grid-preview-area {
        display: flex;
        justify-content: center;
        background: var(--td-bg-color-secondarycontainer);
        padding: 20px;
        border-radius: 8px;

        .grid-container {
            display: grid;
            gap: 10px;

            &.grid-4 {
                grid-template-columns: repeat(2, 200px);
                grid-template-rows: repeat(2, 200px);
            }

            &.grid-6 {
                grid-template-columns: repeat(3, 200px);
                grid-template-rows: repeat(2, 200px);
            }

            &.grid-9 {
                grid-template-columns: repeat(3, 150px);
                grid-template-rows: repeat(3, 150px);
            }

            .grid-cell {
                background: #fff;
                border: 2px dashed var(--td-component-stroke);
                border-radius: 4px;
                display: flex;
                align-items: center;
                justify-content: center;
                cursor: pointer;
                overflow: hidden;
                position: relative;

                &:hover {
                    border-color: var(--td-brand-color);

                    .cell-actions {
                        opacity: 1;
                    }
                }

                .cell-img {
                    width: 100%;
                    height: 100%;
                    object-fit: cover;
                }

                .cell-placeholder {
                    display: flex;
                    flex-direction: column;
                    align-items: center;
                    color: var(--td-text-color-placeholder);
                }

                .cell-actions {
                    position: absolute;
                    top: 4px;
                    right: 4px;
                    opacity: 0;
                    transition: opacity 0.2s;
                    background: rgba(0, 0, 0, 0.5);
                    border-radius: 50%;
                }
            }
        }
    }
}

.image-grid-selector {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 10px;
    max-height: 400px;
    overflow-y: auto;

    .selector-item {
        cursor: pointer;
        border: 2px solid transparent;

        &:hover {
            border-color: var(--td-brand-color);
        }
    }
}

.upload-area-custom {
    height: 200px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: 1px dashed var(--td-component-stroke);
    border-radius: 8px;
    cursor: pointer;

    &:hover {
        border-color: var(--td-brand-color);
        color: var(--td-brand-color);
    }
}
</style>