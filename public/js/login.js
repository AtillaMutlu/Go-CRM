document.getElementById('loginForm').addEventListener('submit', async function(e) {
  e.preventDefault();
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  const errorDiv = document.getElementById('loginError');
  errorDiv.classList.add('d-none');

  try {
    const resp = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    if (!resp.ok) {
      const data = await resp.json();
      throw new Error(data.message || 'Giriş başarısız!');
    }
    const data = await resp.json();
    localStorage.setItem('token', data.token);
    window.location.href = 'dashboard.html';
  } catch (err) {
    errorDiv.textContent = err.message;
    errorDiv.classList.remove('d-none');
  }
}); 