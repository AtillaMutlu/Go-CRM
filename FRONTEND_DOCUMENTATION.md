# Frontend Documentation

## Table of Contents
1. [Overview](#overview)
2. [File Structure](#file-structure)
3. [HTML Pages](#html-pages)
4. [JavaScript Modules](#javascript-modules)
5. [CSS Styling](#css-styling)
6. [API Integration](#api-integration)
7. [Authentication Flow](#authentication-flow)
8. [Error Handling](#error-handling)
9. [User Interface Components](#user-interface-components)
10. [Browser Compatibility](#browser-compatibility)

## Overview

The frontend is built with vanilla HTML, CSS, and JavaScript using ES6 modules. It provides a modern, responsive interface for the CRM application with real-time data management capabilities.

### Technologies
- **HTML5**: Semantic markup
- **CSS3**: Modern styling with Flexbox and Grid
- **Vanilla JavaScript**: ES6 modules, async/await, fetch API
- **Bootstrap**: UI framework for responsive design
- **Local Storage**: Client-side token storage

## File Structure

```
public/
├── index.html          # Login page
├── dashboard.html      # Main application dashboard
├── js/
│   ├── api.js         # API client module
│   ├── login.js       # Login functionality
│   └── dashboard.js   # Dashboard functionality
└── css/               # Stylesheets (if any)
```

## HTML Pages

### Login Page (`index.html`)

**Purpose**: User authentication interface

**Key Elements**:
- Login form with email and password fields
- Error display area
- Responsive design

**Structure**:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CRM Login</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-5">
        <div class="row justify-content-center">
            <div class="col-md-6 col-lg-4">
                <div class="card">
                    <div class="card-body">
                        <h3 class="card-title text-center mb-4">CRM Login</h3>
                        <form id="loginForm">
                            <div class="mb-3">
                                <label for="email" class="form-label">Email</label>
                                <input type="email" class="form-control" id="email" required>
                            </div>
                            <div class="mb-3">
                                <label for="password" class="form-label">Password</label>
                                <input type="password" class="form-control" id="password" required>
                            </div>
                            <div class="alert alert-danger d-none" id="loginError"></div>
                            <button type="submit" class="btn btn-primary w-100">Login</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <script type="module" src="js/login.js"></script>
</body>
</html>
```

### Dashboard Page (`dashboard.html`)

**Purpose**: Main application interface for customer and contact management

**Key Elements**:
- Navigation header
- Customer management section
- Contact management section
- Modal dialogs for forms
- Data tables

**Structure**:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CRM Dashboard</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <!-- Navigation -->
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary">
        <div class="container">
            <a class="navbar-brand" href="#">CRM System</a>
            <button class="btn btn-outline-light" onclick="logout()">Logout</button>
        </div>
    </nav>

    <div class="container mt-4">
        <!-- Customer Section -->
        <div class="row mb-4">
            <div class="col">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h5 class="mb-0">Customers</h5>
                        <button class="btn btn-primary btn-sm" onclick="showAddCustomerModal()">
                            Add Customer
                        </button>
                    </div>
                    <div class="card-body">
                        <div id="customersTable"></div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Contact Section -->
        <div class="row">
            <div class="col">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h5 class="mb-0">Contacts</h5>
                        <button class="btn btn-primary btn-sm" onclick="showAddContactModal()">
                            Add Contact
                        </button>
                    </div>
                    <div class="card-body">
                        <div id="contactsTable"></div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Modals -->
    <!-- Add Customer Modal -->
    <div class="modal fade" id="addCustomerModal" tabindex="-1">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Add Customer</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <form id="addCustomerForm">
                        <div class="mb-3">
                            <label for="customerName" class="form-label">Name</label>
                            <input type="text" class="form-control" id="customerName" required>
                        </div>
                        <div class="mb-3">
                            <label for="customerEmail" class="form-label">Email</label>
                            <input type="email" class="form-control" id="customerEmail" required>
                        </div>
                        <div class="mb-3">
                            <label for="customerPhone" class="form-label">Phone</label>
                            <input type="tel" class="form-control" id="customerPhone">
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
                    <button type="button" class="btn btn-primary" onclick="addCustomer()">Add Customer</button>
                </div>
            </div>
        </div>
    </div>

    <!-- Add Contact Modal -->
    <div class="modal fade" id="addContactModal" tabindex="-1">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Add Contact</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <form id="addContactForm">
                        <div class="mb-3">
                            <label for="contactCustomer" class="form-label">Customer</label>
                            <select class="form-select" id="contactCustomer" required>
                                <option value="">Select Customer</option>
                            </select>
                        </div>
                        <div class="mb-3">
                            <label for="contactMessage" class="form-label">Message</label>
                            <textarea class="form-control" id="contactMessage" rows="3" required></textarea>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
                    <button type="button" class="btn btn-primary" onclick="addContact()">Add Contact</button>
                </div>
            </div>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
    <script type="module" src="js/dashboard.js"></script>
</body>
</html>
```

## JavaScript Modules

### API Client (`js/api.js`)

**Purpose**: Centralized API communication module

**Features**:
- Automatic JWT token inclusion
- Error handling
- Promise-based requests

#### Configuration

```javascript
const API_BASE = '/api';
```

#### Token Management

```javascript
function getToken() {
  return localStorage.getItem('token');
}
```

#### Core Request Function

```javascript
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
```

**Parameters**:
- `path`: API endpoint path
- `options`: Fetch options (method, headers, body)

**Returns**: Promise with JSON response

**Error Handling**: Throws Error with message from API or default message

#### Customer Operations

##### `getCustomers()`
Retrieves all customers.

```javascript
export async function getCustomers() {
  return apiRequest('/customers');
}
```

**Returns**: Promise<Array<Customer>>

**Example Usage**:
```javascript
import { getCustomers } from './api.js';

try {
  const customers = await getCustomers();
  console.log('Customers:', customers);
} catch (error) {
  console.error('Failed to fetch customers:', error);
}
```

##### `addCustomer(customer)`
Creates a new customer.

```javascript
export async function addCustomer(customer) {
  return apiRequest('/customers', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(customer)
  });
}
```

**Parameters**:
- `customer`: Customer object with name, email, phone

**Returns**: Promise<Customer>

**Example Usage**:
```javascript
import { addCustomer } from './api.js';

const newCustomer = {
  name: 'John Doe',
  email: 'john@example.com',
  phone: '+1234567890'
};

try {
  const createdCustomer = await addCustomer(newCustomer);
  console.log('Customer created:', createdCustomer);
} catch (error) {
  console.error('Failed to create customer:', error);
}
```

##### `updateCustomer(id, customer)`
Updates an existing customer.

```javascript
export async function updateCustomer(id, customer) {
  return apiRequest(`/customers/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(customer)
  });
}
```

**Parameters**:
- `id`: Customer ID
- `customer`: Updated customer data

**Returns**: Promise<Customer>

##### `deleteCustomer(id)`
Deletes a customer.

```javascript
export async function deleteCustomer(id) {
  return apiRequest(`/customers/${id}`, { method: 'DELETE' });
}
```

**Parameters**:
- `id`: Customer ID

**Returns**: Promise<void>

#### Contact Operations

##### `getContacts()`
Retrieves all contacts.

```javascript
export async function getContacts() {
  return apiRequest('/contacts');
}
```

**Returns**: Promise<Array<Contact>>

##### `addContact(data)`
Creates a new contact.

```javascript
export async function addContact(data) {
  return apiRequest('/contacts', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });
}
```

**Parameters**:
- `data`: Contact object with customer_id and message

**Returns**: Promise<Contact>

### Login Module (`js/login.js`)

**Purpose**: Handles user authentication

#### Event Listener Setup

```javascript
document.getElementById('loginForm').addEventListener('submit', async function(e) {
  e.preventDefault();
  // Login logic
});
```

#### Login Function

```javascript
async function handleLogin(email, password) {
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
}
```

**Parameters**:
- `email`: User email
- `password`: User password

**Flow**:
1. Sends login request to `/api/login`
2. Stores JWT token in localStorage
3. Redirects to dashboard

**Error Handling**:
- Displays error message in UI
- Handles network errors

### Dashboard Module (`js/dashboard.js`)

**Purpose**: Main application logic and UI management

#### Initialization

```javascript
// Check authentication on page load
if (!localStorage.getItem('token')) {
  window.location.href = 'index.html';
}

// Load initial data
loadCustomers();
loadContacts();
```

#### Customer Management

##### `loadCustomers()`
Loads and displays customers.

```javascript
async function loadCustomers() {
  try {
    const customers = await getCustomers();
    displayCustomers(customers);
  } catch (error) {
    showError('Failed to load customers: ' + error.message);
  }
}
```

##### `displayCustomers(customers)`
Renders customer table.

```javascript
function displayCustomers(customers) {
  const table = document.getElementById('customersTable');
  
  if (customers.length === 0) {
    table.innerHTML = '<p class="text-muted">No customers found.</p>';
    return;
  }
  
  const html = `
    <table class="table table-striped">
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>Email</th>
          <th>Phone</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        ${customers.map(customer => `
          <tr>
            <td>${customer.id}</td>
            <td>${escapeHtml(customer.name)}</td>
            <td>${escapeHtml(customer.email)}</td>
            <td>${escapeHtml(customer.phone || '')}</td>
            <td>${new Date(customer.created_at).toLocaleDateString()}</td>
            <td>
              <button class="btn btn-sm btn-outline-primary" onclick="editCustomer(${customer.id})">
                Edit
              </button>
              <button class="btn btn-sm btn-outline-danger" onclick="deleteCustomer(${customer.id})">
                Delete
              </button>
            </td>
          </tr>
        `).join('')}
      </tbody>
    </table>
  `;
  
  table.innerHTML = html;
}
```

##### `addCustomer()`
Handles customer creation.

```javascript
async function addCustomer() {
  const name = document.getElementById('customerName').value;
  const email = document.getElementById('customerEmail').value;
  const phone = document.getElementById('customerPhone').value;
  
  if (!name || !email) {
    showError('Name and email are required');
    return;
  }
  
  try {
    await addCustomer({ name, email, phone });
    hideModal('addCustomerModal');
    loadCustomers();
    showSuccess('Customer added successfully');
  } catch (error) {
    showError('Failed to add customer: ' + error.message);
  }
}
```

#### Contact Management

##### `loadContacts()`
Loads and displays contacts.

```javascript
async function loadContacts() {
  try {
    const contacts = await getContacts();
    displayContacts(contacts);
  } catch (error) {
    showError('Failed to load contacts: ' + error.message);
  }
}
```

##### `displayContacts(contacts)`
Renders contact table.

```javascript
function displayContacts(contacts) {
  const table = document.getElementById('contactsTable');
  
  if (contacts.length === 0) {
    table.innerHTML = '<p class="text-muted">No contacts found.</p>';
    return;
  }
  
  const html = `
    <table class="table table-striped">
      <thead>
        <tr>
          <th>ID</th>
          <th>Customer</th>
          <th>Message</th>
          <th>Date</th>
        </tr>
      </thead>
      <tbody>
        ${contacts.map(contact => `
          <tr>
            <td>${contact.id}</td>
            <td>${escapeHtml(contact.customer_name)}</td>
            <td>${escapeHtml(contact.message)}</td>
            <td>${contact.date}</td>
          </tr>
        `).join('')}
      </tbody>
    </table>
  `;
  
  table.innerHTML = html;
}
```

##### `addContact()`
Handles contact creation.

```javascript
async function addContact() {
  const customerId = document.getElementById('contactCustomer').value;
  const message = document.getElementById('contactMessage').value;
  
  if (!customerId || !message) {
    showError('Customer and message are required');
    return;
  }
  
  try {
    await addContact({ customer_id: parseInt(customerId), message });
    hideModal('addContactModal');
    loadContacts();
    showSuccess('Contact added successfully');
  } catch (error) {
    showError('Failed to add contact: ' + error.message);
  }
}
```

#### Modal Management

##### `showAddCustomerModal()`
Opens customer creation modal.

```javascript
function showAddCustomerModal() {
  const modal = new bootstrap.Modal(document.getElementById('addCustomerModal'));
  modal.show();
}
```

##### `showAddContactModal()`
Opens contact creation modal and populates customer dropdown.

```javascript
async function showAddContactModal() {
  try {
    const customers = await getCustomers();
    const select = document.getElementById('contactCustomer');
    select.innerHTML = '<option value="">Select Customer</option>' +
      customers.map(c => `<option value="${c.id}">${escapeHtml(c.name)}</option>`).join('');
    
    const modal = new bootstrap.Modal(document.getElementById('addContactModal'));
    modal.show();
  } catch (error) {
    showError('Failed to load customers: ' + error.message);
  }
}
```

##### `hideModal(modalId)`
Closes modal and resets form.

```javascript
function hideModal(modalId) {
  const modal = bootstrap.Modal.getInstance(document.getElementById(modalId));
  modal.hide();
  document.getElementById(modalId.replace('Modal', 'Form')).reset();
}
```

#### Utility Functions

##### `escapeHtml(text)`
Escapes HTML to prevent XSS.

```javascript
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}
```

##### `showError(message)`
Displays error message.

```javascript
function showError(message) {
  // Implementation depends on UI framework
  alert(message); // Simple implementation
}
```

##### `showSuccess(message)`
Displays success message.

```javascript
function showSuccess(message) {
  // Implementation depends on UI framework
  alert(message); // Simple implementation
}
```

##### `logout()`
Handles user logout.

```javascript
function logout() {
  localStorage.removeItem('token');
  window.location.href = 'index.html';
}
```

## CSS Styling

### Bootstrap Integration

The application uses Bootstrap 5.1.3 for responsive design and UI components.

**CDN Link**:
```html
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
```

### Custom Styles

While the application primarily uses Bootstrap classes, custom CSS can be added for specific styling needs.

**Example Custom Styles**:
```css
/* Custom card styling */
.card {
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  border: none;
}

/* Custom button styling */
.btn-primary {
  background-color: #007bff;
  border-color: #007bff;
}

/* Responsive table */
.table-responsive {
  overflow-x: auto;
}

/* Modal backdrop */
.modal-backdrop {
  background-color: rgba(0,0,0,0.5);
}
```

## API Integration

### Request Flow

1. **Authentication**: JWT token stored in localStorage
2. **Request Preparation**: Token added to Authorization header
3. **API Call**: Fetch request to backend
4. **Response Handling**: JSON parsing and error handling
5. **UI Update**: DOM manipulation based on response

### Error Handling

```javascript
try {
  const response = await apiRequest('/customers');
  // Handle success
} catch (error) {
  if (error.message.includes('401')) {
    // Authentication error - redirect to login
    localStorage.removeItem('token');
    window.location.href = 'index.html';
  } else {
    // Other errors - show to user
    showError(error.message);
  }
}
```

### Loading States

```javascript
function showLoading(elementId) {
  const element = document.getElementById(elementId);
  element.innerHTML = '<div class="text-center"><div class="spinner-border" role="status"></div></div>';
}

function hideLoading(elementId) {
  // Reload data or restore original content
}
```

## Authentication Flow

### Login Process

1. User enters credentials
2. Form submission prevented
3. API request to `/api/login`
4. JWT token received and stored
5. Redirect to dashboard

### Token Management

```javascript
// Store token
localStorage.setItem('token', token);

// Retrieve token
const token = localStorage.getItem('token');

// Remove token (logout)
localStorage.removeItem('token');
```

### Authentication Check

```javascript
// Check if user is authenticated
if (!localStorage.getItem('token')) {
  window.location.href = 'index.html';
}
```

### Token Expiration

The frontend doesn't handle token expiration automatically. The backend returns 401 for expired tokens, which triggers a redirect to login.

## Error Handling

### Network Errors

```javascript
catch (error) {
  if (error.name === 'TypeError' && error.message.includes('fetch')) {
    showError('Network error. Please check your connection.');
  } else {
    showError(error.message);
  }
}
```

### API Errors

```javascript
if (!resp.ok) {
  const data = await resp.json().catch(() => ({}));
  throw new Error(data.message || `HTTP ${resp.status}: ${resp.statusText}`);
}
```

### User Feedback

```javascript
function showError(message) {
  // Could use toast notifications, alerts, or inline messages
  alert(message);
}

function showSuccess(message) {
  // Could use toast notifications or success banners
  alert(message);
}
```

## User Interface Components

### Tables

**Customer Table**:
- ID, Name, Email, Phone, Created Date, Actions
- Edit and Delete buttons for each row
- Responsive design

**Contact Table**:
- ID, Customer Name, Message, Date
- Read-only display

### Forms

**Customer Form**:
- Name (required)
- Email (required, email validation)
- Phone (optional)

**Contact Form**:
- Customer selection (dropdown)
- Message (required, textarea)

### Modals

**Add Customer Modal**:
- Form with validation
- Cancel and Submit buttons
- Bootstrap modal implementation

**Add Contact Modal**:
- Dynamic customer dropdown
- Message textarea
- Form validation

### Navigation

**Header**:
- Application title
- Logout button
- Bootstrap navbar

## Browser Compatibility

### Supported Browsers

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

### ES6 Features Used

- `async/await`
- `const/let`
- Arrow functions
- Template literals
- ES6 modules (`import/export`)
- `fetch` API

### Polyfills

For older browsers, consider adding polyfills for:
- `fetch` API
- ES6 Promise
- ES6 modules

### Feature Detection

```javascript
// Check for fetch support
if (!window.fetch) {
  // Load fetch polyfill or show error
  console.error('Fetch API not supported');
}

// Check for localStorage support
if (!window.localStorage) {
  console.error('localStorage not supported');
}
```

## Performance Considerations

### Code Splitting

The application uses ES6 modules for better code organization and potential tree-shaking.

### Lazy Loading

Consider implementing lazy loading for:
- Large datasets
- Modal content
- Images (if added)

### Caching

- API responses could be cached in memory
- Static assets cached by browser
- Service Worker for offline support

### Optimization

- Minimize DOM queries
- Use event delegation
- Debounce API calls
- Virtual scrolling for large lists

## Security Considerations

### XSS Prevention

```javascript
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}
```

### CSRF Protection

The application relies on JWT tokens for CSRF protection.

### Input Validation

```javascript
function validateEmail(email) {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email);
}

function validatePhone(phone) {
  const re = /^[\+]?[1-9][\d]{0,15}$/;
  return re.test(phone);
}
```

### Token Security

- Tokens stored in localStorage (consider httpOnly cookies for production)
- Automatic logout on 401 responses
- Token expiration handled by backend

## Testing

### Unit Testing

```javascript
// Example test structure
describe('API Client', () => {
  test('getCustomers returns array', async () => {
    // Mock fetch
    // Test function
    // Assert result
  });
});
```

### Integration Testing

```javascript
// Example test structure
describe('Dashboard', () => {
  test('loads customers on page load', async () => {
    // Setup DOM
    // Mock API
    // Trigger load
    // Assert UI updates
  });
});
```

### Manual Testing

1. **Login Flow**:
   - Valid credentials
   - Invalid credentials
   - Network errors

2. **Customer Management**:
   - Add customer
   - Edit customer
   - Delete customer
   - Form validation

3. **Contact Management**:
   - Add contact
   - Customer dropdown population
   - Message validation

4. **Authentication**:
   - Token expiration
   - Logout
   - Unauthorized access

## Future Enhancements

1. **Real-time Updates**: WebSocket integration
2. **Offline Support**: Service Worker implementation
3. **Advanced UI**: React/Vue.js migration
4. **State Management**: Redux/Vuex integration
5. **Form Validation**: Client-side validation library
6. **Notifications**: Toast notification system
7. **Search/Filter**: Advanced data filtering
8. **Pagination**: Large dataset handling
9. **Export**: Data export functionality
10. **Themes**: Dark/light mode support