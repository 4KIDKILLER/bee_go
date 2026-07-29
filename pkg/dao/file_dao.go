package dao

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type FileDao struct {
	mysql *sqlx.DB
}

func NewFileDao(mysql *sqlx.DB) (fileDao *FileDao) {
	fileDao = &FileDao{mysql}
	return
}

func (fileDao *FileDao) Insert(parentId, fileId, fileName, filePath string, userId, fileType int, fileSize decimal.Decimal) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("INSERT INTO `bee_file` (`parent_id`,`file_id`,`user_id`,`file_name`,`file_size`,`file_path`,`file_type`) VALUES (?,?,?,?,?,?,?)", parentId, fileId, userId, fileName, fileSize, filePath, fileType)
	return
}
