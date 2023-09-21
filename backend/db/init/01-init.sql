CREATE DATABASE IF NOT EXISTS exchange;
USE exchange;

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stocks (
    symbol VARCHAR(10) UNIQUE PRIMARY KEY NOT NULL,
    INDEX idx_stock_symbol(symbol)
);

CREATE TABLE if NOT EXISTS stock_price_history (
    symbol VARCHAR(10) NOT NULL,
    open DECIMAL(10, 2) NOT NULL,
    high DECIMAL(10, 2) NOT NULL,
    low DECIMAL(10, 2) NOT NULL,
    close DECIMAL(10, 2) NOT NULL,
    recorded_at BIGINT NOT NULL,
    FOREIGN KEY (symbol) REFERENCES Stocks(symbol),
    INDEX idx_stock_time(symbol, recorded_at DESC)
);

CREATE TABLE if NOT EXISTS orders (
    order_id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    order_side ENUM('Buy', 'Sell') NOT NULL,
    order_status ENUM('Open', 'Filled', 'PartiallyFilled', 'Rejected') NOT NULL,
    volume DECIMAL(10, 2) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (symbol) REFERENCES Stocks(symbol)
);