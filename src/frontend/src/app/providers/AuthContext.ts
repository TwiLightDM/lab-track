import { createContext } from "react";
import Store from "@/features/auth/model/store";

type State = {
  store: Store;
};

export const store = new Store();

export const Context = createContext<State>({ store });
