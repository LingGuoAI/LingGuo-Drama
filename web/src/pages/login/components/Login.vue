<template>
  <t-form ref="form" :class="['item-container', `login-${type}`]" :data="formData" :rules="FORM_RULES" label-width="0"
    @submit="onSubmit">
    <template v-if="type == 'password'">
      <t-form-item name="account">
        <t-input v-model="formData.account" size="large" :placeholder="`${t('pages.login.input.account')}`">
          <template #prefix-icon>
            <t-icon name="user" />
          </template>
        </t-input>
      </t-form-item>

      <t-form-item name="password">
        <t-input v-model="formData.password" size="large" :type="showPsw ? 'text' : 'password'" clearable
          :placeholder="`${t('pages.login.input.password')}`">
          <template #prefix-icon>
            <t-icon name="lock-on" />
          </template>
          <template #suffix-icon>
            <t-icon :name="showPsw ? 'browse' : 'browse-off'" @click="showPsw = !showPsw" />
          </template>
        </t-input>
      </t-form-item>
    </template>

    <t-form-item v-if="type !== 'qrcode'" class="btn-container">
      <t-button block size="large" type="submit" :loading="isLoading" :disabled="isLoading">
        {{ isLoading ? "登录中" : "登录" }}
      </t-button>
    </t-form-item>
  </t-form>
</template>

<script setup lang="ts">
import type { FormInstanceFunctions, FormRule, SubmitContext } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { t } from '@/locales';
import { useUserStore } from '@/store';

const userStore = useUserStore();

const INITIAL_DATA = {
  account: '',
  password: ''
};

const FORM_RULES: Record<string, FormRule[]> = {
  account: [{ required: true, message: t('pages.login.required.account'), type: 'error' }],
  password: [{ required: true, message: t('pages.login.required.password'), type: 'error' }],
};

const type = ref('password');
const form = ref<FormInstanceFunctions>();
const formData = ref({ ...INITIAL_DATA });
const showPsw = ref(false);
const isLoading = ref(false);

const router = useRouter();
const route = useRoute();

// 错误信息映射
const getErrorMessage = (error: any): string => {
  // 如果有具体的错误消息，直接使用
  if (error?.message) {
    return error.message;
  }

  // 根据错误代码返回对应的错误信息
  if (error?.code) {
    switch (error.code) {
      case 'USER_NOT_FOUND':
        return t('pages.login.errors.userNotFound');
      case 'INVALID_PASSWORD':
        return t('pages.login.errors.invalidPassword');
      case 'ACCOUNT_LOCKED':
        return t('pages.login.errors.accountLocked');
      case 'NETWORK_ERROR':
        return t('pages.login.errors.networkError');
      case 'SERVER_ERROR':
        return t('pages.login.errors.serverError');
      default:
        return t('pages.login.errors.loginFailed');
    }
  }

  // 根据HTTP状态码返回错误信息
  if (error?.response?.status) {
    switch (error.response.status) {
      case 400:
        return t('pages.login.errors.badRequest');
      case 401:
        return t('pages.login.errors.unauthorized');
      case 403:
        return t('pages.login.errors.forbidden');
      case 404:
        return t('pages.login.errors.notFound');
      case 500:
        return t('pages.login.errors.serverError');
      case 503:
        return t('pages.login.errors.serviceUnavailable');
      default:
        return t('pages.login.errors.unknownError');
    }
  }

  // 默认错误信息
  return t('pages.login.errors.defaultError');
};

const onSubmit = async (ctx: SubmitContext) => {
  if (ctx.validateResult === true) {
    isLoading.value = true;

    try {
      const loginParams = {
        login_id: formData.value.account.trim(),
        password: formData.value.password
      };

      const success = await userStore.login(loginParams);

      if (success) {
        MessagePlugin.success("登录成功");

        // 延迟跳转，让用户看到成功提示
        setTimeout(() => {
          const redirect = route.query.redirect as string;
          const redirectUrl = redirect ? decodeURIComponent(redirect) : '/admin/projects';
          router.push(redirectUrl);
        }, 500);
      }
    } catch (error: any) {
      console.error('登录失败:', error);

      // 使用错误信息映射函数
      const errorMessage = getErrorMessage(error);
      MessagePlugin.error({
        content: errorMessage,
        duration: 3000,
        closeBtn: true
      });

      // 如果是密码错误，清空密码输入框
      if (error?.code === 'INVALID_PASSWORD') {
        formData.value.password = '';
      }
    } finally {
      isLoading.value = false;
    }
  }
};
</script>

<style lang="less" scoped>
@import '../index.less';
</style>