import api from "@/lib/axios";
import type { Item, CreateItemRequest, UpdateItemRequest } from "@/types";

export const itemsServices = {
  getAll: async (): Promise<Item[]> => {
    const response = await api.get<Item[]>("/api/items");
    return response.data;
  },

  getById: async (id: string): Promise<Item> => {
    const response = await api.get<Item>(`/api/items/${id}`);
    return response.data;
  },

  create: async (data: CreateItemRequest): Promise<Item> => {
    const response = await api.post<Item>("/api/items", data);
    return response.data;
  },

  update: async (id: string, data: UpdateItemRequest): Promise<Item> => {
    const response = await api.put<Item>(`/api/items/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<{ message: string }> => {
    const response = await api.delete<{ message: string }>(`/api/items/${id}`);
    return response.data;
  },
};
