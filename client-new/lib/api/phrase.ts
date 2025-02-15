import { AxiosResponse } from "axios";
import { pipe } from "../utils/pipe";
import apiClient from "./client";
import { Phrase, ApiResponse, PhraseResponse } from "./types";

export const getCurrentPhrase = async (): Promise<
  ApiResponse<PhraseResponse>
> => {
  return pipe<AxiosResponse>(
    await apiClient.get<ApiResponse<PhraseResponse>>("/phrase"),
  ).data;
};

export const submitPhrase = async (content: string): Promise<void> => {
  await apiClient.post("/phrase", { content });
};

export const canSubmitPhrase = async (): Promise<boolean> => {
  return pipe<AxiosResponse>(await apiClient.get("/can_submit_phrase")).data
    .can_submit;
};
