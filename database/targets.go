package database

import (
	"fmt"
	"github.com/jvzantvoort/xarchive/target"
	log "github.com/sirupsen/logrus"
)

const (
	SQLString string = "SELECT `target_id`, `target_path`, `target_name`, `target_sha256`, `size_b` FROM `targets`"
)

type Record struct {
	ID   int64
	Path string
	Name string
	Hash string
	Size int64
}

func (r Record) GetDescr() string {
	return fmt.Sprintf("%s:%s", r.Name, r.Hash)
}

// targets table
// +---------------+--------------+------+-----+---------+----------------+
// | Field         | Type         | Null | Key | Default | Extra          |
// +---------------+--------------+------+-----+---------+----------------+
// | target_id     | int(128)     | NO   | PRI | NULL    | auto_increment |
// | target_path   | varchar(512) | NO   | UNI | NULL    |                |
// | target_name   | varchar(128) | NO   |     | NULL    |                |
// | target_sha256 | varchar(64)  | NO   |     | NULL    |                |
// | size_b        | int(32)      | YES  |     | NULL    |                |
// +---------------+--------------+------+-----+---------+----------------+

func (d Database) Query(sqlstr string, args ...any) ([]Record, error) {
	retv := []Record{}

	results, err := d.Connection.Query(sqlstr, args...)
	if err != nil {
		return retv, err
	}

	for results.Next() {
		var rec Record
		err := results.Scan(&rec.ID,
			&rec.Path,
			&rec.Name,
			&rec.Hash,
			&rec.Size)
		if err != nil {
			return retv, err
		}
		retv = append(retv, rec)
	}
	return retv, nil
}

func (d Database) GetTargets(filename, checksum string) ([]Record, error) {
	retv := []Record{}
	sqlstr := SQLString + " WHERE `target_name` = ? AND `target_sha256` = ?"

	results, err := d.Connection.Query(sqlstr, filename, checksum)
	if err != nil {
		return retv, err
	}

	for results.Next() {
		var rec Record
		err := results.Scan(&rec.ID,
			&rec.Path,
			&rec.Name,
			&rec.Hash,
			&rec.Size)
		if err != nil {
			return retv, err
		}
		retv = append(retv, rec)
	}
	return retv, nil
}

func (d Database) LookupTarget(path string) ([]Record, error) {
	sqlstr := SQLString + " WHERE `target_path` = ?"
	args := []any{}

	args = append(args, path)
	return d.Query(sqlstr, args...)
}

func (d Database) InsertTarget(obj *target.Target) (int64, error) {
	log.Debugf("Insert %s, start", obj.Path)
	defer log.Debugf("Insert %s, end", obj.Path)

	sqlstr := "INSERT INTO `targets` ("
	sqlstr += "  `target_path`,"
	sqlstr += "  `target_name`,"
	sqlstr += "  `target_sha256`,"
	sqlstr += "  `size_b` ) VALUES ( ?, ?, ?, ?)"
	insert, err := d.Connection.Prepare(sqlstr)
	if err != nil {
		return 0, err
	}

	result, err := insert.Exec(obj.Path, obj.Name, obj.Hash, obj.Size)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
