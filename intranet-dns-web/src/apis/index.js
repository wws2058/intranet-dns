import request from "./request";

export const localStoreUserDataKey = "userdata";

// 获取所有的api
export async function getAllSysApis() {
  let currentPage = 1;
  let pageSize = 100;
  let apis = [];
  let totalPages = null;
  while (true) {
    const response = await request.get("/api/v1/apis", {
      page: currentPage,
      page_size: pageSize,
    });
    const result = response.data;
    totalPages = Math.ceil(result.pages.total / pageSize);
    apis = apis.concat(result.data);
    if (currentPage >= totalPages) {
      break;
    }
    currentPage++;
  }
  return apis;
}

// 获取所有的系统橘色
export async function getAllSysRoles() {
  let currentPage = 1;
  let pageSize = 100;
  let roles = [];
  let totalPages = null;
  while (true) {
    const response = await request.get("/api/v1/roles", {
      page: currentPage,
      page_size: pageSize,
    });
    const result = response.data;
    totalPages = Math.ceil(result.pages.total / pageSize);
    roles = roles.concat(result.data);
    if (currentPage >= totalPages) {
      break;
    }
    currentPage++;
  }
  return roles;
}
