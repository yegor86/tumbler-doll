<template>
    <header class="job-header">
      <h1>{{ job.Name }}</h1>
      <div class="controls">
        <button class="btn build-now">Build Now</button>
        <button class="btn configure">Configure</button>
        <button class="btn delete">Delete</button>
      </div>
    </header>
  </template>
  
  <script>
  export default {
    name: "JobHeader",
    props: {
      job: {
        type: Object,
        required: true,
      },
    },
  };
  </script>
  
  <style scoped>
  .job-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    background-color: #0077b6;
    color: white;
  }
  .controls {
    display: flex;
    gap: 10px;
  }
  .btn {
    padding: 8px 12px;
    color: white;
    background-color: #005f73;
    border: none;
    cursor: pointer;
  }
  </style>
  