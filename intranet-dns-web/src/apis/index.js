import request, { dnsAxios } from "./request";
import { jwtDecode } from "jwt-decode";

export const localStoreUserDataKey = "userdata";

// 用户登录api, post body {name, password}, response body {name, jwt_token}
export async function userLogin(userMsg) {
  let userdata = {
    name: userMsg.name,
    jwt_token: "",
  };
  const response = await request.post("/api/v1/users/login", userMsg);
  userdata.jwt_token = response.data.data.jwt_token;
  return userdata;
}

// 校验jwt token是否有效
export function isTokenValid(token) {
  try {
    if (token === null || token === undefined) {
      return false;
    }
    const decoded = jwtDecode(token);
    if (!decoded.exp) {
      return false;
    }

    const currentTimestamp = Math.floor(Date.now() / 1000);
    return decoded.exp > currentTimestamp;
  } catch (error) {
    console.log("jwt token parse failed:", error);
    return false;
  }
}

// 获取系统角色, post body {page page_size name_cn}
export async function getSysRoles(queryRoleObj) {
  const response = await request.get("/api/v1/roles", queryRoleObj);
  return response.data;
}

export async function deleteSysRole(roleId) {
  const url = `/api/v1/roles/${roleId}`;
  await dnsAxios.delete(url);
}

export async function getRoleApis(roleId) {
  const url = `/api/v1/roles/${roleId}/apis`;
  const response = await dnsAxios.get(url);
  return response;
}

// 系统api query params, page page_size, path, method, active
export async function getSysApis(params) {
  const response = await request.get("/api/v1/apis", params);
  return response.data;
}

export async function updateSysApi(params) {
  const response = await request.put("/api/v1/apis", params);
  return response.data;
}
