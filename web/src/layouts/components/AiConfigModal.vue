<template>
    <t-dialog v-model:visible="visible" header="AI 服务配置" width="1000px" placement="center" :footer="false"
        destroy-on-close>
        <div class="ai-config-container">
            <div class="header-actions">
                <t-tabs v-model="activeTab" @change="handleTabChange">
                    <t-tab-panel value="text" label="文本模型" />
                    <t-tab-panel value="image" label="图片模型" />
                    <t-tab-panel value="video" label="视频模型" />
                </t-tabs>
                <t-button theme="primary" @click="showCreateDialog">
                    <template #icon><plus-icon /></template>
                    添加配置
                </t-button>
            </div>

            <t-table row-key="id" :data="configs" :columns="columns" :loading="loading" hover stripe
                style="margin-top: 16px;">
                <template #model="{ row }">
                    <t-space size="small" wrap>
                        <t-tag v-for="m in row.model" :key="m" theme="primary" variant="light">{{ m }}</t-tag>
                    </t-space>
                </template>
                <template #is_active="{ row }">
                    <t-switch v-model="row.is_active" :custom-value="[1, 0]" @change="handleToggleActive(row)" />
                </template>
                <template #op="{ row }">
                    <t-space size="small">
                        <t-button v-if="activeTab === 'text'" variant="text" theme="primary"
                            @click="handleTest(row)">测试</t-button>
                        <t-button variant="text" theme="primary" @click="handleEdit(row)">编辑</t-button>
                        <t-button variant="text" theme="danger" @click="handleDelete(row)">删除</t-button>
                    </t-space>
                </template>
            </t-table>
        </div>
    </t-dialog>

    <t-dialog v-model:visible="dialogVisible" :header="isEdit ? '编辑配置' : '添加配置'" width="650px" placement="center"
        attach="body" destroy-on-close :confirm-btn="submitting ? '保存中...' : (isEdit ? '保存' : '创建')"
        @confirm="handleSubmit">
        <t-form ref="formRef" :data="form" :rules="rules" label-width="120px" label-align="right">
            <t-form-item label="配置名称" name="name">
                <t-input v-model="form.name" placeholder="请输入配置名称，如：OpenAI-Text" />
            </t-form-item>

            <t-form-item label="厂商提供商" name="provider">
                <t-select v-model="form.provider" placeholder="请选择服务提供商" @change="handleProviderChange">
                    <t-option v-for="provider in availableProviders" :key="provider.id" :label="provider.name"
                        :value="provider.id" :disabled="provider.disabled" />
                </t-select>
                <template #help>目前可用的厂商类型</template>
            </t-form-item>

            <t-form-item label="优先级" name="priority">
                <t-input-number v-model="form.priority" :min="0" :max="100" :step="1" style="width: 100%" />
                <template #help>数字越大优先级越高</template>
            </t-form-item>

            <t-form-item label="支持的模型" name="model">
                <t-select v-model="form.model" placeholder="请选择或输入模型名后回车" multiple creatable filterable
                    :min-collapsed-num="3">
                    <t-option v-for="model in availableModels" :key="model" :label="model" :value="model" />
                </t-select>
                <template #help>支持多选，可输入自定义模型名称后按回车添加</template>
            </t-form-item>

            <t-form-item label="接口地址" name="base_url">
                <t-input v-model="form.base_url" placeholder="https://api.example.com/v1" />
                <template #help>
                    包含域名和协议的完整基础路径。<br />
                    实际请求路径示例: {{ fullEndpointExample }}
                </template>
            </t-form-item>

            <t-form-item label="API Key" name="api_key">
                <t-input v-model="form.api_key" type="password" clearable placeholder="sk-..." />
            </t-form-item>

            <t-form-item v-if="isEdit" label="状态" name="is_active">
                <t-switch v-model="form.is_active" :custom-value="[1, 0]" />
            </t-form-item>

            <t-form-item style="margin-top: 24px;">
                <t-button variant="outline" theme="default" :loading="testing" @click="testConnection">
                    测试连通性
                </t-button>
            </t-form-item>
        </t-form>
    </t-dialog>

    <t-dialog v-model:visible="testResultVisible" header="测试结果" width="600px" placement="center" :footer="false"
        attach="body" @close="stopPolling">
        <div v-if="testPolling" class="polling-container">
            <t-loading size="medium" text="AI 正在处理中，请耐心稍候..." />
        </div>
        <div v-else class="result-content">
            <div v-if="testResultType === 'text'" class="text-result">
                {{ testResultData }}
            </div>
            <div v-else-if="testResultType === 'image'" class="image-result">
                <t-image :src="testResultData" fit="contain"
                    style="width: 100%; max-height: 400px; border-radius: 8px;" />
            </div>
            <div v-else-if="testResultType === 'video'" class="video-result">
                <video :src="testResultData" controls autoplay playsinline loop
                    style="width: 100%; max-height: 400px; border-radius: 8px; outline: none; background: #000;">
                    您的浏览器不支持视频播放。
                </video>
            </div>
            <div v-else class="text-result">
                {{ testResultData }}
            </div>
        </div>
    </t-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue';
