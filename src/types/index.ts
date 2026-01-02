export interface User {
  id: string;
  username: string;
  email: string;
  phone?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  phone?: string;
}

export interface AuthResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  refresh_token: string;
  user: User;
}

export interface Item {
  id: string;
  title: string;
  description: string;
  price: number;
  completed: boolean;
  createdAt: string;
  updatedAt: string;
  userId: string;
}

export interface CreateItemRequest {
  title: string;
  description?: string;
  price?: number;
}

export interface UpdateItemRequest {
  title?: string;
  description?: string;
  price?: number;
  completed?: boolean;
}

export interface ApiError {
  error: string;
}
