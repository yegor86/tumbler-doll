<template>
    <div>
      <h2>Live Logs</h2>
      <pre>{{ logs }}</pre>
    </div>
  </template>
  
  <script>
  import apiService from "../services/jobService";

  export default {
    data() {
      return {
        logs: ""
      };
    },
    mounted() {
      const eventSource = apiService.streamJobExec(this.$route.fullPath);
  
      eventSource.onmessage = (event) => {
        this.logs += event.data + "\n";
      };
  
      eventSource.onerror = (error) => {
        console.error("EventSource failed:", error);
        eventSource.close();
      };
    }
  };
  </script>
  