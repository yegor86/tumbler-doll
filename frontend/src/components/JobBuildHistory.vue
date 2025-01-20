<template>
    <section class="build-history">
      <h2>Build History</h2>
      <table>
        <thead>
          <tr>
            <th>#</th>
            <th>Status</th>
            <th>Duration</th>
            <th>Timestamp</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="build in builds" :key="build.id">
            <td>{{ build.id }}</td>
            <td>{{ build.status }}</td>
            <td>{{ build.duration }}</td>
            <td>{{ build.timestamp }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </template>
  
  <script>
  export default {
    name: "BuildHistory",
    data() {
      return {
        builds: [
          { id: 1, status: "Success", duration: "5 min", timestamp: "2024-12-23 12:00" },
          { id: 2, status: "Failed", duration: "3 min", timestamp: "2024-12-23 10:00" },
          { id: 3, status: "In Progress", duration: "-", timestamp: "2024-12-23 09:30" },
        ],
      };
    },
  };
  </script>
  
  <style scoped>
  .build-history {
    margin: 20px;
  }
  table {
    width: 100%;
    border-collapse: collapse;
  }
  table th, table td {
    padding: 8px;
    border: 1px solid #ddd;
  }
  table th {
    background-color: #f4f4f4;
  }
  </style>
  