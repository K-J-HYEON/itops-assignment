// src/router/index.js

import { createRouter, createWebHistory } from 'vue-router'
import IssueList from '../pages/IssueList.vue'
import IssueForm from '../pages/IssueForm.vue'

const routes = [
  {
    path: '/',
    name: 'IssueList',
    component: IssueList,
  },
  {
    path: '/issue/new',
    name: 'CreateIssue',
    component: IssueForm,
  },
  {
    path: '/issue/edit/:id',
    name: 'EditIssue',
    component: IssueForm,
    props: true,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
