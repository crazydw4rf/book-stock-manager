CREATE TABLE IF NOT EXISTS book (
    book_id UUID PRIMARY KEY,
    isbn CHAR(17) NOT NULL,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    publisher TEXT NOT NULL,
    published_at DATE NOT NULL,
    stock BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

