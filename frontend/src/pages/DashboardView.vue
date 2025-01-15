<template>
<div>
    <h1>Jenkins-like Dashboard</h1>
    <div class="pipeline-list">
    <PipelineCard
        v-for="pipeline in pipelines"
        :key="pipeline.id"
        :pipeline="pipeline"
    />
    </div>
</div>
</template>
  
<script>
import PipelineCard from "../components/PipelineCard.vue";

export default {
    components: { PipelineCard },
    data() {
        return {
            pipelines: [
                { id: 1, name: "Pipeline 1", status: "success", lastRun: "10m ago" },
                { id: 2, name: "Pipeline 2", status: "failed", lastRun: "30m ago" },
                { id: 3, name: "Pipeline 3", status: "running", lastRun: "5m ago" }
            ]
        };
    }
};
</script>
