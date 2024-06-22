import { createSlice } from "@reduxjs/toolkit";

export interface initialStateI {
  isDarkMode: boolean;
}
const initialState = {
  modelOpen: true,
};

export const isMainSlice = createSlice({
  name: "MainSlice",
  initialState,
  reducers: {
    setIsModelOpen: (state, action) => {
      state.modelOpen = action.payload;
    },
  },
});
export const { setIsModelOpen } = isMainSlice.actions;
export default isMainSlice;
