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

