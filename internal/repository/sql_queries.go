package repository

// TODO: test lagi query update nya
const (
	bookCreate       = `INSERT INTO book(book_id,isbn,title,author,publisher,published_at,stock) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *`
	bookGetById      = `SELECT * FROM book WHERE book_id = $1 LIMIT 1`
	bookGetByISBN    = `SELECT * FROM book WHERE isbn = $1 LIMIT 1`
	bookGetBooksMany = `SELECT * FROM book OFFSET $1 LIMIT $2`
	bookDelete       = `DELETE FROM book WHERE book_id = $1 RETURNING book_id`
	bookUpdate       = `UPDATE book SET
isbn = COALESCE(NULLIF($2, ''), isbn),
title = COALESCE(NULLIF($3, ''), title),
author = COALESCE(NULLIF($4, ''), author),
publisher = COALESCE(NULLIF($5, ''), publisher),
published_at = COALESCE(NULLIF($6, '0001-1-1'::date), published_at),
stock = COALESCE(NULLIF($7, -1), stock) WHERE book_id = $1 RETURNING *`
)