import { MessagePlugin, DialogPlugin, FormRule } from 'tdesign-vue-next';
import { PlusIcon } from 'tdesign-icons-vue-next';
import { request } from "@/utils/request"; // 引入底层 request 以便轮询

import {
    getAiConfigList,
    createAiConfig,
    updateAiConfig,
    deleteAiConfig,
    testTextConfig,
    testImageConfig,
    testVideoConfig
} from '@/api/ai_config';

import type { AIServiceConfig, AIServiceType, CreateAIConfigRequest, UpdateAIConfigRequest } from '@/types/ai';
import { getImageUrl } from '@/utils/format';

const props = defineProps({
    visible: Boolean
});

const emit = defineEmits(['update:visible']);

const visible = computed({
    get: () => props.visible,
    set: (val) => emit('update:visible', val)
});

// --- 列表与核心状态 ---
const activeTab = ref<AIServiceType>('text');
const loading = ref(false);
const configs = ref<any[]>([]);

const columns = [
    { colKey: 'name', title: '配置名称', width: 150 },
    { colKey: 'provider', title: '厂商', width: 120 },
    { colKey: 'model', title: '模型列表', width: 280 },
    { colKey: 'priority', title: '优先级', width: 100, align: 'center' },
    { colKey: 'is_active', title: '状态', width: 100 },
    { colKey: 'op', title: '操作', width: 200, fixed: 'right' }
];

// --- 表单与弹框状态 ---
const dialogVisible = ref(false);
const isEdit = ref(false);
const editingId = ref<number>();
const formRef = ref<any>(null);
const submitting = ref(false);
const testing = ref(false);

const form = reactive({
    service_type: 'text',
    provider: '',
    name: '',
    base_url: '',
    api_key: '',
    model: [],
    priority: 0,
    is_active: 1, // 默认设为 1 (启用)
});

// --- 常量配置 ---
interface ProviderConfig {
    id: string;
    name: string;
    models: string[];
    disabled?: boolean;
}

const providerConfigs: Record<AIServiceType, ProviderConfig[]> = {
    text: [
        { id: 'openai', name: 'OpenAI', models: ['gpt-5.2', 'gemini-3-flash-preview'] },
        { id: 'getgoapi', name: 'GetGo API', models: ['gemini-3-flash-preview', 'claude-sonnet-4-5-20250929', 'doubao-seed-1-8-251228'] },
        { id: 'gemini', name: 'Google Gemini', models: ['gemini-2.5-pro', 'gemini-3-flash-preview'] },
    ],
    image: [
        { id: 'volcengine', name: '火山引擎', models: ['doubao-seedream-4-5-251128', 'doubao-seedream-4-0-250828'] },
        { id: 'getgoapi', name: 'GetGo API', models: ['doubao-seedream-4-5-251128', 'nano-banana-pro'] },
        { id: 'gemini', name: 'Google Gemini', models: ['gemini-3-pro-image-preview'] },
        { id: 'openai', name: 'OpenAI', models: ['dall-e-3', 'dall-e-2'] },
    ],
    video: [
        { id: 'volces', name: '火山引擎', models: ['doubao-seedance-1-5-pro-251215', 'doubao-seedance-1-0-lite-i2v-250428', 'doubao-seedance-1-0-pro-250528'] },
        { id: 'getgoapi', name: 'GetGo API', models: ['doubao-seedance-1-5-pro-251215', 'sora-2', 'sora-2-pro'] },
        { id: 'openai', name: 'OpenAI', models: ['sora-2', 'sora-2-pro'] },
    ],
};

const availableProviders = computed(() => providerConfigs[form.service_type as AIServiceType] || []);

const availableModels = computed(() => {
    if (!form.provider) return [];
    const providerDef = availableProviders.value.find(p => p.id === form.provider);
    return providerDef ? providerDef.models : [];
});

const fullEndpointExample = computed(() => {
    const baseUrl = form.base_url || 'https://api.example.com';
    const provider = form.provider;
    const serviceType = form.service_type;
    let endpoint = '';

    if (serviceType === 'text') {
        endpoint = (provider === 'gemini') ? '/v1beta/models/{model}:generateContent' : '/chat/completions';
    } else if (serviceType === 'image') {
        endpoint = (provider === 'gemini') ? '/v1beta/models/{model}:generateContent' : '/images/generations';
    } else if (serviceType === 'video') {
        if (provider === 'doubao' || provider === 'volcengine' || provider === 'volces') endpoint = '/contents/generations/tasks';
        else if (provider === 'openai') endpoint = '/videos';
        else endpoint = '/video/generations';
    }
    return baseUrl + endpoint;
});

