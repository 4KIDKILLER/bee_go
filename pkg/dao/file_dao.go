package dao

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type FileDao struct {
	mysql *sqlx.DB
}

func NewFileDao(mysql *sqlx.DB) (fileDao *FileDao) {
	fileDao = &FileDao{mysql}
	return
}

func (fileDao *FileDao) Insert(parentId, fileId, fileName, filePath string, userId, fileSize, fileType int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("INSERT INTO `bee_file` (`parent_id`,`file_id`,`user_id`,`file_name`,`file_size`,`file_path`,`file_type`) VALUES (?,?,?,?,?,?,?,?)")
	return
}
