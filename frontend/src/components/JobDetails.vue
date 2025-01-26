<template>
  <div class="job-details">
    <JobHeader :job="job" />
    <BuildStatus :job="job"  />
    <BuildHistory />
  </div>
</template>

<script>
import JobHeader from "../components/JobHeader.vue";
import BuildStatus from "../components/JobBuildStatus.vue";
import BuildHistory from "../components/JobBuildHistory.vue";

export default {
  name: "JobDetails",
  components: {
    JobHeader,
    BuildStatus,
    BuildHistory,
  },
  data() {
    return {
      job: {Name: "/jobs/job1", Status: "pending"},
    };
  },
};
</script>

<style scoped>
.job-details {
  display: flex;
  flex-direction: column;
  height: 100vh;
}
</style>