// TDesign 校验规则
const rules: Record<string, FormRule[]> = {
    name: [{ required: true, message: '请输入配置名称', type: 'error' }],
    provider: [{ required: true, message: '请选择厂商', type: 'error' }],
    base_url: [
        { required: true, message: '请输入 Base URL', type: 'error' },
        { url: true, message: '请输入正确的 URL 格式', type: 'error' }
    ],
    api_key: [{ required: true, message: '请输入 API Key', type: 'error' }],
    model: [
        { required: true, message: '请至少选择一个模型', type: 'error' },
        {
            validator: (val) => (Array.isArray(val) && val.length > 0) || (typeof val === 'string' && val.length > 0),
            message: '请至少选择一个模型',
            type: 'error'
        }
    ]
};

// --- 方法 ---

const loadConfigs = async () => {
    loading.value = true;
    try {
        const res: any = await getAiConfigList({
            service_type: activeTab.value,
            per_page: 100
        });
        configs.value = res.data?.list || [];
    } catch (error: any) {
        MessagePlugin.error(error.message || '加载失败');
    } finally {
        loading.value = false;
    }
};

watch(visible, (val) => {
    if (val) loadConfigs();
});

const handleTabChange = (value: AIServiceType) => {
    activeTab.value = value;
    loadConfigs();
};

const generateConfigName = (provider: string, serviceType: AIServiceType) => {
    const providerNames: Record<string, string> = { getgoapi: 'GetGo', openai: 'OpenAI', gemini: 'Gemini' };
    const serviceNames: Record<AIServiceType, string> = { text: '文本', image: '图片', video: '视频' };
    const randomNum = Math.floor(Math.random() * 10000).toString().padStart(4, '0');
    return `${providerNames[provider] || provider}-${serviceNames[serviceType] || serviceType}-${randomNum}`;
};

const resetForm = () => {
    Object.assign(form, {
        service_type: activeTab.value,
        provider: '',
        name: '',
        base_url: '',
        api_key: '',
        model: [],
        priority: 0,
        is_active: 1,
    });
    formRef.value?.clearValidate();
};

const showCreateDialog = () => {
    isEdit.value = false;
    editingId.value = undefined;
    resetForm();
    form.service_type = activeTab.value;
    form.provider = 'getgoapi';
    form.base_url = 'https://api.getgoapi.com/v1';
    form.name = generateConfigName('getgoapi', activeTab.value);
    dialogVisible.value = true;
};

const handleEdit = (config: any) => {
    isEdit.value = true;
    editingId.value = config.id;
    Object.assign(form, {
        service_type: config.service_type,
        provider: config.provider || 'getgoapi',
        name: config.name,
        base_url: config.base_url,
        api_key: config.api_key,
        model: Array.isArray(config.model) ? config.model : [config.model],
        priority: config.priority || 0,
        is_active: config.is_active,
    });
    dialogVisible.value = true;
};

const handleProviderChange = () => {
    form.model = [];
    form.base_url = (form.provider === 'gemini') ? 'https://generativelanguage.googleapis.com/v1beta' : 'https://api.getgoapi.com/v1';
    if (!isEdit.value) {
        form.name = generateConfigName(form.provider, form.service_type as AIServiceType);
    }
};

const handleDelete = (config: any) => {
    const confirmDia = DialogPlugin.confirm({
        header: '警告',
        body: `确定要删除配置 [${config.name}] 吗？此操作无法恢复。`,
        theme: 'danger',
        onConfirm: async () => {
            try {
                await deleteAiConfig(config.id);
                MessagePlugin.success('删除成功');
                loadConfigs();
            } catch (error: any) {
                MessagePlugin.error(error.message || '删除失败');
            } finally {
                confirmDia.destroy();
            }
        }
    });
};

const handleToggleActive = async (config: any) => {
    try {
        await updateAiConfig(config.id, {
            ...config,
            is_active: config.is_active
        });
        MessagePlugin.success(config.is_active === 1 ? '已启用配置' : '已禁用配置');
    } catch (error: any) {
        config.is_active = config.is_active === 1 ? 0 : 1;
        MessagePlugin.error(error.message || '操作失败');
    }
};

const handleSubmit = async () => {
    const validateResult = await formRef.value.validate();
    if (validateResult !== true) return;

    submitting.value = true;
    try {
        if (isEdit.value && editingId.value) {
            await updateAiConfig(editingId.value, form);
            MessagePlugin.success('更新成功');
        } else {
            await createAiConfig(form);
            MessagePlugin.success('创建成功');
        }
        dialogVisible.value = false;
        loadConfigs();
    } catch (error: any) {
        MessagePlugin.error(error.message || '操作失败');
    } finally {
        submitting.value = false;
    }
};

