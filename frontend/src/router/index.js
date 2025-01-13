import { createRouter, createWebHistory } from "vue-router";
import DashboardView from "../pages/DashboardView.vue";
import PipelineDetails from "../pages/PipelineDetails.vue";

const routes = [
  { path: "/", component: DashboardView },
  { path: "/pipeline/:id", component: PipelineDetails }
];

export default createRouter({
  history: createWebHistory(),
  routes
});
