import { Context } from "@/app/providers/AuthContext";
import type { ILoginData } from "@/features/auth";
import LoginForm from "@/features/auth/ui/LoginForm";
import { useContext } from "react";

const LoginPage = () => {
  const { store } = useContext(Context);

  const handleLogin = async (data: ILoginData) => {
    await store.login(data);
  };

  return (
    <main>
      <LoginForm onSubmit={handleLogin} />
    </main>
  );
};

export default LoginPage;
