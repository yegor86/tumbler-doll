<template>
    <div class="console">
      <h2>Live Logs</h2>
      <pre>{{ logs }}</pre>
    </div>
  </template>
  
  <script>
  import { AnsiUp } from 'ansi_up'
  import apiService from "../services/jobService";

  export default {
    data() {
      return {
        ansi: undefined,
        logs: '',
      };
    },
    beforeMount () {
      this.ansi = new AnsiUp()
    },
    updated () {
      // auto-scroll to the bottom when the DOM is updated
      this.$el.scrollTop = this.$el.scrollHeight
    },
    methods: {
      handleStatusChange(event) {
        
        const eventSource = apiService.streamJobExec(this.$route.fullPath, event.WorkflowID);
        eventSource.onmessage = (event) => {
          this.logs += this.ansi.ansi_to_html(event.data).replace(/\n/gm, '<br>') + "\n";
        };
    
        eventSource.onerror = (error) => {
          console.error("EventSource failed:", error);
          eventSource.close();
        };
      },
    },
  };
  </script>

<style scoped>
.console {
  font-family: monospace;
  text-align: left;
  background-color: black;
  color: #fff;
  overflow-y: auto;
}
</style>
  