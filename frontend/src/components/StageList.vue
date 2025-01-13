<template>
    <div>
      <h2>Stages</h2>
      <ul>
        <li v-for="stage in stages" :key="stage.name">
          <div>
            <strong>{{ stage.name }}</strong>
            <span>Status: {{ stage.status }}</span>
          </div>
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      stages: Array
    }
  };
  </script>
  