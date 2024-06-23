import axios from "axios";
const baseUrl = process.env.BASE_URL;
const apiClient = axios.create({
  baseURL: `${baseUrl}`,
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.response.use(
  function (response) {
    return response;
  },
  function (error) {
    let res = error.response;
    console.error("Looks like there was a problem. Status Code: " + res.status);
    return error;
  }
);

export default apiClient;
