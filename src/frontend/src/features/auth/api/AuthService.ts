import { $api } from "@/shared/api/axios";
import type { AxiosResponse } from "axios";
import type { AuthResponse, ILoginData } from "./types";

export default class AuthService {
  static async login(data: ILoginData): Promise<AxiosResponse<AuthResponse>> {
    return $api.post<AuthResponse>("/public/login", data);
  }

  static async refresh(
    refreshToken: string,
  ): Promise<AxiosResponse<AuthResponse>> {
    return $api.post<AuthResponse>("/public/refresh", null, {
      headers: {
        Authorization: `Bearer ${refreshToken}`,
      },
    });
  }
}
