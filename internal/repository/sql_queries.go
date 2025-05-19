package repository

const (
	bookCreate        = `INSERT INTO books(book_id,isbn,title,author,publisher,published_at,stock) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *`
	bookGetById       = `SELECT * FROM books WHERE book_id = $1 LIMIT 1`
	bookGetByISBN     = `SELECT * FROM books WHERE isbn = $1 LIMIT 1`
	bookGetBooksMany  = `SELECT * FROM books OFFSET $1 LIMIT $2`
	bookDelete        = `DELETE FROM books WHERE book_id = $1 RETURNING book_id`
	bookGetTotalCount = `SELECT COUNT(*) FROM books`
	bookUpdate        = `UPDATE books SET
isbn = COALESCE(NULLIF($2, ''), isbn),
title = COALESCE(NULLIF($3, ''), title),
author = COALESCE(NULLIF($4, ''), author),
publisher = COALESCE(NULLIF($5, ''), publisher),
published_at = COALESCE(NULLIF($6, '0001-01-01'::date), published_at),
stock = CASE WHEN $7 < 0 THEN stock ELSE $7 END,
updated_at = NOW() WHERE book_id = $1 RETURNING *`
)
