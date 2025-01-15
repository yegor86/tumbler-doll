<template>
<div>
    <h1>{{ pipeline.name }}</h1>
    <StageList :stages="pipeline.stages" />
</div>
</template>

<script>
import StageList from "../components/StageList.vue";

export default {
components: { StageList },
data() {
    return {
      pipeline: {
        name: "Pipeline 1",
        stages: [
            { name: "Build", status: "success" },
            { name: "Test", status: "running" },
            { name: "Deploy", status: "pending" }
        ]
      }
    };
}
};
</script>
  