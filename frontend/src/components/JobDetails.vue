<template>
  <div class="job-details">
    <!-- Header -->
    <div class="header">
      <h1>{{ job.Name }}</h1>
    </div>
    <div>
        <p>Status: <span :class=job.Status>{{ job.Status }}</span></p>
    </div>

    <div class="content">
      <!-- Main Panel -->
      <div class="main-panel">
        <section class="build-summary">
          <h2>Last Build</h2>
          <p>Build Number: {{ job && job.LastBuild && job.LastBuild.number }}</p>
          <p>Duration: {{ job && job.LastBuild && job.LastBuild.duration }} seconds</p>
          <p>Status: <span :class="job.Status">{{ job && job.LastBuild && job.LastBuild.status }}</span></p>
        </section>
      </div>
      <LogViewer ref="childRef"/>
    </div>
  </div>
</template>

<script>
import apiService from "../services/jobService";
import LogViewer from "./LogViewer.vue";

export default {
  name: "JobDetails",
  components: {
    LogViewer,
  },
  data() {
    return {
      job: {
        Name: "",
        Status: "SUCCESS",
        LastBuild: {
          number: 123,
          duration: 45,
          status: "SUCCESS",
        },
        Builds: [
          { number: 123, status: "SUCCESS", duration: 45 },
          { number: 122, status: "FAILED", duration: 30 },
          { number: 121, status: "RUNNING", duration: null },
        ],
      },
      showBuildHistory: false,
    };
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    async fetchData() {
        try {
            const response = await apiService.getJobs(this.$route.fullPath);
            const jobs = response.data;
            
            const toFolder = !(jobs.length == 1 && !jobs[0].IsDir);
            this.$emit('navigate', toFolder);

            this.job = jobs && jobs[0];
        } catch (error) {
            this.error = "Error fetching jobs";
            console.error(error);
        }
    },
    handleStatusChange(event) {
      this.$refs.childRef.handleStatusChange(event);
    },
  },
};
</script>

<style scoped>
/* General Styles */
.job-details {
  font-family: Arial, sans-serif;
}

/* Header Styles */
.header {
  display: flex;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  background-color: #f3f4f6;
  padding: 20px;
  border-bottom: 1px solid #ddd;
}
.header h1 {
  margin: 0;
  font-size: 24px;
}
.header span.success {
  color: green;
}
.header span.failure {
  color: red;
}
.header span.running {
  color: orange;
}

.header span.pending {
  color: blueviolet;
}

/* Layout Styles */
.content {
  display: flex;
}

/* Side Panel */
.side-panel {
  width: 250px;
  background-color: #f9fafb;
  border-right: 1px solid #ddd;
  padding: 15px;
}
.side-panel ul {
  list-style: none;
  padding: 0;
}
.side-panel ul li {
  margin-bottom: 10px;
}
.side-panel ul li a {
  color: #007bff;
  text-decoration: none;
}

/* Collapsible Build History */
.build-history {
  margin-top: 20px;
}
.build-history h3 {
  cursor: pointer;
  margin: 0;
  color: #333;
  font-weight: bold;
  display: flex;
  justify-content: space-between;
}
.build-history ul {
  list-style: none;
  padding: 0;
  margin-top: 10px;
}
.build-history ul li {
  margin-bottom: 8px;
}
.build-history ul li span.success {
  color: green;
}
.build-history ul li span.failure {
  color: red;
}
.build-history ul li span.running {
  color: orange;
}

/* Main Panel */
.main-panel {
  flex: 1;
  padding: 50px;
}
</style>
