package dao

import (
	"database/sql"
	"fmt"
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

func (fileDao *FileDao) Insert(parentId, fileId, fileOriginalName, fileExt, filePath, fileThumbPath, cover1, cover2, cover3, remark string, fileSize int64, userId, fileType, isPrivate int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("INSERT INTO bee_file (parent_id,file_id,user_id,file_original_name,file_ext,file_size,file_path,file_thumb_path,file_type,cover_1,cover_2,cover_3,remark,is_private) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", parentId, fileId, userId, fileOriginalName, fileExt, fileSize, filePath, fileThumbPath, fileType, cover1, cover2, cover3, remark, isPrivate)
	return
}

func (fileDao *FileDao) CountRowByParentId(userId int, parentId string, isPrivates []int) (count int64, err error) {

	query, args, err := sqlx.In("SELECT COUNT(*) FROM bee_file WHERE user_id=? AND parent_id=? AND is_private IN (?) AND status=1", userId, parentId, isPrivates)
	if err != nil {
		return 0, err
	}

	query = fileDao.mysql.Rebind(query)

	err = fileDao.mysql.Get(&count, query, args...)
	return
}

func (fileDao *FileDao) SelectRowsLimitByUserId(parentId string, userId int, page, pageSize int64, isPrivates []int) (result []*model.BeeFile, err error) {
	query, args, err := sqlx.In("SELECT parent_id,file_id,user_id,file_original_name,file_ext,file_path,file_thumb_path,file_size,file_type,cover_1,cover_2,cover_3,remark,create_time,update_time,is_private,status FROM bee_file WHERE user_id=? AND parent_id=? AND is_private IN (?) AND status=1 ORDER BY id DESC LIMIT ?, ?", userId, parentId, isPrivates, page, pageSize)

	if err != nil {
		return nil, err
	}

	query = fileDao.mysql.Rebind(query)

	err = fileDao.mysql.Select(&result, query, args...)
	return
}

func (fileDao *FileDao) UpdateRowThumbPathByFileId(thumbPath, fileId string, userId int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("UPDATE bee_file SET file_thumb_path=? WHERE file_id=? AND user_id=?", thumbPath, fileId, userId)
	return
}

func (fileDao *FileDao) UploadRowOriginalNameByFileId(name, fileId string, userId int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("UPDATE bee_file SET file_original_name=? WHERE file_id=? AND user_id=?", name, fileId, userId)
	return
}

func (fileDao *FileDao) UpdateRowStatusByFileIdAndFileType(fileId string, userId, fileType, status int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("UPDATE bee_file SET status=? WHERE file_id=? AND user_id=? AND file_type=?", status, fileId, userId, fileType)
	return
}

func (fileDao *FileDao) SelectRowsByUserId(userId int) (result []*model.BeeFile, err error) {
	err = fileDao.mysql.Select(&result, "SELECT parent_id,file_id,user_id,file_original_name FROM bee_file WHERE user_id=? AND file_type=1 AND status=1", userId)
	return
}

func (fileDao *FileDao) SelectRowsRecursionByFileId(fileId string, userId int) (result []string, err error) {
	//这里的result返回的类型是[]string,使用select需要注意只能有一列，
	//即select file_id，不能多列select file_id,file_type
	if fileId != "" {
		/*
			如果fileId存在，则查询当前文件夹下的所有文件/文件夹
		*/
		err = fileDao.mysql.Select(&result, `WITH RECURSIVE file_tree AS (
						SELECT file_id,file_type,file_original_name FROM bee_file WHERE file_id=? AND user_id=?
						UNION ALL
						SELECT f.file_id,f.file_type,f.file_original_name FROM bee_file f INNER JOIN file_tree t ON f.parent_id = t.file_id
						)
						SELECT file_id FROM file_tree`,
			fileId,
			userId,
		)
	} else {
		/*
			否则查询当前用户下的所有正常状态的文件/文件夹
		*/
		err = fileDao.mysql.Select(&result, "SELECT file_id FROM bee_file WHERE user_id=? AND `status`=1",
			userId,
		)
	}

	return
}

func (fileDao *FileDao) UpdateRowsStatusByFileIdInIds(fileIds []string, status, userId int) (result int64, err error) {
	/*
		使用sqlx.In函数构建in语法语句.这里会将IN (?)构建为fileIds一样长度的模版字符串
		例如fileIds中长度为3,则构建的语法中IN (?)会变成IN (?,?,?),方便后续执行进行赋值
		其中：query为构建后的sql模板，args为参数切片
	*/
	query, args, err := sqlx.In("UPDATE bee_file SET status=? WHERE file_id IN (?) AND user_id=?", status, fileIds, userId)
	if err != nil {
		return 0, err
	}

	query = fileDao.mysql.Rebind(query)

	execRes, err := fileDao.mysql.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	//获取受影响的行数
	result, err = execRes.RowsAffected()
	if err != nil {
		return 0, err
	}
	return
}

func (fileDao *FileDao) SelectRowsByStatusOrderByFileType(status, userId int) (result []*model.BeeFile, err error) {
	err = fileDao.mysql.Select(&result, "SELECT file_id,file_type,file_ext,file_path,file_thumb_path FROM bee_file WHERE status=? AND user_id=? ORDER BY file_type ASC", status, userId)
	return
}

func (fileDao *FileDao) DeleteRowsByStatus(status int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("DELETE FROM bee_file WHERE status=?", status)
	return
}

func (fileDao *FileDao) UpdateRemarkByFileId(fileId, remark string, userId int) (result sql.Result, err error) {
	result, err = fileDao.mysql.Exec("UPDATE bee_file SET remark=? WHERE file_id=? AND user_id=?", remark, fileId, userId)
	return
}

func (fileDao *FileDao) UpdateCoverByFileId(cover, fileId string, position, userId int) (result sql.Result, err error) {

	query := fmt.Sprintf("UPDATE bee_file SET cover_%d=? WHERE file_id=? AND user_id=?", position)
	result, err = fileDao.mysql.Exec(query, cover, fileId, userId)
	return
}

func (fileDao *FileDao) UpdateRowsIsPrivateByFileIds(fileIds []string, userId int) (result int64, err error) {
	query, args, err := sqlx.In(`
		WITH RECURSIVE file_tree AS (
			SELECT file_id,file_type 
			FROM bee_file
			WHERE user_id=? AND file_id IN (?)

			UNION DISTINCT

			SELECT child.file_id,child.file_type 
			FROM bee_file AS child
			INNER JOIN file_tree AS parent
				ON child.parent_id=parent.file_id
			WHERE child.user_id=? AND parent.file_type=1
		)
		UPDATE bee_file AS file 
		INNER JOIN file_tree AS tree
			ON file.file_id=true.file_id
		SET file.is_private=2
		WHERE file.user_id=?
	`, userId, fileIds, userId, userId)

	if err != nil {
		return 0, err
	}

	query = fileDao.mysql.Rebind(query)

	execRes, err := fileDao.mysql.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	result, err = execRes.RowsAffected()
	return
}
