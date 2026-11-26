import type { IUser } from "@/entities/user";

export type ILoginData = {
  email: string;
  password: string;
};

export type AuthTokens = {
  access_token: string;
  refresh_token: string;
};

export type AuthResponse = IUser & {
  tokens: AuthTokens;
};