import axios from "axios";
import { localStoreUserDataKey } from ".";
import { message } from "ant-design-vue";
// import { useRouter } from "vue-router";

const dnsAxios = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 5000,
  responseType: "json",
  responseEncoding: "utf8",
});

// 添加请求拦截器: 请求前获取token
dnsAxios.interceptors.request.use(
  (config) => {
    const userdata = localStorage.getItem(localStoreUserDataKey);
    if (userdata) {
      const userdataObj = JSON.parse(userdata);
      config.headers["token"] = userdataObj.jwt_token;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 添加响应拦截器: 非200请求, 200请求中result为false则请求失败, 全局失败响应信息提示
dnsAxios.interceptors.response.use(
  (response) => {
    if (response.status !== 200 || !response.data.status) {
      message.error(`request_id: ${response.data.request_id}`);
      return Promise.reject(
        new Error(`request_id: ${response.data.request_id}`)
      );
    }
    return response;
  },
  (error) => {
    if (error.response) {
      message.error(`${error.response.config.url} ${error.response.status}`);
    }
    return Promise.reject(error);
  }
);

const request = {
  get(url, params = {}) {
    return dnsAxios.get(url, { params });
  },
  post(url, data = {}) {
    return dnsAxios.post(url, data);
  },
  put(url, data = {}) {
    return dnsAxios.put(url, data);
  },
  delete(url, params = {}) {
    return dnsAxios.delete(url, { params });
  },
};

export { dnsAxios };
export default request;
