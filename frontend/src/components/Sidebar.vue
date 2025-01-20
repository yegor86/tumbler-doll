<template>
    <aside class="sidebar">
      <ul>
        <li><a href="#">Dashboard</a></li>
        <li><a href="#">Manage Jenkins</a></li>
        <li><a href="#">My Views</a></li>
        <li><a href="#">New View</a></li>
      </ul>
    </aside>
  </template>
  
  <script>
  export default {
    name: "SideBar",
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
  