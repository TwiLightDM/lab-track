import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";

import LoginPage from "@/screens/LoginPage";
import ProfilePage from "@/screens/ProfilePage";

import ProtectedRoute from "./guards/ProtectedRoute";
import { paths } from "./paths";

const AppRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path={paths.login.route} element={<LoginPage />} />

        <Route
          path={paths.profile.route}
          element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/"
          element={<Navigate to={paths.profile.route} replace />}
        />

        <Route
          path="*"
          element={<Navigate to={paths.profile.route} replace />}
        />
      </Routes>
    </BrowserRouter>
  );
};

export default AppRouter;
