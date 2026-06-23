CREATE TABLE IF NOT EXISTS management (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(255),
    area_of_work TEXT,
    photo_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employee_accounts (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    iin VARCHAR(12) NOT NULL UNIQUE,
    position VARCHAR(255) NOT NULL,
    department VARCHAR(255) NOT NULL,
    management VARCHAR(255) NOT NULL,
    cabinet VARCHAR(50),
    phone_work VARCHAR(50),
    phone_personal VARCHAR(50),
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO management (fullname, position, email, phone, area_of_work, photo_url) VALUES
('Сарсембаев Ерлан Жаксылыкович', 'Министр юстиции Республики Казахстан', '-', '+7 (7172) 74-02-01', '-', 'https://www.gov.kz/uploads/2025/1/9/a0ca29b7dadf93f5571066e3a45d453d_1280x720.jpeg'),
('Мерсалимова Лаура Канатовна', 'Вице-министр юстиции Республики Казахстан', 'la.mersalimova@adilet.gov.kz', '+7 (7172) 74-02-15', 'Нормотворчество', 'https://www.gov.kz/uploads/2024/3/28/7de87829b6d1af1e2afc6c57d3ec83a3_1280x720.jpeg'),
('Молдабеков Бекболат Серикович', 'Вице-министр юстиции Республики Казахстан', 'b.moldabekov@adilet.gov.kz', '+7 (7172) 74-02-26', 'Цифровизация', 'https://www.gov.kz/uploads/2024/10/23/dc0a875ca474b57d0332e7ee324d784b_1280x720.jpg'),
('Жакселекова Ботагоз Шаймардановна', 'Вице-министр юстиции Республики Казахстан', 'b.zhakselekova@adilet.gov.kz', '+7 (7172) 74-02-14', 'Международное сотрудничество, Регистрационная служба и организация юридических услуг, Интеллектуальная собственность, Государственная регистрация нормативных правовых актов', 'https://www.gov.kz/uploads/2024/10/23/abdbdae1250cf5483ece47ea28d8563d_1280x720.jpg'),
('Ваисов Даниель Мереевич', 'Вице-министр юстиции Республики Казахстан', 'd.vaissov@adilet.gov.kz', '+7 (7172) 74-01-21', 'Судебно экспертная деятельность, Исполнение исполнительных документов, Защита имущественных прав государства', 'https://www.gov.kz/uploads/2025/9/23/180d4ee20428e16fa85f0b244b5dff15_1280x720.jpg'),
('Ерсеитова Сандугаш Абдразаковна', 'Руководитель аппарата Министерства юстиции Республики Казахстан', 's.erseitova@adilet.gov.kz', '+7 (7172) 74-01-06', 'Кадровое обеспечение, Экономика и финансы, Связь с общественностью, Регистрация и контроль документооборота', 'https://www.gov.kz/uploads/2023/2/24/77ad3baa6b0ce58c76af059e7196c37e_1280x720.jpg');
