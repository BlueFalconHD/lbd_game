import apiClient from "./client";
import { User } from "./types";

export const getUserStatistics = async (): Promise<any[]> => {
  const response = await apiClient.get<{ statistics: any[] }>(
    "/admin/stats/users",
  );
  return response.data.statistics;
};

export const resurrectUser = async (userId: number): Promise<void> => {
  await apiClient.put(`/admin/user/${userId}/resurrect`);
};

export const promoteUser = async (userId: number): Promise<void> => {
  await apiClient.put(`/superadmin/user/${userId}/promote`);
};

export const demoteUser = async (userId: number): Promise<void> => {
  await apiClient.put(`/superadmin/user/${userId}/demote`);
};

export const editPhrase = async (content: string): Promise<void> => {
  await apiClient.put("/admin/edit_phrase", { content: content });
};

export const unsubmitPhrase = async (): Promise<void> => {
  await apiClient.put("/admin/unsubmit_phrase");
};

export const manualReset = async (openTime: number): Promise<void> => {
  await apiClient.put("/admin/manual_reset", { open_time: openTime });
};
