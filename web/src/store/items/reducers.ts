import { createSlice } from "@reduxjs/toolkit";
import type { Item } from "@/types";
import {
  fetchItems,
  fetchItemById,
  createItem,
  updateItem,
  deleteItem,
} from "./actions";

interface ItemsState {
  items: Item[];
  selectedItem: Item | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: ItemsState = {
  items: [],
  selectedItem: null,
  isLoading: false,
  error: null,
};

const itemsSlice = createSlice({
  name: "items",
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    clearSelectedItem: (state) => {
      state.selectedItem = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchItems.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchItems.fulfilled, (state, action) => {
        state.isLoading = false;
        state.items = action.payload;
      })
      .addCase(fetchItems.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      .addCase(fetchItemById.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchItemById.fulfilled, (state, action) => {
        state.isLoading = false;
        state.selectedItem = action.payload;
      })
      .addCase(fetchItemById.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      .addCase(createItem.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createItem.fulfilled, (state, action) => {
        state.isLoading = false;
        state.items.push(action.payload);
      })
      .addCase(createItem.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      .addCase(updateItem.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(updateItem.fulfilled, (state, action) => {
        state.isLoading = false;
        const index = state.items.findIndex(
          (item) => item.id === action.payload.id
        );
        if (index !== -1) {
          state.items[index] = action.payload;
        }
        if (state.selectedItem?.id === action.payload.id) {
          state.selectedItem = action.payload;
        }
      })
      .addCase(updateItem.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      .addCase(deleteItem.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(deleteItem.fulfilled, (state, action) => {
        state.isLoading = false;
        state.items = state.items.filter((item) => item.id !== action.payload);
        if (state.selectedItem?.id === action.payload) {
          state.selectedItem = null;
        }
      })
      .addCase(deleteItem.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });
  },
});

export const { clearError, clearSelectedItem } = itemsSlice.actions;
export default itemsSlice.reducer;
