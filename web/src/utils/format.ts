import { MessagePlugin } from "tdesign-vue-next";

// 格式化日期
export const formatDate = (date: string | Date): string => {
  if (!date) return "";
  const d = new Date(date);
  return d
    .toLocaleString("zh-CN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    })
    .replace(/\//g, "-");
};

// 格式化布尔值
export const formatBoolean = (value: boolean): string => {
  return value ? "是" : "否";
};

// 字典过滤
export const filterDict = (
  value: any,
  options: Array<{ label: string; value: any }>,
): string => {
  const option = options.find((item) => item.value === value);
  return option ? option.label : String(value);
};

// 数据源过滤
export const filterDataSource = (
  dataSource: Array<{ label: string; value: any }>,
  value: any,
): string | string[] => {
  if (Array.isArray(value)) {
    return value.map((v) => {
      const item = dataSource.find((d) => d.value === v);
      return item ? item.label : String(v);
    });
  } else {
    const item = dataSource.find((d) => d.value === value);
    return item ? item.label : String(value);
  }
};

export const getImageUrl = (url: string): string => {
  if (!url) return "";

  // 检查是否为完整的HTTP/HTTPS URL
  if (url.startsWith("http://") || url.startsWith("https://")) {
    return url;
  }

  // 相对路径则拼接API基础URL
  return import.meta.env.VITE_API_URL + `/${url}`;
};

// 文件下载
export const downloadFile = (url: string, filename?: string): void => {
  const link = document.createElement("a");
  link.href = url;
  link.target = "_blank";
  if (filename) {
    link.download = filename;
  }
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};

// 获取字典数据的模拟函数
export const getDictFunc = async (
  dictType: string,
): Promise<Array<{ label: string; value: any }>> => {
  // 这里应该调用实际的字典API
  console.log("获取字典数据:", dictType);
  return [];
};

// 返回图片数组用于预览
export const returnArrImg = (url: string | string[]): string[] => {
  if (Array.isArray(url)) {
    return url.map(getImageUrl);
  }
  return url ? [getImageUrl(url)] : [];
};
