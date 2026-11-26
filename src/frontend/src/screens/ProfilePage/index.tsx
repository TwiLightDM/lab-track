import { observer } from "mobx-react-lite";
import { Navigate } from "react-router-dom";

import { useAuth } from "@/entities/user/hook/useAuth";
import { paths } from "@/app/router/paths";

import AdminProfile from "./ui/AdminProfile";
import LabAssistantProfile from "./ui/LabAssistantProfile";
import LabManagerProfile from "./ui/LabManagerProfile";
import RepairmanProfile from "./ui/RepairmanProfile";

const ProfilePage = observer(() => {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to={paths.login.route} replace />;
  }

  switch (user.role) {
    case "admin":
      return <AdminProfile />;

    case "lab_assistant":
      return <LabAssistantProfile />;

    case "lab_manager":
      return <LabManagerProfile />;

    case "repairman":
      return <RepairmanProfile />;

    default:
      return <Navigate to={paths.login.route} replace />;
  }
});

export default ProfilePage;
