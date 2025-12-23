import api from "@/lib/axios";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
} from "@/types";

interface BackendAuthResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  refresh_token: string;
  user: {
    id: string;
    email: string;
  };
}

export const authServices = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<BackendAuthResponse>(
      "/api/auth/login",
      data
    );
    const authData = response.data;

    const expiresAt = new Date(
      Date.now() + authData.expires_in * 1000
    ).toISOString();

    return {
      message: "Login realizado com sucesso",
      token: authData.access_token,
      expiresAt: expiresAt,
      user: {
        id: authData.user.id,
        username: "",
        email: authData.user.email,
      },
    };
  },

  register: async (
    data: RegisterRequest
  ): Promise<{ message: string; user: User }> => {
    const response = await api.post<BackendAuthResponse>("/api/auth/register", {
      email: data.email,
      password: data.password,
      data: {
        username: data.username,
      },
    });
    const authData = response.data;

    return {
      message: "Usuário criado com sucesso",
      user: {
        id: authData.user.id,
        username: data.username,
        email: authData.user.email,
      },
    };
  },

  logout: async (): Promise<{ message: string }> => {
    // Client-side logout only (remove token)
    return { message: "Logout realizado com sucesso" };
  },

  getProfile: async (token?: string): Promise<User> => {
    const response = await api.get<User>("/api/auth/profile", {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    return response.data;
  },
};
