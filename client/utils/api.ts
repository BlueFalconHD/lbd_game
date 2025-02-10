// utils/api.ts

import axios from "axios";
import Cookies from "js-cookie";

const API_URL = "http://localhost:8080"; // Update if your backend runs on a different host or port

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
