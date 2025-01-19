<template>
<div>
    <h1>Jenkins-like Dashboard</h1>
    <div class="pipeline-list">
    <pipeline-card v-for="pipeline in pipelines"
        :key="pipeline.id"
        :pipeline="pipeline"
    />
    </div>
</div>
</template>
  
<script>
import PipelineCard from "../components/PipelineCard.vue";
import apiService from "../services/jobService";

export default {
    components: { PipelineCard },
    data() {
        return {
            pipelines: [],
            error: null,
        }
    },
    // data() {
    //     return {
    //         pipelines: [
    //             { id: 1, name: "Pipeline 1", status: "success", lastRun: "10m ago" },
    //             { id: 2, name: "Pipeline 2", status: "failed", lastRun: "30m ago" },
    //             { id: 3, name: "Pipeline 3", status: "running", lastRun: "5m ago" }
    //         ]
    //     };
    // }
    mounted() {
        this.fetchData();
    },
    methods: {
        async fetchData() {
            try {
                const response = await apiService.getJobs();
                this.pipelines = response.data;
            } catch (error) {
                this.error = "Error fetching jobs";
                console.error(error);
            }
        }
    },
};
</script>
