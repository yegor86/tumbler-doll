<template>
  <div class="job-list">
    <h2>Jobs</h2>
    <table>
    <thead>
        <tr>
        <th>Name</th>
        <th>Status</th>
        <th>Last Build</th>
        </tr>
    </thead>
    <tbody>
        <JobItem v-for="job in jobs" :key="job.id" :job="job"/>
    </tbody>
    </table>
  </div>
</template>
  
<script>
import JobItem from "../components/JobItem.vue";
import apiService from "../services/jobService";

export default {
    name: "JobList",
    components: { 
      JobItem
    },
    data() {
      return {
        jobs: [],
        error: null,
      }
    },
    mounted() {
        this.fetchData();
    },
    methods: {
        async fetchData() {
            try {
                const response = await apiService.getJobs(this.$route.fullPath);
                this.jobs = response.data;
            } catch (error) {
                this.error = "Error fetching jobs";
                console.error(error);
            }
        }
    },
};
</script>
  
<style scoped>
  .job-list {
    flex-grow: 1;
    padding: 20px;
  }
  ul {
    list-style: none;
    padding-left: 20px;
  }
  table {
    width: 100%;
    border-collapse: collapse;
  }
  table th, table td {
    border: 1px solid #ddd;
    padding: 8px;
  }
  table th {
    background-color: #f4f4f4;
  }
</style>
  