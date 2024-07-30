CREATE TABLE court_case (
	cnj				TEXT NOT NULL,
	plaintiff		TEXT NOT NULL,
	defendant		TEXT NOT NULL,
	court_of_origin TEXT NOT NULL,
	start_date		TIMESTAMPTZ NOT NULL,
	created_at		TIMESTAMPTZ NOT NULL
);
