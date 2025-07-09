-- Customers tablosu
CREATE TABLE IF NOT EXISTS customers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone VARCHAR(20),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

-- Contacts tablosu
CREATE TABLE IF NOT EXISTS contacts (
  id SERIAL PRIMARY KEY,
  customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
  message TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now()
);

-- İndeksler
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(name);
CREATE INDEX IF NOT EXISTS idx_contacts_customer_id ON contacts(customer_id);
CREATE INDEX IF NOT EXISTS idx_contacts_created_at ON contacts(created_at);

-- Test için örnek data
INSERT INTO customers (name, email, phone) VALUES 
('Demo Müşteri 1', 'demo1@example.com', '+905551111111'),
('Demo Müşteri 2', 'demo2@example.com', '+905552222222'),
('Demo Müşteri 3', 'demo3@example.com', '+905553333333')
ON CONFLICT DO NOTHING;

-- Test contacts
INSERT INTO contacts (customer_id, message) 
SELECT 
  c.id,
  'Demo iletişim mesajı #' || c.id
FROM customers c 
WHERE c.email LIKE 'demo%@example.com'
ON CONFLICT DO NOTHING; 