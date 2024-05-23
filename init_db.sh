#!/bin/sh
DB_PATH=/data/exchange_rates.db

sqlite3 $DB_PATH <<EOF
CREATE TABLE IF NOT EXISTS exchange_rates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT,
    codein TEXT,
    name TEXT,
    high TEXT,
    low TEXT,
    varBid TEXT,
    pctChange TEXT,
    bid TEXT,
    ask TEXT,
    timestamp TEXT,
    create_date TEXT
);
EOF

echo "Database and table initialization complete."
