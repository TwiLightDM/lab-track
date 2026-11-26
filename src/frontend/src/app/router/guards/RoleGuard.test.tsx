import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, useLocation } from "react-router-dom";

import RoleGuard from "./RoleGuard";
import { useAuth } from "@/entities/user/hook/useAuth";
import type { IUser, UserRole } from "@/entities/user";
import { paths } from "@/app/router/paths";

vi.mock("@/entities/user/hook/useAuth", () => ({
  useAuth: vi.fn(),
}));

const LocationDisplay = () => {
  const location = useLocation();

  return <div data-testid="location">{location.pathname}</div>;
};

const createUser = (role: UserRole): IUser => ({
  id: 1,
  name: "Иван",
  surname: "Иванов",
  email: "ivan@example.com",
  role,
  is_active: true,
});

const renderRoleGuard = (allowRoles: UserRole[]) =>
  render(
    <MemoryRouter initialEntries={["/protected"]}>
      <LocationDisplay />
      <RoleGuard allowRoles={allowRoles}>
        <div>Страница с ограниченным доступом</div>
      </RoleGuard>
    </MemoryRouter>,
  );

describe("RoleGuard", () => {
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

    renderRoleGuard(["admin"]);

    expect(screen.getByText("Проверяем авторизацию...")).toBeInTheDocument();
    expect(
      screen.queryByText("Страница с ограниченным доступом"),
    ).not.toBeInTheDocument();
  });

  it("перенаправляет на вход, если пользователь не авторизован", () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      isLoading: false,
      isAuthenticated: false,
      refetchUser: vi.fn(),
    });

    renderRoleGuard(["admin"]);

    expect(screen.getByTestId("location")).toHaveTextContent(paths.login.route);
  });

  it("перенаправляет на главную, если у пользователя нет нужной роли", () => {
    vi.mocked(useAuth).mockReturnValue({
      user: createUser("lab_assistant"),
      isLoading: false,
      isAuthenticated: true,
      refetchUser: vi.fn(),
    });

    renderRoleGuard(["admin"]);

    expect(screen.getByTestId("location")).toHaveTextContent(paths.home.route);
    expect(
      screen.queryByText("Страница с ограниченным доступом"),
    ).not.toBeInTheDocument();
  });

  it("показывает содержимое, если роль пользователя разрешена", () => {
    vi.mocked(useAuth).mockReturnValue({
      user: createUser("admin"),
      isLoading: false,
      isAuthenticated: true,
      refetchUser: vi.fn(),
    });

    renderRoleGuard(["admin", "lab_manager"]);

    expect(
      screen.getByText("Страница с ограниченным доступом"),
    ).toBeInTheDocument();
    expect(screen.getByTestId("location")).toHaveTextContent("/protected");
  });
});
