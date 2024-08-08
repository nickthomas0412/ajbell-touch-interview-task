-- Create the Client table
CREATE TABLE IF NOT EXISTS Clients (
    client_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Create the Deposit table
CREATE TABLE IF NOT EXISTS Deposits (
    deposit_id TEXT PRIMARY KEY,
    client_id TEXT,
    nominal REAL NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES Clients(client_id)
);

-- Create the Receipt table
CREATE TABLE IF NOT EXISTS Receipts (
    receipt_id TEXT PRIMARY KEY,
    deposit_id TEXT,
    amount REAL NOT NULL,
    pot TEXT NOT NULL,
    wrapper_type TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (deposit_id) REFERENCES Deposits(deposit_id)
);
