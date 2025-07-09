// API base URL (gerekirse değiştir)
const API_BASE = '/api';

function getToken() {
  return localStorage.getItem('token');
}

async function apiRequest(path, options = {}) {
  const headers = options.headers || {};
  const token = getToken();
  if (token) headers['Authorization'] = 'Bearer ' + token;
  options.headers = headers;
  const resp = await fetch(API_BASE + path, options);
  if (!resp.ok) {
    const data = await resp.json().catch(() => ({}));
    throw new Error(data.message || 'API hatası');
  }
  return resp.json();
}

// Müşteri işlemleri
export async function getCustomers() {
  return apiRequest('/customers');
}
export async function addCustomer(customer) {
  return apiRequest('/customers', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(customer)
  });
}
export async function updateCustomer(id, customer) {
  return apiRequest(`/customers/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(customer)
  });
}
export async function deleteCustomer(id) {
  return apiRequest(`/customers/${id}`, { method: 'DELETE' });
}

// İletişim işlemleri
export async function getContacts() {
  return apiRequest('/contacts');
} 