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
  if (response.status === 200 && response.data.status) {
    userdata.jwt_token = response.data.data.jwt_token;
  }
  return userdata;
}

// 校验jwt token是否有效
export function isTokenValid(token) {
  if (token === null || token === undefined) {
    return false;
  }
  const decoded = jwtDecode(token);
  if (!decoded.exp) {
    return false;
  }

  const currentTimestamp = Math.floor(Date.now() / 1000);
  return decoded.exp > currentTimestamp;
}
