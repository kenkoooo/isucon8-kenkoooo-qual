CREATE TABLE IF NOT EXISTS users (
    id          INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    nickname    VARCHAR(128) NOT NULL,
    login_name  VARCHAR(128) NOT NULL,
    pass_hash   VARCHAR(128) NOT NULL,
    UNIQUE KEY login_name_uniq (login_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS events (
    id          INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    title       VARCHAR(128)     NOT NULL,
    public_fg   TINYINT(1)       NOT NULL,
    closed_fg   TINYINT(1)       NOT NULL,
    price       INTEGER UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sheets (
    id          INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `rank`      VARCHAR(128)     NOT NULL,
    num         INTEGER UNSIGNED NOT NULL,
    price       INTEGER UNSIGNED NOT NULL,
    UNIQUE KEY rank_num_uniq (`rank`, num)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS reservations (
    id          INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    event_id    INTEGER UNSIGNED NOT NULL,
    sheet_id    INTEGER UNSIGNED NOT NULL,
    user_id     INTEGER UNSIGNED NOT NULL,
    reserved_at DATETIME(6)      NOT NULL,
    canceled_at DATETIME(6)      DEFAULT NULL,
    KEY event_id_and_sheet_id_idx (event_id, sheet_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE INDEX idx_reservations_user_id ON reservations (user_id);
CREATE INDEX idx_reservations_event_id ON reservations (event_id);
CREATE INDEX idx_reservations_sheet_id ON reservations (sheet_id);

CREATE TABLE IF NOT EXISTS administrators (
    id          INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    nickname    VARCHAR(128) NOT NULL,
    login_name  VARCHAR(128) NOT NULL,
    pass_hash   VARCHAR(128) NOT NULL,
    UNIQUE KEY login_name_uniq (login_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sold (
    event_id    INTEGER         NOT NULL,
    sheet_rank  VARCHAR(128)    NOT NULL,
    sold_count  INTEGER         NOT NULL,
    PRIMARY KEY pkey_event_id_sheet_rank (event_id, sheet_rank)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DELIMITER //
CREATE TRIGGER sold_insert_reservation
AFTER INSERT ON reservations FOR EACH ROW
BEGIN
    IF NEW.canceled_at IS NULL THEN
        INSERT INTO sold (event_id, sheet_rank, sold_count)
        SELECT NEW.event_id, s.`rank`, 1 FROM sheets s WHERE s.id=NEW.sheet_id
        ON DUPLICATE KEY UPDATE sold_count = sold_count+1;
    END IF;
END; //

DELIMITER //
CREATE TRIGGER sold_update_reservation
AFTER UPDATE ON reservations FOR EACH ROW
BEGIN
    IF NEW.canceled_at IS NOT NULL THEN
        INSERT INTO sold (event_id, sheet_rank, sold_count)
        SELECT NEW.event_id, s.`rank`, -1 FROM sheets s WHERE s.id=NEW.sheet_id
        ON DUPLICATE KEY UPDATE sold_count = sold_count-1;
    END IF;
END; //
