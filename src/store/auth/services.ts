import api from "@/lib/axios";
import { getSupabaseClient } from "@/lib/supabaseClient";
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
} from "@/types";

export const authServices = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const supabase = getSupabaseClient();
    const { data: authData, error } = await supabase.auth.signInWithPassword({
      email: data.email,
      password: data.password,
    });
    if (error || !authData.session) {
      throw new Error(error?.message || "Erro ao fazer login");
    }

    return {
      message: "Login realizado com sucesso",
      token: authData.session.access_token,
      expiresAt: authData.session.expires_at
        ? new Date(authData.session.expires_at * 1000).toISOString()
        : "",
      user: {
        id: authData.user?.id || "",
        username: "",
        email: authData.user?.email || data.email,
      },
    };
  },

  register: async (
    data: RegisterRequest
  ): Promise<{ message: string; user: User }> => {
    const supabase = getSupabaseClient();
    const { data: authData, error } = await supabase.auth.signUp({
      email: data.email,
      password: data.password,
      options: {
        data: {
          username: data.username,
        },
      },
    });
    if (error || !authData.user) {
      throw new Error(error?.message || "Erro ao registrar");
    }

    return {
      message: "Usuário criado com sucesso",
      user: {
        id: authData.user.id,
        username: data.username,
        email: authData.user.email || data.email,
      },
    };
  },

  logout: async (): Promise<{ message: string }> => {
    const supabase = getSupabaseClient();
    const { error } = await supabase.auth.signOut();
    if (error) {
      throw new Error(error.message);
    }
    return { message: "Logout realizado com sucesso" };
  },

  getProfile: async (token?: string): Promise<User> => {
    const response = await api.get<User>("/api/auth/profile", {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    return response.data;
  },
};
