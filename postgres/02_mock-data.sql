INSERT INTO court_case (cnj, plaintiff, defendant, court_of_origin, start_date) VALUES
('5001682-88.2024.8.13.0672', 'Alice Johnson', 'Bob Smith', 'TJSP', '2024-01-15'),
('3562061-02.2024.8.13.0431', 'Michael Brown', 'Sarah Davis', 'TJSP', '2024-02-20'),
('2247987-03.2024.8.13.0231', 'Emily Wilson', 'David Lee', 'TJSP', '2024-03-10'),
('6772130-04.2024.8.13.0161', 'Chris Miller', 'Jessica Taylor', 'TJSP', '2024-04-05'),
('2802936-05.2024.8.13.0991', 'Samantha Garcia', 'Matthew Martinez', 'TJSP', '2024-05-18'),
('6712851-06.2024.8.13.0441', 'John Anderson', 'Laura Hernandez', 'TJSP', '2024-06-22'),
('4446422-07.2024.8.13.0221', 'Olivia Perez', 'Kevin Lopez', 'TJSP', '2024-07-08'),
('6992827-08.2024.8.13.0141', 'Daniel Young', 'Megan Gonzalez', 'TJSP', '2024-08-03'),
('7169867-09.2024.8.13.0981', 'Isabella Martinez', 'James Allen', 'TJSP', '2024-09-15'),
('2716321-10.2024.8.13.0871', 'William King', 'Sophia Wright', 'TJSP', '2024-10-29');

INSERT INTO case_update (cnj, update_date, update_details) VALUES
('5001682-88.2024.8.13.0672', '2024-07-31T10:00:00Z', 'Initial hearing scheduled for August 15, 2024.'),
('5001682-88.2024.8.13.0672', '2024-08-01T14:30:00Z', 'Plaintiff submitted additional evidence.'),
('5001682-88.2024.8.13.0672', '2024-08-02T09:00:00Z', 'Defendant requested a delay for response.'),
('3562061-02.2024.8.13.0431', '2024-07-30T11:00:00Z', 'Case file reviewed by judge.'),
('3562061-02.2024.8.13.0431', '2024-08-01T09:30:00Z', 'Witness statements collected.'),
('3562061-02.2024.8.13.0431', '2024-08-02T13:45:00Z', 'Defendant’s lawyer filed a motion for dismissal.'),
('3562061-02.2024.8.13.0431', '2024-08-03T10:00:00Z', 'Hearing date scheduled for August 10, 2024.');
