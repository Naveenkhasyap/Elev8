import apiClient from "@/utils/apiClient";
import { NextApiRequest, NextApiResponse } from "next";
export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method === "POST") {
    let response = await apiClient.get(`/token/v1/fetch/all/0`);
    if (response.status === 200) {
      return res.status(200).json({
        ...response.data,
      });
    } else {
      res.status(404).json({ success: false });
    }
  } else {
    res.status(404).json({ status: false, data: {}, message: "not found" });
  }
}
