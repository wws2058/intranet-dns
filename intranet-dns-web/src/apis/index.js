import request from "./request";
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
  await request.delete(url);
}

// put {name name_cn id api_ids id}
export async function updateSysRole(params) {
  await request.put("/api/v1/roles", params);
}

// post
export async function newSysRoles(params) {
  await request.post("/api/v1/roles", params);
}

export async function getRoleApis(roleId) {
  const url = `/api/v1/roles/${roleId}/apis`;
  const response = await request.get(url);
  return response.data;
}

// 系统api query params, page page_size, path, method, active
export async function getSysApis(params) {
  const response = await request.get("/api/v1/apis", params);
  return response.data;
}

// 获取所有的api
export async function getAllSysApis() {
  let currentPage = 1;
  let pageSize = 100;
  let apis = [];
  let totalPages = null;
  while (true) {
    const result = await getSysApis({ page: currentPage, page_size: pageSize });
    totalPages = Math.ceil(result.pages.total / pageSize);
    apis = apis.concat(result.data);
    if (currentPage >= totalPages) {
      break;
    }
    currentPage++;
  }
  return apis;
}

export async function updateSysApi(params) {
  const response = await request.put("/api/v1/apis", params);
  return response.data;
}
