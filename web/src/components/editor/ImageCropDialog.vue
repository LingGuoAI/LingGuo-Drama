<template>
    <t-dialog v-model:visible="visible" header="图片裁剪" width="600px" @confirm="handleSave">
        <div style="height: 300px; display: flex; align-items: center; justify-content: center; background: #eee;">
            <p>裁剪功能暂未集成 (Placeholder)</p>
        </div>
    </t-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
    modelValue: boolean;
    imageUrl: string;
}>();

const emit = defineEmits(['update:modelValue', 'save']);

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
});

const handleSave = () => {
    // Mock save
    emit('save', []);
    visible.value = false;
}
</script>