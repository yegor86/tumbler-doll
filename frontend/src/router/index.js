import { createRouter, createWebHistory } from "vue-router";
import DashboardView from "../pages/DashboardView.vue";
import PipelineDetails from "../pages/PipelineDetails.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/",             component: DashboardView },
    { path: "/pipeline/:id", component: PipelineDetails }
  ]
});

export default router;
