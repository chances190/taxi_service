<script>
  let name = '';
  let email = '';
  let password = '';

  async function register() {
    const response = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, email, password }),
    });
    const data = await response.json();
    if (response.ok) {
      alert('Registration successful!');
      window.location.href = '/';
    } else {
      alert(data.error || 'Registration failed');
    }
  }
</script>

<main>
  <h1>Register</h1>
  <form on:submit|preventDefault={register}>
    <label>Name: <input type="text" bind:value={name} /></label>
    <label>Email: <input type="email" bind:value={email} /></label>
    <label>Password: <input type="password" bind:value={password} /></label>
    <button type="submit">Register</button>
  </form>
</main>
