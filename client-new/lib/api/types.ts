export interface User {
  id: number;
  username: string;
  email: string;
  privilege: number;
  isEliminated: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Phrase {
  id: number;
  content: string;
  submittedBy: number;
  date: string;
  createdAt: string;
  updatedAt: string;
}

export interface Verification {
  id: number;
  verifiedUserId: number;
  verifierId: number;
  date: string;
  createdAt: string;
  updatedAt: string;
}

export interface TodayVerification {
  verification_id: number;
  verifier_id: number;
  verifier_name: string;
  verified_id: number;
  verified_name: string;
  created_at: string; // ISO 8601 date string
}
export type TodayVerificationsResponse = TodayVerification[];

export interface UnverifiedUser {
  id: number;
  username: string;
}
export type UnverifiedUsersResponse = UnverifiedUser[];

export interface AuthResponse {
  token: string;
}

export interface LoginResponse {
  token: string;
}

export interface SignUpResponse {} // doesn't return anything unique

export interface PhraseResponse {
  content: string;
  submittedBy: string;
  time_until_reset: number; // Seconds
}

export interface PrivelegeResponse {
  privilege: number;
}

export interface ApiResponse<T> {
  message?: string;
  error?: string;
  data?: T;
  [key: string]: any;
}
