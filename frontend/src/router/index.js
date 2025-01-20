import { createRouter, createWebHistory } from "vue-router";
import DashboardView from "../pages/DashboardView.vue";
import JobDetails from "../pages/JobDetailsView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/",             component: DashboardView },
    { path: "/jobs/:id",   component: JobDetails }
  ]
});

export default router;
