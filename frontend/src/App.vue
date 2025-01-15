<template>
<div class="pipeline-card" @click="navigateToPipeline">
    <h3>{{ pipeline?.name || 'No Pipeline' }}</h3>
    <p>Status: <span :class="pipeline.status">{{ pipeline?.status || 'Status' }}</span></p>
    <p>Last Run: {{ pipeline?.lastRun || 'Never' }}</p>
</div>
</template>
  
<script>
export default {
    data() {
      return {
        pipeline: Object
      }
    },
    methods: {
      navigateToPipeline() {
          this.$router.push(`/pipeline/${this.pipeline.id}`);
      }
    }
};
</script>

<style>
.pipeline-card {
    border: 1px solid #ccc;
    padding: 16px;
    border-radius: 8px;
    margin: 8px;
    cursor: pointer;
    transition: box-shadow 0.3s;
}
.pipeline-card:hover {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}
.pipeline-card .success {
    color: green;
}
.pipeline-card .failed {
    color: red;
}
</style>
