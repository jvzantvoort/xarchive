# xarchive

This is a simple tool to maintain larger amounts of binary files for
my hugo based website.

configfile: ``~/.xarchive.yaml``

```yaml
---
database:
  username: targets
  password: targets
  hostname: localhost
  database: targets
  port: 3306
```

## Database schema

```sql
CREATE TABLE `config` (
  `config_name` varchar(128) NOT NULL,
  `config_value` varchar(512) NOT NULL,
   PRIMARY KEY (`config_name`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `targets` (
  `target_id` int(128) NOT NULL AUTO_INCREMENT,
  `target_path` varchar(512) NOT NULL,
  `target_name` varchar(128) NOT NULL,
  `target_sha256` varchar(64) NOT NULL,
  `size_b` int(32) DEFAULT NULL,
  PRIMARY KEY (`target_id`),
  UNIQUE KEY `target_path` (`target_path`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
```
