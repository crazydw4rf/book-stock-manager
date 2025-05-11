package repository

const (
	bookCreate        = `INSERT INTO book(book_id,isbn,title,author,publisher,published_at,stock) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *`
	bookGetById       = `SELECT * FROM book WHERE book_id = $1 LIMIT 1`
	bookGetByISBN     = `SELECT * FROM book WHERE isbn = $1 LIMIT 1`
	bookGetBooksMany  = `SELECT * FROM book OFFSET $1 LIMIT $2`
	bookDelete        = `DELETE FROM book WHERE book_id = $1 RETURNING book_id`
	bookGetTotalCount = `SELECT COUNT(*) FROM book`
	bookUpdate        = `UPDATE book SET
isbn = COALESCE(NULLIF($2, ''), isbn),
title = COALESCE(NULLIF($3, ''), title),
author = COALESCE(NULLIF($4, ''), author),
publisher = COALESCE(NULLIF($5, ''), publisher),
published_at = COALESCE($6, published_at),
stock = COALESCE($7, stock)
WHERE book_id = $1 RETURNING *`
)
