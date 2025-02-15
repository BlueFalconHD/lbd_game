import apiClient from "./client";
import {
  ApiResponse,
  TodayVerificationsResponse,
  UnverifiedUsersResponse,
  Verification,
} from "./types";

export const verifyUser = async (verifiedUserId: number): Promise<void> => {
  await apiClient.post("/verify", { verified_user_id: verifiedUserId });
};

export const getTodayVerifications = async (): Promise<
  ApiResponse<TodayVerificationsResponse>
> => {
  const response =
    await apiClient.get<TodayVerificationsResponse>("/verifications");
  console.log("TodayVerifications", response);
  return response;
};

export const getUnverifiedUsers = async (): Promise<
  ApiResponse<UnverifiedUsersResponse>
> => {
  const response =
    await apiClient.get<UnverifiedUsersResponse>("/unverified_users");
  console.log("Unverified", response);
  return response;
};
