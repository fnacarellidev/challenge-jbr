CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE court_case (
	id				UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	cnj				TEXT NOT NULL,
	plaintiff		TEXT NOT NULL,
	defendant		TEXT NOT NULL,
	court_of_origin TEXT NOT NULL,
	start_date		DATE NOT NULL,
	created_at		TIMESTAMPTZ DEFAULT NOW()
);
