<template>
    <nav class="navbar">
      <div class="logo">Tumbler</div>
      <div class="actions">
        <button class="btn">New Item</button>
        <button class="btn">People</button>
        <button class="btn">Build History</button>
      </div>
    </nav>
  </template>
  
  <script>
  export default {
    name: "NavBar",
  };
  </script>
  
  <style scoped>
  .navbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background-color: #0077b6;
    color: white;
    padding: 10px 20px;
  }
  .btn {
    margin-left: 10px;
    background: #005f73;
    color: white;
    border: none;
    padding: 5px 10px;
    cursor: pointer;
  }
  </style>
  