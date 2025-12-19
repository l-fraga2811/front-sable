import api from "@/lib/axios";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
} from "@/types";

export const authServices = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>("/auth/login", data);
    return response.data;
  },

  register: async (
    data: RegisterRequest
  ): Promise<{ message: string; user: User }> => {
    const response = await api.post<{ message: string; user: User }>(
      "/auth/register",
      data
    );
    return response.data;
  },

  logout: async (): Promise<{ message: string }> => {
    const response = await api.post<{ message: string }>("/api/auth/logout");
    return response.data;
  },

  getProfile: async (): Promise<User> => {
    const response = await api.get<User>("/api/auth/profile");
    return response.data;
  },
};
