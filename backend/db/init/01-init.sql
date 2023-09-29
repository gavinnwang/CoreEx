CREATE DATABASE IF NOT EXISTS exchange;
USE exchange;

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    cash_balance DECIMAL(10, 2) NOT NULL DEFAULT 100000,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stocks (
    symbol VARCHAR(10) UNIQUE PRIMARY KEY NOT NULL
);

CREATE TABLE if NOT EXISTS stock_history (
    symbol VARCHAR(10) NOT NULL,
    open DECIMAL(10, 2) NOT NULL,
    high DECIMAL(10, 2) NOT NULL,
    low DECIMAL(10, 2) NOT NULL,
    close DECIMAL(10, 2) NOT NULL,
    bid_volume DECIMAL(10, 2) NOT NULL,
    ask_volume DECIMAL(10, 2) NOT NULL,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (symbol) REFERENCES stocks(symbol),
    INDEX idx_stock_time(symbol, recorded_at DESC)
);

CREATE TABLE if NOT EXISTS orders (
    order_id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    order_side ENUM('Buy', 'Sell') NOT NULL,
    order_status ENUM('Open', 'Filled', 'PartiallyFilled', 'Rejected') NOT NULL,
    order_type ENUM('Market', 'Limit') NOT NULL,
    filled_at DECIMAL(10, 2),
    total_processed DECIMAL(10, 2) DEFAULT 0,
    volume DECIMAL(10, 2) NOT NULL,
    initial_volume DECIMAL(10, 2) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (symbol) REFERENCES stocks(symbol)
);

CREATE TABLE if NOT EXISTS holdings (
    user_id VARCHAR(26) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    volume DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, symbol), -- ensures that a user can only have one holding of a stock, combination is unique
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (symbol) REFERENCES stocks(symbol)
);

DELIMITER //
CREATE PROCEDURE InsertOrUpdateHoldingThenDeleteZeroVolume(
    IN p_user_id VARCHAR(26), 
    IN p_symbol VARCHAR(10), 
    IN p_volume DECIMAL(10, 2)
)
BEGIN
    DECLARE updated_volume DECIMAL(10, 2);

    -- Insert or update
    INSERT INTO holdings (user_id, symbol, volume)
    VALUES (p_user_id, p_symbol, p_volume)
    ON DUPLICATE KEY UPDATE 
    volume = volume + p_volume;

    SELECT volume INTO updated_volume
    FROM holdings
    WHERE user_id = p_user_id AND symbol = p_symbol;

    -- Delete rows with volume of zero for the specific user and symbol
    IF updated_volume = 0 THEN
        DELETE FROM holdings WHERE user_id = p_user_id AND symbol = p_symbol;
    END IF;
END //
DELIMITER ;

