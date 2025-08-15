<script>
  import { Router, Route, Link } from 'svelte-routing';
  import Login from './pages/Login.svelte';
  import Register from './pages/Register.svelte';
  import Profile from './pages/Profile.svelte';
  let message = 'Welcome to Taxi Service';
</script>

<Router>
  <nav>
    <Link to="/login">Login</Link>
    <Link to="/register">Register</Link>
    <Link to="/profile">Profile</Link>
  </nav>

  <main>
    <h1>{message}</h1>
    <Route path="/login" component={Login} />
    <Route path="/register" component={Register} />
    <Route path="/profile" component={Profile} />
  </main>
</Router>

<style>
  nav {
    display: flex;
    justify-content: center;
    gap: 1em;
    margin-bottom: 1em;
  }
  main {
    text-align: center;
    padding: 1em;
    max-width: 240px;
    margin: 0 auto;
  }
  h1 {
    color: #ff3e00;
    font-size: 2em;
  }
</style>
