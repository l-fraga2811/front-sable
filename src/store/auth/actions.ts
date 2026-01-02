import { createAsyncThunk } from "@reduxjs/toolkit";
import type { LoginRequest, RegisterRequest } from "@/types";
import { authServices } from "./services";

export const login = createAsyncThunk(
  "auth/login",
  async (data: LoginRequest, { rejectWithValue }) => {
    try {
      const response = await authServices.login(data);
      localStorage.setItem("token", response.access_token);

      const profile = await authServices.getProfile(response.access_token);
      localStorage.setItem("user", JSON.stringify(profile));

      return {
        ...response,
        user: profile,
      };
    } catch (error: unknown) {
      const err = error as { message?: string };
      if (err.message?.includes("email_not_confirmed")) {
        return rejectWithValue(
          "Por favor, confirme seu e-mail antes de fazer login. Verifique sua caixa de entrada."
        );
      }
      return rejectWithValue(err.message || "Erro ao fazer login");
    }
  }
);

export const register = createAsyncThunk(
  "auth/register",
  async (data: RegisterRequest, { rejectWithValue }) => {
    try {
      const response = await authServices.register(data);
      return response;
    } catch (error: unknown) {
      const err = error as { message?: string };
      return rejectWithValue(err.message || "Erro ao registrar");
    }
  }
);

export const logout = createAsyncThunk(
  "auth/logout",
  async (_, { rejectWithValue }) => {
    try {
      await authServices.logout();
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      return null;
    } catch (error: unknown) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      const err = error as { message?: string };
      return rejectWithValue(err.message || "Erro ao fazer logout");
    }
  }
);

export const getProfile = createAsyncThunk(
  "auth/getProfile",
  async (_, { rejectWithValue }) => {
    try {
      const token =
        typeof window !== "undefined" ? localStorage.getItem("token") : null;
      if (!token) {
        return rejectWithValue("No authentication token found");
      }
      const response = await authServices.getProfile(token);
      localStorage.setItem("user", JSON.stringify(response));
      return response;
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      return rejectWithValue(
        err.response?.data?.error || "Erro ao buscar perfil"
      );
    }
  }
);
