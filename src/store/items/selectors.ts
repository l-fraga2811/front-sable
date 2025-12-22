import type { RootState } from "../index";

export const selectItems = (state: RootState) => state.items.items;
export const selectSelectedItem = (state: RootState) =>
  state.items.selectedItem;
export const selectItemsLoading = (state: RootState) => state.items.isLoading;
export const selectItemsError = (state: RootState) => state.items.error;
