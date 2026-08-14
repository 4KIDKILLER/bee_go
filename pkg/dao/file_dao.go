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

func (fileDao *FileDao) Insert(parentId, fileId, fileOriginalName, fileExt, filePath, fileThumbPath, tags, cover1, cover2, cover3, remark string, fileSize int64, userId, fileType int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("INSERT INTO `bee_file` (`parent_id`,`file_id`,`user_id`,`file_original_name`,`file_ext`,`file_size`,`file_path`,`file_thumb_path`,`file_type`,`tags`,`cover_1`,`cover_2`,`cover_3`,`remark`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", parentId, fileId, userId, fileOriginalName, fileExt, fileSize, filePath, fileThumbPath, fileType, tags, cover1, cover2, cover3, remark)
	return
}

func (fileDao *FileDao) CountFileByParentId(userId int, parentId string) (count int, err error) {
	err = fileDao.mysql.Get(&count, "SELECT COUNT(`file_id`) FROM `bee_file` WHERE `user_id`=? AND `parent_id`=? AND `status`=1", userId, parentId)
	return
}

func (fileDao *FileDao) QueryUserFiles(parentId string, userId, page, pageSize int) (result []*model.BeeFile, err error) {
	err = fileDao.mysql.Select(&result, "SELECT `parent_id`,`file_id`,`user_id`,`file_original_name`,`file_ext`,`file_path`,`file_thumb_path`,`file_size`,`file_type`,`tags`,`cover_1`,`cover_2`,`cover_3`,`remark`,`create_time`,`update_time` FROM `bee_file` WHERE `user_id`=? AND `parent_id`=? AND `status`=1 ORDER BY `id` DESC LIMIT ?, ?", userId, parentId, page, pageSize)

	return
}

func (fileDao *FileDao) UpdateThumbPathByFileId(thumbPath, fileId string, userId int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("UPDATE `bee_file` SET `file_thumb_path`=? WHERE `file_id`=? AND `user_id`=?", thumbPath, fileId, userId)
	return
}

func (fileDao *FileDao) UploadOriginalNameByFileId(name, fileId string, userId int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("UPDATE `bee_file` SET `file_original_name`=? WHERE `file_id`=? AND `user_id`=?", name, fileId, userId)
	return
}

func (fileDao *FileDao) UpdateStatusByFileIdAndFileType(fileId string, userId, fileType, status int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("UPDATE `bee_file` SET `status`=? WHERE `file_id`=? AND `user_id`=? AND `file_type`=?", status, fileId, userId, fileType)
	return
}

func (fileDao *FileDao) QueryUserFolders(userId int) (result []*model.BeeFile, err error) {
	err = fileDao.mysql.Select(&result, "SELECT `parent_id`,`file_id`,`user_id`,`file_original_name` FROM `bee_file` WHERE `user_id`=? AND `file_type`=1 AND `status`=1", userId)

	return
}
