import { AxiosResponse } from "axios";
import apiClient from "./client";
import {
  AuthResponse,
  LoginResponse,
  SignUpResponse,
  User,
  ApiResponse,
  PrivelegeResponse,
} from "./types";

interface SignUpData {
  username: string;
  password: string;
}

interface LoginData {
  username: string;
  password: string;
}

export const signUp = async (
  data: SignUpData,
): Promise<ApiResponse<SignUpResponse>> => {
  return (await apiClient.post("/signup", data)).data;
};

export const login = async (
  data: LoginData,
): Promise<ApiResponse<LoginResponse>> => {
  return (await apiClient.post("/login", data)).data;
};

export const getPrivelege = async (): Promise<
  ApiResponse<PrivelegeResponse>
> => {
  return (await apiClient.get("/privilege")).data;
};
