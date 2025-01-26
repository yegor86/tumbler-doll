<template>
<tr @click="navigateToJob" class="job-card">
    <td class="folder-name">
      <span v-if="isFolder">📂</span>
      <span v-else>📄</span>
    </td>
    <td>{{ truncateJobName(job.Name) }}</td>
    <td>{{ job.Status }}</td>
    <td>{{ job.LastBuild }}</td>
</tr>
</template>

<script>
  export default {
    name: 'JobItem',
    props: {
      job: {
        type: Object,
        required: true,
      },
    },
    computed: {
      isFolder() {
        return this.job.IsDir;
      },
    },
    methods: {
      navigateToJob() {
        // '/' -> /jobs/
        // '/Public' -> /jobs/Public
        // '/Public/' -> /jobs/Public/jobs/
        // '/Public/mltibranch' -> /jobs/Public/jobs/mltibranch
        const jobPath = this.job.Name.replace(/\//g, /jobs/);
        this.$router.push(jobPath);
      },

      // '/jobs/Public' -> Public
      // '/Public/mltibranch' -> mltibranch
      // '/Public/mltibranch/' -> mltibranch
      truncateJobName(jobName) {
        return jobName;
      },
    },
  };
</script>

<style>
  .job-card {
    border: 1px solid #ccc;
    padding: 16px;
    border-radius: 8px;
    margin: 8px;
    cursor: pointer;
    transition: box-shadow 0.3s;
  }
  .job-card:hover {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  }
  .job-card .success {
    color: green;
  }
  .job-card .failed {
    color: red;
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