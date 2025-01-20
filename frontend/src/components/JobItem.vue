<template>
<tr @click="navigateToJob" class="job-card">
    <td>{{ job.Name }}</td>
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
    methods: {
      navigateToJob() {
        this.$router.push(`${this.job.Name}`);
      }
    }
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