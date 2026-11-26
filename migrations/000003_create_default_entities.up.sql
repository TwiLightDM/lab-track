-- Пароль для администратора Admin1!!
-- Пароль для всех остальных Password1!
INSERT INTO users (name, surname, email, password, role, is_active, created_at, updated_at)
VALUES ('Админ', 'Админов', 'admin@example.com',
        '$2a$10$ZZ/uqkd846qMNgCUjyqHUOXDBkIKPSGaGR39G6jmvSPCdI6xXfzku',
        'admin', TRUE, NOW(), NOW()),

       ('Алиса', 'Дроздова', 'assistant1@example.com',
        '$2a$10$Hkuy1ZQ4UHtDFGqn7eYoK.MQ6cUjFl4OwTSu8J2NexvZprSMHehMK',
        'lab_assistant', TRUE, NOW(), NOW()),
       ('Степан', 'Тетерев', 'assistant2@example.com',
        '$2a$10$Hkuy1ZQ4UHtDFGqn7eYoK.MQ6cUjFl4OwTSu8J2NexvZprSMHehMK',
        'lab_assistant', TRUE, NOW(), NOW()),

       ('Мария', 'Галкина', 'manager@example.com',
        '$2a$10$Hkuy1ZQ4UHtDFGqn7eYoK.MQ6cUjFl4OwTSu8J2NexvZprSMHehMK',
        'lab_manager', TRUE, NOW(), NOW()),

       ('Алексей', 'Воробьёв', 'repairman@example.com',
        '$2a$10$Hkuy1ZQ4UHtDFGqn7eYoK.MQ6cUjFl4OwTSu8J2NexvZprSMHehMK',
        'repairman', TRUE, NOW(), NOW());

INSERT INTO locations (address, room, is_active, created_at, updated_at)
VALUES ('ул. Ленина, д. 10', 101, TRUE, NOW(), NOW()),
       ('ул. Ленина, д. 10', 102, TRUE, NOW(), NOW()),
       ('ул. Ленина, д. 10', 103, TRUE, NOW(), NOW()),

       ('пр. Мира, д. 25', 201, TRUE, NOW(), NOW()),
       ('пр. Мира, д. 25', 202, TRUE, NOW(), NOW());

INSERT INTO equipment_categories (name, descriptions, created_at, updated_at)
VALUES ('Микроскопы', 'Оптические и электронные микроскопы для лабораторных исследований',
        NOW(), NOW()),
       ('Центрифуги', 'Лабораторные центрифуги для разделения жидких сред',
        NOW(), NOW()),
       ('Измерительные приборы', 'Приборы для измерения температуры, давления, влажности и других параметров',
        NOW(), NOW());

INSERT INTO equipments (name, inventory_number, serial_number, state, description,
                        category_id, location_id, created_at, updated_at)
VALUES ('Микроскоп оптический M-100', '8facfb3d-746f-4684-959e-ef875b3a3a28', 'SN-M100-001', 'idle',
        'Бинокулярный микроскоп, увеличение 40–1000x', 1, 1, NOW(), NOW()),
       ('Микроскоп электронный E-200', 'e9c2435a-141f-482d-9f5d-9b8167aacd22', 'SN-E200-001', 'in_use',
        'Сканирующий электронный микроскоп', 1, 2, NOW(), NOW()),

       ('Центрифуга лабораторная C-300', 'aad579a3-01d6-4e53-bd66-e51a9d1720ef', 'SN-C300-001', 'idle',
        'Настольная центрифуга, до 6000 об/мин', 2, 4, NOW(), NOW()),
       ('Центрифуга рефрижераторная C-400', 'e085eea1-e5ed-489a-801e-0bdc6fd3a160', 'SN-C400-001',
        'maintenance_required',
        'Центрифуга с охлаждением, требует ТО', 2, 5, NOW(), NOW()),

       ('Термометр цифровой T-10', '81e8df75-7c48-4354-b9c2-c4aa88cafe79', 'SN-T10-001', 'idle',
        'Термометр с диапазоном -50…+200 °C', 3, 3, NOW(), NOW()),
       ('Манометр прецизионный P-20', 'ed8610f7-39a3-4944-ad48-377e73729c44', 'SN-P20-001', 'in_repair',
        'Прецизионный манометр, 0–10 МПа, в ремонте', 3, 4, NOW(), NOW())