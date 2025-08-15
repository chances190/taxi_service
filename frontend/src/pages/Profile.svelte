<script>
  let profile = {};

  async function fetchProfile() {
    const response = await fetch('/api/profile/:id');
    if (response.ok) {
      profile = await response.json();
    } else {
      alert('Failed to fetch profile');
    }
  }

  fetchProfile();
</script>

<main>
  <h1>Profile</h1>
  <p>Name: {profile.name}</p>
  <p>Email: {profile.email}</p>
</main>
