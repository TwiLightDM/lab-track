import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, useLocation } from "react-router-dom";

import ProtectedRoute from "./ProtectedRoute";
import { useAuth } from "@/entities/user/hook/useAuth";
import { paths } from "@/app/router/paths";

vi.mock("@/entities/user/hook/useAuth", () => ({
  useAuth: vi.fn(),
}));

const LocationDisplay = () => {
  const location = useLocation();

  return <div data-testid="location">{location.pathname}</div>;
};

const renderProtectedRoute = () =>
  render(
    <MemoryRouter initialEntries={["/protected"]}>
      <LocationDisplay />
      <ProtectedRoute>
        <div>Защищённая страница</div>
      </ProtectedRoute>
    </MemoryRouter>,
  );

describe("ProtectedRoute", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("показывает загрузку, пока проверяется авторизация", () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      isLoading: true,
      isAuthenticated: false,
      refetchUser: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.getByText("Проверяем авторизацию...")).toBeInTheDocument();
    expect(screen.queryByText("Защищённая страница")).not.toBeInTheDocument();
  });

  it("перенаправляет на страницу входа, если пользователь не авторизован", () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      isLoading: false,
      isAuthenticated: false,
      refetchUser: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.getByTestId("location")).toHaveTextContent(paths.login.route);
    expect(screen.queryByText("Защищённая страница")).not.toBeInTheDocument();
  });

  it("показывает содержимое, если пользователь авторизован", () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      isLoading: false,
      isAuthenticated: true,
      refetchUser: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.getByText("Защищённая страница")).toBeInTheDocument();
    expect(screen.getByTestId("location")).toHaveTextContent("/protected");
  });
});
