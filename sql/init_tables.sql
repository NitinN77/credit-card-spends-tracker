CREATE TABLE IF NOT EXISTS cc_transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL,
    card_name TEXT,
    last_4 TEXT,
    merchant TEXT,
    date TEXT
);

CREATE TABLE IF NOT EXISTS fetched_date_ranges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_date TEXT,
    end_date TEXT
);

CREATE TABLE IF NOT EXISTS merchant_aliases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias TEXT,
    merchant TEXT
);

