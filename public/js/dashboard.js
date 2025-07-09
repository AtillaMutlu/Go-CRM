import { getCustomers, addCustomer, updateCustomer, deleteCustomer, getContacts, addContact } from './api.js';

// Token yoksa login'e yönlendir
if (!localStorage.getItem('token')) {
  window.location.href = 'index.html';
}

document.getElementById('logoutBtn').onclick = function() {
  localStorage.removeItem('token');
  window.location.href = 'index.html';
};

// Müşteri Tablosu
const customerTableBody = document.querySelector('#customerTable tbody');
const customerForm = document.getElementById('customerForm');
const customerModal = new bootstrap.Modal(document.getElementById('customerModal'));
const contactModal = new bootstrap.Modal(document.getElementById('contactModal'));
const contactForm = document.getElementById('contactForm');

let editingId = null;

async function loadCustomers() {
  customerTableBody.innerHTML = '<tr><td colspan="4">Yükleniyor...</td></tr>';
  try {
    const customers = await getCustomers();
    customerTableBody.innerHTML = '';
    customers.forEach(c => {
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td>${c.name}</td>
        <td>${c.email}</td>
        <td>${c.phone || ''}</td>
        <td>
          <button class="btn btn-sm btn-info" onclick="addContactHandler('${c.id}')">İletişim Ekle</button>
          <button class="btn btn-sm btn-warning me-1" onclick="editCustomer('${c.id}')">Düzenle</button>
          <button class="btn btn-sm btn-danger" onclick="deleteCustomerHandler('${c.id}')">Sil</button>
        </td>
      `;
      customerTableBody.appendChild(tr);
    });
  } catch (err) {
    customerTableBody.innerHTML = `<tr><td colspan="4">${err.message}</td></tr>`;
  }
}

window.editCustomer = function(id) {
  getCustomers().then(customers => {
    const c = customers.find(x => x.id == id);
    if (!c) return;
    editingId = c.id;
    document.getElementById('customerId').value = c.id;
    document.getElementById('customerName').value = c.name;
    document.getElementById('customerEmail').value = c.email;
    document.getElementById('customerPhone').value = c.phone || '';
    document.getElementById('customerModalLabel').textContent = 'Müşteri Düzenle';
    customerModal.show();
  });
};

window.deleteCustomerHandler = async function(id) {
  if (!confirm('Silmek istediğine emin misin?')) return;
  await deleteCustomer(id);
  loadCustomers();
};

window.addContactHandler = function(customerId) {
  contactForm.reset();
  document.getElementById('contactCustomerId').value = customerId;
  contactModal.show();
};

contactForm.onsubmit = async function(e) {
  e.preventDefault();
  const customer_id = parseInt(document.getElementById('contactCustomerId').value, 10);
  const message = document.getElementById('contactMessage').value;
  
  try {
    await addContact({ customer_id, message });
    contactModal.hide();
    // İletişimler sekmesini aktif hale getir ve listeyi yenile
    const contactTab = document.querySelector('a[href="#contacts"]');
    if(contactTab) {
      new bootstrap.Tab(contactTab).show();
      loadContacts();
    }
  } catch (err) {
    alert('Hata: ' + err.message);
  }
};

customerForm.onsubmit = async function(e) {
  e.preventDefault();
  const id = document.getElementById('customerId').value;
  const name = document.getElementById('customerName').value;
  const email = document.getElementById('customerEmail').value;
  const phone = document.getElementById('customerPhone').value;
  if (id) {
    await updateCustomer(id, { name, email, phone });
  } else {
    await addCustomer({ name, email, phone });
  }
  customerModal.hide();
  customerForm.reset();
  editingId = null;
  loadCustomers();
};

document.getElementById('customerModal').addEventListener('hidden.bs.modal', function() {
  customerForm.reset();
  editingId = null;
  document.getElementById('customerModalLabel').textContent = 'Müşteri Ekle';
});

// İletişim Tablosu
const contactTableBody = document.querySelector('#contactTable tbody');
async function loadContacts() {
  contactTableBody.innerHTML = '<tr><td colspan="3">Yükleniyor...</td></tr>';
  try {
    const contacts = await getContacts();
    contactTableBody.innerHTML = '';
    // Gelen veri null değilse ve bir dizi ise döngüye gir
    if (contacts && Array.isArray(contacts)) {
      contacts.forEach(c => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
          <td>${c.customer_name}</td>
          <td>${c.message}</td>
          <td>${c.date}</td>
        `;
        contactTableBody.appendChild(tr);
      });
    }
  } catch (err) {
    contactTableBody.innerHTML = `<tr><td colspan="3">${err.message}</td></tr>`;
  }
}

// Tab değişiminde ilgili tabloyu yükle
const tabMenu = document.getElementById('tabMenu');
tabMenu.addEventListener('click', function(e) {
  if (e.target.matches('[data-bs-toggle="tab"]')) {
    if (e.target.getAttribute('href') === '#customers') loadCustomers();
    if (e.target.getAttribute('href') === '#contacts') loadContacts();
  }
});

// İlk yüklemede müşteri tablosunu getir
loadCustomers(); 