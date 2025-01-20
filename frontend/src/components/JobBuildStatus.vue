<template>
<div class="build-status">
    <h2>Current Build Status</h2>
    <div :class="statusClass">{{ job.Status }}</div>
</div>
</template>

<script>
export default {
name: "BuildStatus",
props: {
    job: {
    type: Object,
    required: true,
    },
},
computed: {
    statusClass() {
    return {
        success: this.job.Status === "Success",
        failed: this.job.Status === "Failed",
        progress: this.job.Status === "In Progress",
        pending: this.job.Status === "Pending",
    };
    },
},
};
</script>

<style scoped>
.build-status {
margin: 20px;
text-align: center;
}
.success {
color: green;
}
.failed {
color: red;
}
.progress {
color: orange;
}
</style>
