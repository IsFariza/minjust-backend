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
    iin VARCHAR(12) NOT NULL UNIQUE,
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50) NOT NULL,
    department VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO management (fullname, position, email, phone, area_of_work, photo_url) VALUES
('Sarsembayev Erlan Zhaksylykovich', 'Minister of Justice of the Republic of Kazakhstan', '', '+7 (7172) 74-02-01', '', 'https://www.gov.kz/uploads/2025/1/9/a0ca29b7dadf93f5571066e3a45d453d_1280x720.jpeg'),
('Mersalimova Laura Kanatovna', 'Vice Minister of Justice of the Republic of Kazakhstan', 'la.mersalimova@adilet.gov.kz', '+7 (7172) 74-02-15', 'Rulemaking', 'https://www.gov.kz/uploads/2024/3/28/7de87829b6d1af1e2afc6c57d3ec83a3_1280x720.jpeg'),
('Moldabekov Bekbolat Serikovich', 'Vice Minister of Justice of the Republic of Kazakhstan', 'b.moldabekov@adilet.gov.kz', '+7 (7172) 74-02-26', 'Digitalization', 'https://www.gov.kz/uploads/2024/10/23/dc0a875ca474b57d0332e7ee324d784b_1280x720.jpg'),
('Zhakselekova Botagoz Shaimardanovna', 'Vice Minister of Justice of the Republic of Kazakhstan', 'b.zhakselekova@adilet.gov.kz', '+7 (7172) 74-02-14', 'International cooperation, registration service and legal services, intellectual property', 'https://www.gov.kz/uploads/2024/10/23/abdbdae1250cf5483ece47ea28d8563d_1280x720.jpg'),
('Vaissov Daniel Mereevich', 'Vice Minister of Justice of the Republic of Kazakhstan', 'd.vaissov@adilet.gov.kz', '+7 (7172) 74-01-21', 'Forensic expert activity, enforcement documents, protection of state property rights', 'https://www.gov.kz/uploads/2025/9/23/180d4ee20428e16fa85f0b244b5dff15_1280x720.jpg'),
('Erseitova Sandugash Abdrazakovna', 'Chief of Staff of the Ministry of Justice of the Republic of Kazakhstan', 's.erseitova@adilet.gov.kz', '+7 (7172) 74-01-06', 'Human resources, economics and finance, public relations, document control', 'https://www.gov.kz/uploads/2023/2/24/77ad3baa6b0ce58c76af059e7196c37e_1280x720.jpg');
