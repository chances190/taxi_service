<script>
  let email = '';
  let password = '';

  async function login() {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
    const data = await response.json();
    if (response.ok) {
      alert('Login successful!');
      window.location.href = '/profile';
    } else {
      alert(data.error || 'Login failed');
    }
  }
</script>

<main>
  <h1>Login</h1>
  <form on:submit|preventDefault={login}>
    <label>Email: <input type="email" bind:value={email} /></label>
    <label>Password: <input type="password" bind:value={password} /></label>
    <button type="submit">Login</button>
  </form>
</main>
