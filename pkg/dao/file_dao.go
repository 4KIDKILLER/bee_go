package dao

import (
	"database/sql"
	"goserver/pkg/model"

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

func (fileDao *FileDao) QueryFileByName(FileName string) (err error) {
	var folder model.BeeFile
	err = fileDao.mysql.Get(&folder, "SELECT `parent_id`,`file_id`,`user_id`,`file_name`,`file_size`,`file_path`,`file_type`,`create_time`,`update_time` WHERE `file_name`=? AND `file_type`=2 AND `status`=1", FileName)
	return
}
