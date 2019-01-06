import Vue from "vue";
import Router from "vue-router";
import Upload from "./views/Upload.vue";

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: "/",
      name: "upload",
      component: Upload,
    },
  ],
});
