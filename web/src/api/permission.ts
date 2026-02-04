import type { MenuListResult } from '@/api/model/permissionModel';
import { request } from '@/utils/request';

const Api = {
  MenuList: '/getMenuList',
};

export function getMenuList() {
  return request.get<MenuListResult>({
    url: '/sys_base_menuses/getMenuList',
  });
}
