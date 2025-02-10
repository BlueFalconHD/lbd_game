import axios from "axios";
import Cookies from "js-cookie";
import { env } from "process";

const API_URL = env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

const api = axios.create({
  baseURL: API_URL,
});

// Add a request interceptor to include the token in headers
api.interceptors.request.use(
  (config) => {
    const token = Cookies.get("token");
    if (token && config.headers) {
      config.headers["Authorization"] = token;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

export default api;
