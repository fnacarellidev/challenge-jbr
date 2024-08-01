CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE court_case (
	id				UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	cnj				TEXT NOT NULL UNIQUE,
	plaintiff		TEXT NOT NULL,
	defendant		TEXT NOT NULL,
	court_of_origin TEXT NOT NULL,
	start_date		DATE NOT NULL,
	created_at		TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE case_update (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cnj              TEXT NOT NULL,
    update_date      TIMESTAMPTZ NOT NULL,
    update_details   TEXT NOT NULL,
	created_at       TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (cnj) REFERENCES court_case(cnj) ON DELETE CASCADE
);

INSERT INTO court_case (cnj, plaintiff, defendant, court_of_origin, start_date) VALUES
('5001682-88.2024.8.13.0672', 'Alice Johnson', 'Bob Smith', 'TJSP', '2024-01-15'),
('6772130-04.2024.8.13.0161', 'Chris Miller', 'Jessica Taylor', 'TJSP', '2024-04-05'),
('3562061-02.2024.8.13.0431', 'Michael Brown', 'Sarah Davis', 'TJSP', '2024-02-20');

INSERT INTO case_update (cnj, update_date, update_details) VALUES
('5001682-88.2024.8.13.0672', '2024-07-31T10:00:00Z', 'Initial hearing scheduled for August 15, 2024.'),
('5001682-88.2024.8.13.0672', '2024-08-01T14:30:00Z', 'Plaintiff submitted additional evidence.'),
('5001682-88.2024.8.13.0672', '2024-08-02T09:00:00Z', 'Defendant requested a delay for response.'),
('6772130-04.2024.8.13.0161', '2024-07-30T11:00:00Z', 'Case file reviewed by judge.'),
('6772130-04.2024.8.13.0161', '2024-08-01T09:30:00Z', 'Witness statements collected.'),
('6772130-04.2024.8.13.0161', '2024-08-02T13:45:00Z', 'Defendant’s lawyer filed a motion for dismissal.'),
('6772130-04.2024.8.13.0161', '2024-08-03T10:00:00Z', 'Hearing date scheduled for August 10, 2024.');
