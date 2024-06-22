import nextClient from "./nextClient";

export type CreateTokenResponse = {
  success: boolean;
  data: any;
};

export const createToken = async (data: any) => {
  try {
    const response = await nextClient.post("create-token", data);
    return response;
  } catch (err) {
    console.log(err);
  }
};

export const fetchTokens = async () => {
  try {
    const response = await nextClient.post("fetch-tokens");
    return response;
  } catch (err) {
    console.log(err);
  }
};
