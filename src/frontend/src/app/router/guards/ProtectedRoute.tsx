import { useAuth } from "@/entities/user/hook/useAuth";
import { observer } from "mobx-react-lite";
import { Navigate } from "react-router-dom";
import { paths } from "@/app/router/paths";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute = observer(({ children }: ProtectedRouteProps) => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <div>Проверяем авторизацию...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to={paths.login.route} replace />;
  }

  return <>{children}</>;
});

export default ProtectedRoute;
