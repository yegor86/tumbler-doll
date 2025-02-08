<template>
    <aside class="sidebar">
      <ul>
        <!-- Folder-specific menu items -->
        <template v-if="isFolder">
          <li><a href="#" @click.prevent="navigateTo('/')">Dashboard</a></li>
          <li><a href="#" @click.prevent="navigateTo('new-job')">New Job</a></li>
          <li><a href="#" @click.prevent="navigateTo('manage-folder')">Manage Folder</a></li>
          <li><a href="#">Manage Jenkins</a></li>
          <li><a href="#">My Views</a></li>
          <li><a href="#">New View</a></li>
        </template>

        <!-- Job-specific menu items -->
        <template v-else>
          <li><a href="#" @click.prevent="submitJob()">Run</a></li>
          <li><a href="#" @click.prevent="navigateTo('configure')">Configure</a></li>
          <li><a href="#" @click.prevent="navigateTo('build-history')">Build History</a></li>
          <li><a href="#" @click.prevent="navigateTo('delete-job')">Delete Job</a></li>
        </template>
      </ul>
    </aside>
  </template>
  
  <script>
  import apiService from "../services/jobService";

  export default {
    name: "SideBar",
    props: {
      isFolder: {
        type: Boolean,
        required: true,
      },
    },
    methods: {
      navigateTo(path) {
        this.$router.push(path);
      },
      async submitJob() {
        
        try {
            const jobPath = this.$route.fullPath; // .replace(/\/jobs\//g, '/')
            const response = await apiService.submitJob(jobPath);

            apiService.streamJobExec(this.$route.fullPath, response.data.WorkflowID);
            
            this.$emit('statusChange', { Status: 'running', Message: response.data });

        } catch (error) {
            this.error = "Error fetching jobs";
            console.error(error);
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .sidebar {
    width: 200px;
    background-color: #f0f0f0;
    padding: 10px;
    height: 100vh;
  }
  .sidebar ul {
    list-style-type: none;
    padding: 0;
  }
  .sidebar li {
    margin: 10px 0;
  }
  .sidebar a {
    text-decoration: none;
    color: #0077b6;
  }
  </style>
  