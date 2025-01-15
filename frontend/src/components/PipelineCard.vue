<template>
  <div class="pipeline-card" @click="navigateToPipeline">
    <h3>{{ pipeline.name }}</h3>
    <p>Status: <span :class="pipeline.status">{{ pipeline.status }}</span></p>
    <p>Last Run: {{ pipeline.lastRun }}</p>
  </div>
</template>

<script>
  export default {
    name: 'PipelineCard',
    props: {
      pipeline: {
        type: Object,
        required: true,
      },
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