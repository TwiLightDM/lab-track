import { useAuth } from "@/entities/user/hook/useAuth";
import type { UserRole } from "@/entities/user";
import { observer } from "mobx-react-lite";
import { Navigate } from "react-router-dom";
import { paths } from "@/app/router/paths";

interface RoleGuardProps {
  children: React.ReactNode;
  allowRoles: UserRole[];
}

const RoleGuard = observer(({ children, allowRoles }: RoleGuardProps) => {
  const { user, isLoading, isAuthenticated } = useAuth();

  if (isLoading) {
    return <div>Проверяем авторизацию...</div>;
  }

  if (!isAuthenticated || !user) {
    return <Navigate to={paths.login.route} replace />;
  }

  if (!allowRoles.includes(user.role)) {
    return <Navigate to={paths.home.route} replace />;
  }

  return <>{children}</>;
});

export default RoleGuard;
