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

func (fileDao *FileDao) QueryFilesByUserId(parentId string, userId, page, pageSize int) (result []*model.BeeFile, err error) {

	err = fileDao.mysql.Select(&result, "SELECT `parent_id`,`file_id`,`user_id`,`file_name`,`file_size`,`file_type`,`create_time`,`update_time` FROM `bee_file` WHERE `user_id`=? AND `parent_id`=? AND `status`=1 ORDER BY `create_time` DESC LIMIT ?, ?", userId, parentId, page, pageSize)

	return
}
