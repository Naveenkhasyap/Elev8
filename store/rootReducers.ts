// store/rootReducer.ts
import { combineReducers } from "@reduxjs/toolkit";
// Import your reducers here
import exampleReducer from "./slices/walletSlice";

const rootReducer = combineReducers({
  // Add your reducers here
  example: exampleReducer,
});

export default rootReducer;
