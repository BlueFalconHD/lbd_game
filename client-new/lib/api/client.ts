import axios from "axios";
import Cookies from "js-cookie";
import { GetToken } from "../auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8040";

const validateStatus = (status: number) => {
  // If status is 401, do not reject the promise so login page can handle it
  return (status >= 200 && status < 300) || status === 401;
};

const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
  validateStatus: validateStatus,
});

apiClient.interceptors.request.use(
  (config) => {
    const token = GetToken();
    if (token && config.headers) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

export default apiClient;
