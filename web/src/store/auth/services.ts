import api from "@/lib/axios";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
} from "@/types";

interface ApiError {
  response?: {
    data?: {
      error?: string;
    };
  };
}

export const authServices = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    try {
      const response = await api.post<AuthResponse>("/api/auth/signin", data);
      return response.data;
    } catch (error: unknown) {
      const apiError = error as ApiError;
      throw new Error(apiError.response?.data?.error || "Erro ao fazer login");
    }
  },

  register: async (
    data: RegisterRequest
  ): Promise<{ message: string; user: User }> => {
    try {
      const response = await api.post<{ message: string; user: User }>(
        "/api/auth/signup",
        data
      );
      return response.data;
    } catch (error: unknown) {
      const apiError = error as ApiError;
      throw new Error(apiError.response?.data?.error || "Erro ao registrar");
    }
  },

  logout: async (): Promise<{ message: string }> => {
    try {
      const response = await api.post<{ message: string }>("/api/auth/logout");
      return response.data;
    } catch (error: unknown) {
      const apiError = error as ApiError;
      throw new Error(apiError.response?.data?.error || "Erro ao fazer logout");
    }
  },

  getProfile: async (token?: string): Promise<User> => {
    const response = await api.get<User>("/api/auth/profile", {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    return response.data;
  },
};
