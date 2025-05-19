-- Police Stations
INSERT INTO police_stations (name, location, contact) VALUES
('Kampala Central Police Station', 'Kampala', '+256 414 255 888'),
('Jinja Road Police Station', 'Jinja Road, Kampala', '+256 414 222 333'),
('Entebbe Police Station', 'Entebbe', '+256 414 320 123');

-- Police Officers
INSERT INTO police_officers (name, username, password_hash, rank, force_number, telephone, email, station_id) VALUES
('John Okello', 'jokello', 'hashedpassword1', 'Inspector', 'UPF12345', '+256 772 111111', 'jokello@police.go.ug', 1),
('Grace Namutebi', 'gnamutebi', 'hashedpassword2', 'Sergeant', 'UPF54321', '+256 772 222222', 'gnamutebi@police.go.ug', 2),
('Samuel Mugisha', 'smugisha', 'hashedpassword3', 'Constable', 'UPF67890', '+256 772 333333', 'smugisha@police.go.ug', 3);

-- Victims
INSERT INTO victims (name, sex, occupation, marital_status, residence, case_number, police_unit, date_of_incident, narrator, relationship, apparent_age, age_estimation_note, incident_history, general_condition, mental_status, head_neck, chest_breast, abdomen_back) VALUES
('Aisha Nansubuga', 'Female', 'Student', 'Single', 'Makindye', 'CASE2024-001', 'Kampala Central Police Station', '2024-05-01', 'Sarah Nakato', 'Mother', '17', 'Based on physical exam', 'Victim was assaulted on her way from school.', 'Stable', 'Calm', 'No visible injuries', 'No visible injuries', 'No visible injuries'),
('Peter Ssemanda', 'Male', 'Boda Boda Rider', 'Married', 'Kawempe', 'CASE2024-002', 'Jinja Road Police Station', '2024-05-03', 'James Lwanga', 'Brother', '28', 'Based on ID', 'Victim was attacked at night.', 'Bruises on arm', 'Distressed', 'Bruising on left arm', 'No visible injuries', 'No visible injuries');

-- Accused
INSERT INTO accuseds (name, sex, occupation, residence, work_place, telephone, case_number, police_unit, date_of_incident, apparent_age, age_estimation_note, hiv_test_results, general_condition, mental_status, head_neck, chest_breast, abdomen_back, upper_lower_limbs, ano_genital) VALUES
('Brian Kato', 'Male', 'Mechanic', 'Nansana', 'Nansana Garage', '+256 772 444444', 'CASE2024-001', 'Kampala Central Police Station', '2024-05-01', '32', 'Based on statement', 'Negative', 'Stable', 'Normal', 'No visible injuries', 'No visible injuries', 'No visible injuries', 'No visible injuries', 'No visible injuries'),
('Ritah Nakimera', 'Female', 'Vendor', 'Kireka', 'Kireka Market', '+256 772 555555', 'CASE2024-002', 'Jinja Road Police Station', '2024-05-03', '24', 'Based on appearance', 'Negative', 'Stable', 'Normal', 'No visible injuries', 'No visible injuries', 'No visible injuries', 'No visible injuries', 'No visible injuries'); 