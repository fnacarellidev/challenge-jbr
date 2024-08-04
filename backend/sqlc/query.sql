-- name: InsertCourtCase :one
INSERT INTO court_case (
	cnj, plaintiff, defendant, court_of_origin, start_date
) VALUES (
	$1, $2, $3, $4, $5
)
RETURNING cnj;

-- name: GetCourtCase :one
SELECT * from court_case WHERE cnj = $1 LIMIT 1;

-- name: GetCaseUpdates :many
SELECT update_date, update_details FROM case_update WHERE cnj = $1 ORDER BY update_date DESC;

-- name: InsertCaseUpdate :exec
INSERT INTO case_update (
	cnj, update_date, update_details
) VALUES (
	$1, $2, $3
);
