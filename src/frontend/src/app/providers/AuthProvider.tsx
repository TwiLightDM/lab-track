import { useEffect } from "react";
import type { ReactNode } from "react";

import { Context, store } from "./AuthContext";

type AuthProviderProps = {
  children: ReactNode;
};

export const AuthProvider = ({ children }: AuthProviderProps) => {
  useEffect(() => {
    const refreshToken = localStorage.getItem("refreshToken");

    if (refreshToken) {
      void store.checkAuth();
    }
  }, []);

  return <Context.Provider value={{ store }}>{children}</Context.Provider>;
};
