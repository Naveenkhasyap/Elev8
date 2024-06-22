import axios from "axios";
const baseUrl = process.env.NEXT_URL;

const nextClient = axios.create({
  baseURL: baseUrl,
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
  },
});

nextClient.interceptors.response.use(
  function (response) {
    return response;
  },
  function (error) {
    console.log(error, "from axios intercepter");
    let res = error.response;
    return res;
  }
);

export default nextClient;
