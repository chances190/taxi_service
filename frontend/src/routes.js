export const routes = [
  { path: '/', component: () => import('./pages/Login.svelte') },
  { path: '/register', component: () => import('./pages/Register.svelte') },
  { path: '/profile', component: () => import('./pages/Profile.svelte') },
];
