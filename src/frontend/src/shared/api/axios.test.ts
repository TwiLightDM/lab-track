import axios, {
  AxiosError,
  type AxiosAdapter,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { $api, API_URL } from "./axios";

type RefreshTokenResponse = {
  tokens: {
    access_token: string;
    refresh_token: string;
  };
};

const createResponse = (
  config: InternalAxiosRequestConfig,
  data: unknown = { success: true },
): AxiosResponse => ({
  data,
  status: 201,
  statusText: "OK",
  headers: {},
  config,
});

const createUnauthorizedError = (config: InternalAxiosRequestConfig) => {
  const response: AxiosResponse = {
    data: { message: "Unauthorized" },
    status: 401,
    statusText: "Unauthorized",
    headers: {},
    config,
  };

  return new AxiosError(
    "Unauthorized",
    "ERR_BAD_REQUEST",
    config,
    undefined,
    response,
  );
};

const mockRefreshResponse = (
  accessToken = "new-access-token",
  refreshToken = "new-refresh-token",
) => {
  const response: AxiosResponse<RefreshTokenResponse> = {
    data: {
      tokens: {
        access_token: accessToken,
        refresh_token: refreshToken,
      },
    },
    status: 200,
    statusText: "OK",
    headers: {},
    config: {} as InternalAxiosRequestConfig,
  };

  return vi.spyOn(axios, "post").mockResolvedValue(response);
};

describe("Axios instance", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("добавляет accessToken в заголовок Authorization", async () => {
    localStorage.setItem("accessToken", "my-access-token");

    let requestConfig: InternalAxiosRequestConfig | undefined;

    $api.defaults.adapter = (async (config) => {
      requestConfig = config;
      return createResponse(config);
    }) as AxiosAdapter;

    await $api.get("/users");

    expect(requestConfig?.headers.Authorization).toBe("Bearer my-access-token");
  });

  it("не перезаписывает уже установленный Authorization", async () => {
    localStorage.setItem("accessToken", "my-access-token");

    let requestConfig: InternalAxiosRequestConfig | undefined;

    $api.defaults.adapter = (async (config) => {
      requestConfig = config;
      return createResponse(config);
    }) as AxiosAdapter;

    await $api.get("/users", {
      headers: {
        Authorization: "Bearer custom-token",
      },
    });

    expect(requestConfig?.headers.Authorization).toBe("Bearer custom-token");
  });

  it("обновляет токены и повторяет запрос после получения 401", async () => {
    localStorage.setItem("accessToken", "old-access-token");
    localStorage.setItem("refreshToken", "old-refresh-token");

    const refreshSpy = mockRefreshResponse();

    const requestAttempts: string[] = [];

    $api.defaults.adapter = (async (config) => {
      requestAttempts.push(config.headers.Authorization as string);

      if (requestAttempts.length === 1) {
        throw createUnauthorizedError(config);
      }

      return createResponse(config);
    }) as AxiosAdapter;

    const response = await $api.get("/users");

    expect(refreshSpy).toHaveBeenCalledTimes(1);
    expect(refreshSpy).toHaveBeenCalledWith(`${API_URL}/public/refresh`, null, {
      headers: {
        Authorization: "Bearer old-refresh-token",
      },
    });

    expect(localStorage.getItem("accessToken")).toBe("new-access-token");
    expect(localStorage.getItem("refreshToken")).toBe("new-refresh-token");

    expect(requestAttempts).toEqual([
      "Bearer old-access-token",
      "Bearer new-access-token",
    ]);

    expect(response.status).toBe(201);
  });

  it("не пытается обновлять токены при ошибке авторизации на login", async () => {
    const refreshSpy = vi.spyOn(axios, "post");

    $api.defaults.adapter = (async (config) => {
      throw createUnauthorizedError(config);
    }) as AxiosAdapter;

    await expect(
      $api.post("/public/login", {
        email: "test@example.com",
        password: "password",
      }),
    ).rejects.toMatchObject({
      response: {
        status: 401,
      },
    });

    expect(refreshSpy).not.toHaveBeenCalled();
  });

  it("не пытается обновлять токены повторно при ошибке refresh", async () => {
    const refreshSpy = vi.spyOn(axios, "post");

    $api.defaults.adapter = (async (config) => {
      throw createUnauthorizedError(config);
    }) as AxiosAdapter;

    await expect($api.post("/public/refresh", null)).rejects.toMatchObject({
      response: {
        status: 401,
      },
    });

    expect(refreshSpy).not.toHaveBeenCalled();
  });

  it("обновляет токены только один раз при нескольких одновременных 401", async () => {
    localStorage.setItem("accessToken", "old-access-token");
    localStorage.setItem("refreshToken", "old-refresh-token");

    let resolveRefresh!: (value: AxiosResponse<RefreshTokenResponse>) => void;

    const refreshPromise = new Promise<AxiosResponse<RefreshTokenResponse>>(
      (resolve) => {
        resolveRefresh = resolve;
      },
    );

    const refreshSpy = vi
      .spyOn(axios, "post")
      .mockImplementation(() => refreshPromise);

    const attempts = new Map<string, number>();

    $api.defaults.adapter = (async (config) => {
      const url = config.url ?? "";
      const count = (attempts.get(url) ?? 0) + 1;

      attempts.set(url, count);

      if (count === 1) {
        throw createUnauthorizedError(config);
      }

      return createResponse(config);
    }) as AxiosAdapter;

    const firstRequest = $api.get("/users");
    const secondRequest = $api.get("/equipment");

    await vi.waitFor(() => {
      expect(refreshSpy).toHaveBeenCalledTimes(1);
    });

    resolveRefresh({
      data: {
        tokens: {
          access_token: "new-access-token",
          refresh_token: "new-refresh-token",
        },
      },
      status: 200,
      statusText: "OK",
      headers: {},
      config: {} as InternalAxiosRequestConfig,
    });

    const [firstResponse, secondResponse] = await Promise.all([
      firstRequest,
      secondRequest,
    ]);

    expect(firstResponse.status).toBe(201);
    expect(secondResponse.status).toBe(201);

    expect(refreshSpy).toHaveBeenCalledTimes(1);

    expect(attempts.get("/users")).toBe(2);
    expect(attempts.get("/equipment")).toBe(2);
  });
});
