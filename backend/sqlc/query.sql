-- name: InsertCourtCase :one
INSERT INTO court_case (
	cnj, plaintiff, defendant, court_of_origin, start_date
) VALUES (
	$1, $2, $3, $4, $5
)
RETURNING id;

-- name: GetCourtCase :one
SELECT * from court_case WHERE cnj = $1 LIMIT 1;
