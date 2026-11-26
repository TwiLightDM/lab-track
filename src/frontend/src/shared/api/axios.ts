import axios from "axios";
import type { InternalAxiosRequestConfig } from "axios";

type RefreshTokenResponse = {
  tokens: {
    access_token: string;
    refresh_token: string;
  };
};

type RetryRequestConfig = InternalAxiosRequestConfig & {
  _isRetry?: boolean;
};

type FailedRequest = {
  resolve: (token: string) => void;
  reject: (error: unknown) => void;
};

export const API_URL = "http://localhost:9080/api/v1";

export const $api = axios.create({
  baseURL: API_URL,
});

let isRefreshing = false;

let failedQueue: FailedRequest[] = [];

const processQueue = (error: unknown | null, token?: string) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else if (token) {
      resolve(token);
    }
  });

  failedQueue = [];
};

const clearTokens = () => {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
};

$api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem("accessToken");

  if (token && !config.headers.Authorization) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

$api.interceptors.response.use(
  (response) => response,

  async (error) => {
    const originalRequest = error.config as RetryRequestConfig | undefined;

    if (
      error.response?.status !== 401 ||
      !originalRequest ||
      originalRequest._isRetry ||
      originalRequest.url?.includes("/public/login") ||
      originalRequest.url?.includes("/public/refresh")
    ) {
      return Promise.reject(error);
    }

    if (isRefreshing) {
      return new Promise<string>((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      }).then((token) => {
        originalRequest.headers.Authorization = `Bearer ${token}`;
        return $api(originalRequest);
      });
    }

    originalRequest._isRetry = true;
    isRefreshing = true;

    const refreshToken = localStorage.getItem("refreshToken");

    if (!refreshToken) {
      clearTokens();
      processQueue(error);
      isRefreshing = false;

      window.location.href = "/login";

      return Promise.reject(error);
    }

    try {
      const response = await axios.post<RefreshTokenResponse>(
        `${API_URL}/public/refresh`,
        null,
        {
          headers: {
            Authorization: `Bearer ${refreshToken}`,
          },
        },
      );

      const { access_token: newAccessToken, refresh_token: newRefreshToken } =
        response.data.tokens;

      localStorage.setItem("accessToken", newAccessToken);
      localStorage.setItem("refreshToken", newRefreshToken);

      processQueue(null, newAccessToken);

      originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;

      return $api(originalRequest);
    } catch (refreshError) {
      processQueue(refreshError);
      clearTokens();

      window.location.href = "/login";

      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  },
);
