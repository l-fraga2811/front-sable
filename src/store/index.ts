import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./auth/reducers";
import itemsReducer from "./items/reducers";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    items: itemsReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
