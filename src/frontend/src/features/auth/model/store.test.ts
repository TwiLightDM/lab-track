import { beforeEach, describe, expect, it, vi } from "vitest";

import Store from "./store";
import AuthService from "../api/AuthService";
import { userApi } from "@/entities/user/api";

import type { ILoginData } from "@/features/auth";
import type { IUser } from "@/entities/user";
import type { AuthResponse } from "../api/types";

vi.mock("../api/AuthService", () => ({
  default: {
    login: vi.fn(),
    refresh: vi.fn(),
  },
}));

vi.mock("@/entities/user/api", () => ({
  userApi: {
    getCurrentUser: vi.fn(),
  },
}));

const mockUser: IUser = {
  id: 1,
  name: "Иван",
  surname: "Иванов",
  email: "test@example.com",
  role: "admin",
  is_active: true,
};

const mockAuthResponse: AuthResponse = {
  ...mockUser,
  tokens: {
    access_token: "access-token",
    refresh_token: "refresh-token",
  },
};

describe("Store", () => {
  let store: Store;

  beforeEach(() => {
    localStorage.clear();
    vi.clearAllMocks();

    store = new Store();
  });

  describe("constructor", () => {
    it("создаёт Store с начальным состоянием", () => {
      expect(store.user).toBeNull();
      expect(store.isAuth).toBe(false);
      expect(store.isLoading).toBe(false);
    });

    it("устанавливает isAuth и isLoading, если accessToken существует", () => {
      localStorage.setItem("accessToken", "access-token");

      const newStore = new Store();

      expect(newStore.isAuth).toBe(true);
      expect(newStore.isLoading).toBe(true);
    });
  });

  describe("login", () => {
    const loginData: ILoginData = {
      email: "test@example.com",
      password: "password123",
    };

    it("успешно авторизует пользователя", async () => {
      vi.mocked(AuthService.login).mockResolvedValue({
        data: mockAuthResponse,
      } as never);

      await store.login(loginData);

      expect(AuthService.login).toHaveBeenCalledWith(loginData);

      expect(localStorage.getItem("accessToken")).toBe("access-token");
      expect(localStorage.getItem("refreshToken")).toBe("refresh-token");

      expect(store.user).toEqual(mockUser);
      expect(store.isAuth).toBe(true);
    });

    it("пробрасывает ошибку, если авторизация не удалась", async () => {
      const error = new Error("Неверный email или пароль");

      vi.mocked(AuthService.login).mockRejectedValue(error);

      await expect(store.login(loginData)).rejects.toThrow(
        "Неверный email или пароль",
      );

      expect(store.isAuth).toBe(false);
      expect(store.user).toBeNull();

      expect(localStorage.getItem("accessToken")).toBeNull();
      expect(localStorage.getItem("refreshToken")).toBeNull();
    });
  });

  describe("logout", () => {
    it("удаляет токены и сбрасывает состояние авторизации", () => {
      localStorage.setItem("accessToken", "access-token");
      localStorage.setItem("refreshToken", "refresh-token");

      store.setAuth(true);
      store.setUser(mockUser);
      store.setLoading(true);

      store.logout();

      expect(localStorage.getItem("accessToken")).toBeNull();
      expect(localStorage.getItem("refreshToken")).toBeNull();

      expect(store.user).toBeNull();
      expect(store.isAuth).toBe(false);
      expect(store.isLoading).toBe(false);
    });
  });

  describe("checkAuth", () => {
    it("восстанавливает авторизацию, если запрос пользователя успешен", async () => {
      vi.mocked(userApi.getCurrentUser).mockResolvedValue({
        data: mockUser,
      } as never);

      const promise = store.checkAuth();

      expect(store.isLoading).toBe(true);

      await promise;

      expect(userApi.getCurrentUser).toHaveBeenCalled();

      expect(store.user).toEqual(mockUser);
      expect(store.isAuth).toBe(true);
      expect(store.isLoading).toBe(false);
    });

    it("выполняет logout, если запрос пользователя завершился ошибкой", async () => {
      localStorage.setItem("accessToken", "access-token");
      localStorage.setItem("refreshToken", "refresh-token");

      store.setAuth(true);
      store.setUser(mockUser);

      vi.mocked(userApi.getCurrentUser).mockRejectedValue(
        new Error("Ошибка получения пользователя"),
      );

      await store.checkAuth();

      expect(store.user).toBeNull();
      expect(store.isAuth).toBe(false);
      expect(store.isLoading).toBe(false);

      expect(localStorage.getItem("accessToken")).toBeNull();
      expect(localStorage.getItem("refreshToken")).toBeNull();
    });
  });
});
