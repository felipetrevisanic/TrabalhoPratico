INSERT INTO public.products ("Name", "Description", "Category", "Images", "Price", "StockQuantity", "CreatedAt")
VALUES
    ('Notebook', 'Notebook para desenvolvimento', 'Computers', ARRAY['https://cdn.example.com/notebook/front.png', 'https://cdn.example.com/notebook/back.png'], 4500.00, 10, NOW()),
    ('Mouse', 'Mouse sem fio', 'Accessories', ARRAY['https://cdn.example.com/mouse/front.png', 'https://cdn.example.com/mouse/side.png'], 120.50, 25, NOW());
