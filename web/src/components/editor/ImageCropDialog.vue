<template>
    <t-dialog v-model:visible="visible" header="图片裁剪" width="900px"
        :confirm-btn="{ content: '确认裁剪并保存', theme: 'primary', loading: saving }" :close-on-overlay-click="false"
        @confirm="handleSave" @close="handleClose" @opened="initCropper">
        <div class="cropper-wrapper">
            <div class="crop-container">
                <img ref="imageRef" :src="imageUrl" alt="Source Image" v-if="imageUrl" style="max-width: 100%;" />
                <div v-else class="empty-state">无图片数据</div>
            </div>

            <div class="toolbar">
                <div class="tool-group">
                    <t-tooltip content="放大">
                        <t-button variant="outline" shape="square" @click="handleZoom(0.1)">
                            <template #icon><t-icon name="zoom-in" /></template>
                        </t-button>
                    </t-tooltip>
                    <t-tooltip content="缩小">
                        <t-button variant="outline" shape="square" @click="handleZoom(-0.1)">
                            <template #icon><t-icon name="zoom-out" /></template>
                        </t-button>
                    </t-tooltip>
                </div>

                <div class="tool-group">
                    <t-tooltip content="向左旋转90°">
                        <t-button variant="outline" shape="square" @click="handleRotate(-90)">
                            <template #icon><t-icon name="history" /></template>
                        </t-button>
                    </t-tooltip>
                    <t-tooltip content="向右旋转90°">
                        <t-button variant="outline" shape="square" @click="handleRotate(90)">
                            <template #icon><t-icon name="refresh" /></template>
                        </t-button>
                    </t-tooltip>
                </div>

                <div class="tool-group">
                    <t-tooltip content="水平翻转">
                        <t-button variant="outline" shape="square" @click="handleScaleX">
                            <template #icon><t-icon name="swap" /></template>
                        </t-button>
                    </t-tooltip>
                    <t-tooltip content="垂直翻转">
                        <t-button variant="outline" shape="square" @click="handleScaleY">
                            <template #icon><t-icon name="swap" style="transform: rotate(90deg)" /></template>
                        </t-button>
                    </t-tooltip>
                </div>

                <div class="tool-group ratio-group">
                    <t-radio-group variant="default-filled" v-model="aspectRatio" @change="updateRatio">
                        <t-radio-button :value="NaN">自由</t-radio-button>
                        <t-radio-button :value="16 / 9">16:9</t-radio-button>
                        <t-radio-button :value="4 / 3">4:3</t-radio-button>
                        <t-radio-button :value="1">1:1</t-radio-button>
                        <t-radio-button :value="9 / 16">9:16</t-radio-button>
                    </t-radio-group>
                </div>

                <div class="tool-group">
                    <t-button theme="default" @click="handleReset">重置</t-button>
                </div>
            </div>
        </div>
    </t-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onUnmounted } from 'vue';
import Cropper from 'cropperjs'
import '@/assets/cropper.css'
import {
    ZoomInIcon, ZoomOutIcon, HistoryIcon, RefreshIcon, SwapIcon
} from 'tdesign-icons-vue-next';

const props = defineProps<{
    modelValue: boolean;
    imageUrl: string;
}>();

// 定义 save 事件，返回的数据结构需匹配父组件 handleCropSave
const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'save', images: { blob: Blob; frameType: string }[]): void;
}>();

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
});

const imageRef = ref<HTMLImageElement | null>(null);
const cropper = ref<Cropper | null>(null);
const saving = ref(false);
const aspectRatio = ref(NaN); // 默认自由比例
const scaleX = ref(1);
const scaleY = ref(1);

// 初始化 Cropper
const initCropper = async () => {
    await nextTick();
    if (!imageRef.value || !props.imageUrl) return;

    // 如果已存在实例，先销毁
    if (cropper.value) {
        cropper.value.destroy();
    }

    cropper.value = new Cropper(imageRef.value, {
        viewMode: 1, // 限制裁剪框不超过画布
        dragMode: 'move', // 默认模式：移动画布
        aspectRatio: NaN, // 默认自由比例
        autoCropArea: 0.8, // 自动裁剪区域大小
        background: true, // 显示网格背景
        responsive: true,
        restore: false,
        checkCrossOrigin: false,
    });
};

// 监听弹窗关闭，销毁实例
const handleClose = () => {
    if (cropper.value) {
        cropper.value.destroy();
        cropper.value = null;
    }
    visible.value = false;
};

// 监听图片变化，重新初始化
watch(() => props.imageUrl, () => {
    if (visible.value) {
        initCropper();
    }
});

// === 工具栏操作 ===

const handleZoom = (ratio: number) => {
    cropper.value?.zoom(ratio);
};

const handleRotate = (degree: number) => {
    cropper.value?.rotate(degree);
};

const handleScaleX = () => {
    scaleX.value = scaleX.value === 1 ? -1 : 1;
    cropper.value?.scaleX(scaleX.value);
};

const handleScaleY = () => {
    scaleY.value = scaleY.value === 1 ? -1 : 1;
    cropper.value?.scaleY(scaleY.value);
};

const updateRatio = (val: number) => {
    cropper.value?.setAspectRatio(val);
};

const handleReset = () => {
    cropper.value?.reset();
    scaleX.value = 1;
    scaleY.value = 1;
    aspectRatio.value = NaN;
};

// === 保存逻辑 ===
const handleSave = () => {
    if (!cropper.value) return;

    saving.value = true;

    // 获取裁剪后的 Canvas
    const canvas = cropper.value.getCroppedCanvas({
        // 限制最大输出尺寸，避免过大导致性能问题
        maxWidth: 2048,
        maxHeight: 2048,
        imageSmoothingEnabled: true,
        imageSmoothingQuality: 'high',
    });

    if (!canvas) {
        saving.value = false;
        return;
    }

    // 转换为 Blob
    canvas.toBlob((blob) => {
        if (blob) {
            // 按照 ScriptEditor.vue 中 handleCropSave 的期望格式返回数组
            // 这里 frameType 默认为 'action' 或者你可以通过 props 传入当前选中的 frameType
            const result = [
                {
                    blob: blob,
                    frameType: 'action' // 通常裁剪用于动作序列修正，或者作为新的参考图
                }
            ];

            emit('save', result);
            handleClose();
        }
        saving.value = false;
    }, 'image/jpeg', 0.9); // 0.9 质量 JPEG
};

onUnmounted(() => {
    if (cropper.value) {
        cropper.value.destroy();
    }
});
</script>

<style scoped lang="less">
.cropper-wrapper {
    display: flex;
    flex-direction: column;
    gap: 16px;

    .crop-container {
        height: 450px;
        background-color: #f0f0f0;
        overflow: hidden;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 4px;

        img {
            max-width: 100%;
            /* 必须加上这个，否则 cropperjs 初始化会有问题 */
            display: block;
        }

        .empty-state {
            color: #999;
        }
    }

    .toolbar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 8px;
        background: var(--td-bg-color-secondarycontainer);
        border-radius: 4px;
        flex-wrap: wrap;
        gap: 8px;

        .tool-group {
            display: flex;
            gap: 8px;
            align-items: center;

            &.ratio-group {
                border-left: 1px solid var(--td-component-stroke);
                border-right: 1px solid var(--td-component-stroke);
                padding: 0 12px;
            }
        }
    }
}
</style>