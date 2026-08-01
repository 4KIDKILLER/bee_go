package dao

import (
	"database/sql"
	"goserver/pkg/model"

	"github.com/jmoiron/sqlx"
)

type FileDao struct {
	mysql *sqlx.DB
}

func NewFileDao(mysql *sqlx.DB) (fileDao *FileDao) {
	fileDao = &FileDao{mysql}
	return
}

func (fileDao *FileDao) Insert(parentId, fileId, fileOriginalName, fileExt, filePath, tags, cover1, cover2, cover3, remark string, fileSize int64, userId, fileType int) (result sql.Result, err error) {

	result, err = fileDao.mysql.Exec("INSERT INTO `bee_file` (`parent_id`,`file_id`,`user_id`,`file_original_name`,`file_ext`,`file_size`,`file_path`,`file_type`,`tags`,`cover_1`,`cover_2`,`cover_3`,`remark`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", parentId, fileId, userId, fileOriginalName, fileExt, fileSize, filePath, fileType, tags, cover1, cover2, cover3, remark)
	return
}

func (fileDao *FileDao) CountFileByUserIdAndParentId(userId int, parentId string) (count int, err error) {
	err = fileDao.mysql.Get(&count, "SELECT COUNT(`file_id`) FROM `bee_file` WHERE `user_id`=? AND `parent_id`=? AND `status`=1", userId, parentId)
	return
}

func (fileDao *FileDao) QueryFilesByUserId(parentId string, userId, page, pageSize int) (result []*model.BeeFile, err error) {

	err = fileDao.mysql.Select(&result, "SELECT `parent_id`,`file_id`,`user_id`,`file_original_name`,`file_ext`,`file_path`,`file_size`,`file_type`,`tags`,`cover_1`,`cover_2`,`cover_3`,`remark`,`create_time`,`update_time` FROM `bee_file` WHERE `user_id`=? AND `parent_id`=? AND `status`=1 ORDER BY `create_time` DESC LIMIT ?, ?", userId, parentId, page, pageSize)

	return
}
