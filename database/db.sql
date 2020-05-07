CREATE TABLE `speedtest` (
    `id`         INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `dt`         VARCHAR(20),
    `latency`    REAL,
    `jitter`     REAL,
    `download`   REAL NOT NULL,
    `upload`     REAL NOT NULL,
    `packetLoss` INT,
    `hardware`   VARCHAR(10) DEFAULT "not found",
    `serverId`   INT,
    `port`       INT,
    `ip`         VARCHAR(15),
    `name`       VARCHAR(50),
    `location`  VARCHAR(30),
    `host`       VARCHAR(50),
    PRIMARY KEY(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;