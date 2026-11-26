import { useContext } from "react";
import type { IUser } from "@/entities/user";
import { Context } from "@/app/providers/AuthContext";

export type AuthContextType = {
  user: IUser | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  refetchUser: () => Promise<void>;
};

export const useAuth = (): AuthContextType => {
  const { store } = useContext(Context);

  return {
    user: store.user,
    isLoading: store.isLoading,
    isAuthenticated: store.isAuth,
    refetchUser: async () => {
      await store.checkAuth();
    },
  };
};
