import { defineStore } from 'pinia';
import { MessagePlugin } from "tdesign-vue-next";

import { usePermissionStore } from '@/store';
import type { UserInfo } from '@/types/interface';
import { request } from '@/utils/request';

const InitUserInfo: UserInfo = {
  name: '', // 用户名，用于展示在页面右上角头像处
  roles: ['all'], // 前端权限模型使用 如果使用请配置modules/permission-fe.ts使用
};

// 定义错误类型
interface LoginError {
  code?: string;
  message?: string;
  response?: {
    status?: number;
    data?: any;
  };
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '', // 默认token不走权限
    userInfo: {} as UserInfo,
    loginAttempts: 0, // 登录尝试次数
    lastLoginAttempt: null as Date | null, // 最后一次登录尝试时间
  }),
  getters: {
    roles: (state) => {
      return state.userInfo?.roles;
    },
    isLoggedIn: (state) => {
      return !!state.token;
    },
  },
  actions: {
    async login(loginParams: any): Promise<boolean> {
      try {
        // 检查登录尝试次数限制
        if (this.loginAttempts >= 5 && this.lastLoginAttempt) {
          const timeSinceLastAttempt = Date.now() - this.lastLoginAttempt.getTime();
          const waitTime = 5 * 60 * 1000; // 5分钟

          if (timeSinceLastAttempt < waitTime) {
            const remainingTime = Math.ceil((waitTime - timeSinceLastAttempt) / 1000 / 60);
            throw {
              code: 'TOO_MANY_ATTEMPTS',
              message: `登录尝试次数过多，请${remainingTime}分钟后再试`
            };
          } else {
            // 重置登录尝试次数
            this.loginAttempts = 0;
          }
        }

        const response = await request.post<any>({
          url: '/auth/login/using-phone',
          data: loginParams
        });

        if (response.code === 0) {
          // 登录成功
          this.token = response.data.token;
          this.userInfo = response.data.user_info;

          // 保存到本地存储
          localStorage.setItem('token', response.data.token);
          localStorage.setItem('userInfo', JSON.stringify(response.data.user_info));

          // 重置登录尝试次数
          this.loginAttempts = 0;
          this.lastLoginAttempt = null;

          return true;
        } else {
          // 服务器返回错误
          this.loginAttempts++;
          this.lastLoginAttempt = new Date();

          // 根据返回的错误代码抛出相应的错误
          const error: LoginError = {
            code: response.code,
            message: response.message || '登录失败'
          };

          // 特殊错误码处理
          switch (response.code) {
            case 1001:
              error.code = 'USER_NOT_FOUND';
              error.message = '用户不存在或者密码错误';
              break;
            default:
              // 使用服务器返回的消息
              break;
          }

          throw error;
        }
      } catch (error: any) {
        // 网络错误或其他错误
        if (!error.code) {
          // 这是一个网络错误
          const networkError: LoginError = {
            code: 'NETWORK_ERROR',
            message: '网络连接失败，请检查网络设置',
            response: error.response
          };

          if (error.response) {
            switch (error.response.status) {
              case 400:
                networkError.message = '请求参数错误';
                break;
              case 401:
                networkError.message = '认证失败';
                break;
              case 403:
                networkError.message = '没有权限访问';
                break;
              case 404:
                networkError.message = '服务接口不存在';
                break;
              case 500:
                networkError.code = 'SERVER_ERROR';
                networkError.message = '服务器内部错误';
                break;
              case 502:
                networkError.message = '网关错误';
                break;
              case 503:
                networkError.message = '服务暂时不可用';
                break;
              default:
                networkError.message = `服务器错误 (${error.response.status})`;
            }
          } else if (error.message === 'Network Error') {
            networkError.message = '无法连接到服务器，请检查网络连接';
          }

          throw networkError;
        }

        // 重新抛出错误
        throw error;
      }
    },

    async register(registerParams: any): Promise<boolean> {
      try {
        const response = await request.post<any>({
          url: '/auth/register',
          data: registerParams
        });

        if (response.code === 201 || response.code === 0 || response.success) {
          MessagePlugin.success("注册成功");
          return true;
        } else {
          MessagePlugin.error(response.message || '注册失败');
          return false;
        }
      } catch (error: any) {
        console.error('注册失败:', error);
        MessagePlugin.error(error.message || '网络错误，请稍后尝试');
        return false;
      }
    },

    async getUserInfo() {
      try {
        const userInfoStr = localStorage.getItem('userInfo');
        if (userInfoStr) {
          const users = JSON.parse(userInfoStr);
          users.name = users.username || users.name;
          users.roles = users.roles || ['all'];
          this.userInfo = users;
        } else {
          // 如果本地没有用户信息，可以从服务器获取
          const response = await request.get<any>({
            url: '/user/info'
          });

          if (response.code === 0) {
            this.userInfo = response.data;
            localStorage.setItem('userInfo', JSON.stringify(response.data));
          }
        }
      } catch (error) {
        console.error('获取用户信息失败:', error);
        // 如果获取失败，可能需要重新登录
        // this.logout();
        throw error;
      }
    },

    async logout() {
      // 清理本地存储和状态
      localStorage.removeItem("token");
      localStorage.removeItem("userInfo");
      this.token = '';
      this.userInfo = {} as UserInfo;
      this.loginAttempts = 0;
      this.lastLoginAttempt = null;
      // try {
      //   // 调用后端登出接口（如果有的话）
      //   await request.post({
      //     url: '/auth/logout'
      //   });
      // } catch (error) {
      //   // 即使后端登出失败，也要清理本地状态
      //   console.error('登出请求失败:', error);
      // } finally {
      //   // 清理本地存储和状态
      //   localStorage.removeItem("token");
      //   localStorage.removeItem("userInfo");
      //   this.token = '';
      //   this.userInfo = {} as UserInfo;
      //   this.loginAttempts = 0;
      //   this.lastLoginAttempt = null;
      // }
    },

    // 检查Token是否有效
    async checkToken(): Promise<boolean> {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          return false;
        }

        // 可以调用一个验证token的接口
        const response = await request.get({
          url: '/auth/check-token'
        });

        return response.code === 0;
      } catch (error) {
        return false;
      }
    }
  },
  // persist: {
  //   afterRestore: () => {
  //     const permissionStore = usePermissionStore();
  //     permissionStore.initRoutes();
  //   },
  //   key: 'user',
  //   paths: ['token'],
  // },
});