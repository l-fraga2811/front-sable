import { createAsyncThunk } from "@reduxjs/toolkit";
import type { CreateItemRequest, UpdateItemRequest } from "@/types";
import { itemsServices } from "./services";

export const fetchItems = createAsyncThunk(
  "items/fetchAll",
  async (_, { rejectWithValue }) => {
    try {
      const response = await itemsServices.getAll();
      return response;
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      return rejectWithValue(
        err.response?.data?.error || "Erro ao buscar items"
      );
    }
  }
);

export const fetchItemById = createAsyncThunk(
  "items/fetchById",
  async (id: string, { rejectWithValue }) => {
    try {
      const response = await itemsServices.getById(id);
      return response;
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      return rejectWithValue(
        err.response?.data?.error || "Erro ao buscar item"
      );
    }
  }
);

export const createItem = createAsyncThunk(
  "items/create",
  async (data: CreateItemRequest, { rejectWithValue }) => {
    try {
      const response = await itemsServices.create(data);
      return response;
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      return rejectWithValue(err.response?.data?.error || "Erro ao criar item");
    }
  }
);

export const updateItem = createAsyncThunk(
  "items/update",
  async (
    { id, data }: { id: string; data: UpdateItemRequest },
    { rejectWithValue }
  ) => {
    try {
      const response = await itemsServices.update(id, data);
      return response;
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      return rejectWithValue(
        err.response?.data?.error || "Erro ao atualizar item"
      );
    }
  }
);

export const deleteItem = createAsyncThunk(
  "items/delete",
  async (id: string, { rejectWithValue }) => {
    try {
      await itemsServices.delete(id);
      return id;
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      return rejectWithValue(
        err.response?.data?.error || "Erro ao deletar item"
      );
    }
  }
);
