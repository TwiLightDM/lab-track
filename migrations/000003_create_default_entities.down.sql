DELETE
FROM equipments
WHERE inventory_number IN (
                           'ed8610f7-39a3-4944-ad48-377e73729c44', '8facfb3d-746f-4684-959e-ef875b3a3a28',
                           'e9c2435a-141f-482d-9f5d-9b8167aacd22', 'aad579a3-01d6-4e53-bd66-e51a9d1720ef',
                           'e085eea1-e5ed-489a-801e-0bdc6fd3a160', '81e8df75-7c48-4354-b9c2-c4aa88cafe79'
    );

DELETE
FROM equipment_categories
WHERE name IN (
               'Микроскопы',
               'Центрифуги',
               'Измерительные приборы'
    );

DELETE
FROM locations
WHERE (address, room) IN (
                          ('ул. Ленина, д. 10', 101),
                          ('ул. Ленина, д. 10', 102),
                          ('ул. Ленина, д. 10', 103),
                          ('пр. Мира, д. 25', 201),
                          ('пр. Мира, д. 25', 202)
    );

DELETE
FROM users
WHERE email IN (
                'admin@example.com',
                'assistant1@example.com',
                'assistant2@example.com',
                'manager@example.com',
                'repairman@example.com'
    );