// ==========================================
// 🔴 测试连接与轮询逻辑 
// ==========================================
const testResultVisible = ref(false);
const testPolling = ref(false);
const testResultType = ref('');
const testResultData = ref('');
let pollTimer: any = null;

// 清除轮询器
const stopPolling = () => {
    if (pollTimer) {
        clearTimeout(pollTimer);
        pollTimer = null;
    }
};

// 开始轮询 TaskID 的结果
const startPolling = (taskId: number | string, type: string) => {
    testResultVisible.value = true;
    testPolling.value = true;
    testResultType.value = type;
    testResultData.value = '';

    const poll = async () => {
        try {
            const res: any = await request.get({ url: `/tasks/${taskId}` });
            const task = res.data;

            if ((task.status === 2 || task.status === 3) && task.process === 100) {
                testPolling.value = false;

                let resultObj: any = {};
                try {
                    resultObj = JSON.parse(task.result || '{}');
                } catch (e) {
                    resultObj = { raw: task.result };
                }

                const rawUrl = resultObj.video_url || resultObj.url || resultObj.image_url || task.result;

                if (type === 'text') {
                    testResultData.value = resultObj.reply || task.result;
                } else {
                    // 图片和视频都需要转换路径
                    testResultData.value = getImageUrl(rawUrl);
                    console.log('testResultData', testResultData)
                }

                MessagePlugin.success('AI 处理完成！');
                stopPolling(); // 停止定时器

            } else if (task.status === 4 || task.status === -1) {
                // 失败状态处理
                testPolling.value = false;
                testResultData.value = task.error_msg || '任务执行失败';
                MessagePlugin.error(`测试失败: ${testResultData.value}`);
                stopPolling();
            } else {
                // 继续轮询
                pollTimer = setTimeout(poll, 3000);
            }
        } catch (error) {
            testPolling.value = false;
            MessagePlugin.error('查询任务状态异常');
            stopPolling();
        }
    };

    poll();
};

const executeTest = async (payload: any, serviceType: string) => {
    testing.value = true;
    stopPolling(); // 确保清理旧的轮询

    try {
        let res: any;
        if (serviceType === 'text') {
            res = await testTextConfig(payload);
        } else if (serviceType === 'image') {
            res = await testImageConfig(payload);
        } else if (serviceType === 'video') {
            res = await testVideoConfig(payload);
        }

        // 💡 核心分发逻辑
        // 如果后端同步直接返回了数据 (如 reply 或 image_url)
        if (res.data?.reply || res.data?.image_url) {
            testResultVisible.value = true;
            testPolling.value = false;
            testResultType.value = serviceType;
            testResultData.value = res.data.reply || res.data.image_url;
            MessagePlugin.success('测试成功！');
        }
        // 如果后端走的是异步队列，返回了 task_id
        else if (res.data?.task_id) {
            MessagePlugin.info(`任务已提交 [TaskID: ${res.data.task_id}]，正在等待结果...`);
            startPolling(res.data.task_id, serviceType);
        }
        // 其他兜底情况
        else {
            MessagePlugin.success('连接测试成功，但未返回内容。');
        }

    } catch (error: any) {
        MessagePlugin.error(error.message || '连接测试请求失败');
    } finally {
        testing.value = false;
    }
};

const testConnection = async () => {
    const validateResult = await formRef.value.validate();
    if (validateResult !== true) return;
    executeTest(
        { base_url: form.base_url, api_key: form.api_key, model: form.model, provider: form.provider },
        form.service_type
    );
};

const handleTest = (config: any) => {
    executeTest(
        { base_url: config.base_url, api_key: config.api_key, model: config.model, provider: config.provider },
        config.service_type
    );
};

// 处理资源 URL
const getResourceUrl = (url: string) => {
    if (!url) return '';
    if (url.startsWith('http') || url.startsWith('data:')) return url;
    const baseUrl = (import.meta.env.VITE_APP_RESOURCE_URL || '').replace(/\/$/, '');
    return `${baseUrl}/${url.replace(/^\//, '')}`;
};
</script>

<style scoped>
.header-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--td-component-stroke);
    padding-bottom: 12px;
}

:deep(.t-form__help) {
    margin-top: 4px;
    line-height: 1.4;
    color: var(--td-text-color-placeholder);
}

.polling-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 150px;
}

.result-content {
    min-height: 100px;
    padding: 16px;
    background: var(--td-bg-color-container-active);
    border-radius: 8px;
    display: flex;
    justify-content: center;
    align-items: center;
}

.text-result {
    font-size: 14px;
    line-height: 1.6;
    white-space: pre-wrap;
    word-break: break-all;
    width: 100%;
}
</style>