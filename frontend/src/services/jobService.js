import axios from "axios";

const apiClient = axios.create({
  baseURL: "http://localhost:8888",
  headers: {
    "Content-Type": "application/json",
  },
});

export default {
  getJobs() {
    return apiClient.get("/jobs");
  },
};