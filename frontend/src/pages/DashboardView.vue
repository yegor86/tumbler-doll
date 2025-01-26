<template>
<div class="dashboard">
    <Navbar />
    <div class="main-content">
    <Sidebar />
    <transition name="fade" mode="out-in">
      <component @sendJobType="handleJobType"
        :is="currentComponent"
        :key="componentKey"
      />
    </transition>
    <!-- <JobList :key="componentKey"/> -->
    </div>
</div>
</template>
  
<script>
import Navbar from "../components/Navbar.vue";
import Sidebar from "../components/Sidebar.vue";
import JobList from "../components/JobList.vue";
import JobDetails from "../components/JobDetails.vue";

export default {
    name: "DashboardView",
    components: {
        Navbar,
        Sidebar,
        JobList,
        JobDetails,
    },
    data() {
      return {
        selectedJob: {IsDir: true},
        componentKey: 0,
      };
    },
    watch: {
      $route(to, from) {
        console.log('Route changed from ', from.fullPath, ' to ', to.fullPath);
        this.refresh()
      }
    },
    computed: {
      currentComponent() {
        // Dynamically decide which component to show
        return this.selectedJob.IsDir ? 'JobList' : 'JobDetails';
      },
    },
    methods: {
      handleJobType(jobType) {
        this.selectedJob.IsDir = jobType;
      },
      // Update the key to refresh
      refresh() {
        this.componentKey += 1;
      },
    },
};
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  height: 100vh;
}
.main-content {
  display: flex;
  flex-grow: 1;
}
</style>