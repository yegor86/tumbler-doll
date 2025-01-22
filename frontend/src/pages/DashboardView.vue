<template>
<div class="dashboard">
    <Navbar />
    <div class="main-content">
    <Sidebar />
    <JobList :key="componentKey"/>
    </div>
</div>
</template>
  
<script>
import Navbar from "../components/Navbar.vue";
import Sidebar from "../components/Sidebar.vue";
import JobList from "../components/JobList.vue";

export default {
    name: "DashboardView",
    components: {
        Navbar,
        Sidebar,
        JobList,
    },
    data() {
      return {
        componentKey: 0,
      };
    },
    watch: {
      $route(to, from) {
        console.log('Route changed from ', from.fullPath, ' to ', to.fullPath);
        this.refresh()
      }
    },
    methods: {
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