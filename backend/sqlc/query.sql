-- name: InsertCourtCase :one
INSERT INTO court_case (
	cnj, plaintiff, defendant, court_of_origin, start_date
) VALUES (
	$1, $2, $3, $4, $5
)
RETURNING id;
