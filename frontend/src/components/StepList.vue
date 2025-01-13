<template>
<div>
<h2>Steps</h2>
<ul>
    <li v-for="step in steps" :key="step.name">
    <span>{{ step.name }}</span>
    <span>Status: {{ step.status }}</span>
    </li>
</ul>
</div>
</template>

<script>
export default {
props: {
steps: Array
}
};
</script>
  