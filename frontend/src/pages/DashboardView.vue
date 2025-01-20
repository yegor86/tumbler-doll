<template>
<div class="dashboard">
    <Navbar />
    <div class="main-content">
    <Sidebar />
    <JobList />
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