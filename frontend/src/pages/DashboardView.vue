<template>
<div class="dashboard">
    <Navbar />
    <div class="main-content">
    <Sidebar :isFolder="isFolder()" @statusChange="handleStatusChange"/>
    <transition name="fade" mode="out-in">
      <component @navigate="handleNavigation"
        ref="childRef"
        :is="currentComponent"
        :key="componentKey"
      />
    </transition>
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
      handleNavigation(destination) {
        this.selectedJob.IsDir = destination;
      },
      handleStatusChange(event) {
        this.$refs.childRef.handleStatusChange(event);
      },
      isFolder() {
        return this.selectedJob.IsDir;
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
}
.main-content {
  display: flex;
  flex-grow: 1;
}
</style>