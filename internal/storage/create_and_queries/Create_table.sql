CREATE TABLE IF NOT EXISTS employees(
    id SERIAL PRIMARY KEY,
    fio TEXT NOT NULL,
    phone TEXT NOT NULL,
    salary NUMERIC(10,2) NOT NULL
    );
    INSERT INTO employees (fio, phone, salary) VALUES
    ('Иван Иванов', '123-456-7890', 50000.00),
    ('Мария Петрова', '987-654-3210', 60000.00),
    ('Алексей Смирнов', '555-666-7777', 55000.00);
   


