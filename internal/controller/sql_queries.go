package controller

const (
	createBook = `INSERT INTO book(book_id,isbn,title,author,publisher,published_at,stock) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *`
)
