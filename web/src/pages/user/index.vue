<template>
  <t-row :gutter="[24, 24]">
    <t-col :flex="3">
      <div class="user-left-greeting">
        <div>
          Hi，{{userInfo.username}}。
          <span class="regular"> 所有的伟大，都源于一个勇敢的开始。</span>
        </div>
      </div>

      <t-card class="user-info-list" title="个人信息" :bordered="false">
        <template #actions>
          <t-button theme="default" shape="square" variant="text">
            <t-icon name="ellipsis" />
          </t-button>
        </template>
        <t-descriptions :column="4" item-layout="vertical">
          <t-descriptions-item v-for="(item, index) in userInfoList" :key="index" :label="item.title">
            {{ item.content }}
          </t-descriptions-item>
        </t-descriptions>
      </t-card>
    </t-col>
  </t-row>
</template>
<script lang="ts">
export default {
  name: 'UserIndex',
};
</script>
<script setup lang="ts">
const userInfo = JSON.parse(localStorage.getItem('userInfo'));

const userInfoList = [
  {
    title:'用户名',
    content:userInfo.username
  },
  {
    title:'手机号',
    content:userInfo.mobile
  },
  {
    title:'邮箱',
    content:userInfo.email
  },
  {
    title:'创建时间',
    content:userInfo.createdAt
  }
]


</script>

<style lang="less" scoped>


.user-left-greeting {
  padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl);
  font: var(--td-font-title-large);
  background: var(--td-bg-color-container);
  color: var(--td-text-color-primary);
  text-align: left;
  border-radius: var(--td-radius-medium);
  display: flex;
  justify-content: space-between;
  align-items: center;

  .regular {
    margin-left: var(--td-comp-margin-xl);
    font: var(--td-font-body-medium);
  }

  .logo {
    width: 168px;
  }
}

.user-info-list {
  margin-top: var(--td-comp-margin-l);
  padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl);

  .content {
    width: 90%;
  }

  :deep(.t-card__header) {
    padding: 0;
  }

  :deep(.t-card__body) {
    padding: 0;
  }

}

.t-descriptions {
  margin-top: 24px;
}
</style>
