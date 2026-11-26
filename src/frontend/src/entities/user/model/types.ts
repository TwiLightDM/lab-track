export type UserRole = "admin" | "lab_assistant" | "lab_manager" | "repairman";

export type IUser = {
  id: number;
  name: string;
  surname: string;
  email: string;
  role: UserRole;
  is_active: boolean;
};
