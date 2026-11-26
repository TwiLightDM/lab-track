import { makeAutoObservable } from "mobx";
import type { ILoginData } from "@/features/auth";
import { AxiosError } from "axios";
import type { IUser } from "@/entities/user";
import { userApi } from "@/entities/user/api";
import AuthService from "../api/AuthService";

export default class Store {
  user: IUser | null = null;
  isAuth = false;
  isLoading = false;

  constructor() {
    const token = localStorage.getItem("accessToken");

    this.isAuth = Boolean(token);
    this.isLoading = Boolean(token);

    makeAutoObservable(this);
  }

  setAuth(bool: boolean) {
    this.isAuth = bool;
  }

  setUser(user: IUser | null) {
    this.user = user;
  }

  setLoading(bool: boolean) {
    this.isLoading = bool;
  }

  async login(data: ILoginData) {
    try {
      const response = await AuthService.login(data);

      const { tokens, ...user } = response.data;

      localStorage.setItem("accessToken", tokens.access_token);
      localStorage.setItem("refreshToken", tokens.refresh_token);

      this.setAuth(true);
      this.setUser(user);
    } catch (e) {
      if (e instanceof AxiosError) {
        console.log("Login error:", e.response?.data?.message);
      } else if (e instanceof Error) {
        console.log(e.message);
      }

      throw e;
    }
  }

  logout() {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");

    this.setAuth(false);
    this.setUser(null);
    this.setLoading(false);
  }

  async checkAuth() {
    this.setLoading(true);

    try {
      const response = await userApi.getCurrentUser();

      this.setAuth(true);
      this.setUser(response.data);
    } catch {
      this.logout();
    } finally {
      this.setLoading(false);
    }
  }
}
