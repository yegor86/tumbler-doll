import axios from "axios";

const apiClient = axios.create({
  baseURL: "http://localhost:8888",
  headers: {
    "Content-Type": "application/json",
  },
});

export default {
  getJobs(jobPath) {
    if (!jobPath || jobPath == "/") {
      return apiClient.get("/jobs");  
    }
    return apiClient.get(jobPath);
  },

  submitJob(jobPath) {
    return apiClient.post(`/submit/${jobPath}`);
  },

  streamJobExec(jobPath, workflowId) {
    const eventSource = new EventSource(`http://localhost:8888/stream/${jobPath}?workflowId=${workflowId}`);
    return eventSource;
  },
};