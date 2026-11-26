import { beforeEach, describe, expect, it, vi } from "vitest";
import type { AxiosResponse } from "axios";

import { $api } from "@/shared/api/axios";
import AuthService from "./AuthService";
import type { AuthResponse, ILoginData } from "./types";

vi.mock("@/shared/api/axios", () => ({
  $api: {
    post: vi.fn(),
  },
}));

describe("AuthService", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("login", () => {
    it("отправляет POST-запрос на /public/login с данными пользователя", async () => {
      const loginData: ILoginData = {
        email: "test@example.com",
        password: "password123",
      };

      const response = {
        data: {
          id: 1,
          name: "Иван",
          surname: "Иванов",
          email: "test@example.com",
          role: "admin",
          is_active: true,
          tokens: {
            access_token: "access-token",
            refresh_token: "refresh-token",
          },
        },
      } as AxiosResponse<AuthResponse>;

      vi.mocked($api.post).mockResolvedValue(response);

      const result = await AuthService.login(loginData);

      expect($api.post).toHaveBeenCalledWith("/public/login", loginData);
      expect(result).toBe(response);
    });
  });

  describe("refresh", () => {
    it("отправляет POST-запрос на /public/refresh с refresh-токеном в заголовке", async () => {
      const refreshToken = "refresh-token";

      const response = {
        data: {
          id: 1,
          name: "Иван",
          surname: "Иванов",
          email: "test@example.com",
          role: "admin",
          is_active: true,
          tokens: {
            access_token: "new-access-token",
            refresh_token: "new-refresh-token",
          },
        },
      } as AxiosResponse<AuthResponse>;

      vi.mocked($api.post).mockResolvedValue(response);

      const result = await AuthService.refresh(refreshToken);

      expect($api.post).toHaveBeenCalledWith("/public/refresh", null, {
        headers: {
          Authorization: `Bearer ${refreshToken}`,
        },
      });

      expect(result).toBe(response);
    });
  });
});
