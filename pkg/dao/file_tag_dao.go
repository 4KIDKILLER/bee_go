package dao

import (
	"database/sql"
	"goserver/pkg/model"

	"github.com/jmoiron/sqlx"
)

type FileTagDao struct {
	mysql *sqlx.DB
}

func NewFileTagDao(mysql *sqlx.DB) (fileTagDao *FileTagDao) {
	fileTagDao = &FileTagDao{mysql}
	return
}

func (fileTagDao *FileTagDao) Insert(tagName, tagId, fileId string, userId int) (result sql.Result, err error) {
	result, err = fileTagDao.mysql.Exec("INSERT INTO bee_file_tag (tag_id,file_id,user_id,tag_name) VALUES (?,?,?,?)", tagId, fileId, userId, tagName)

	return
}

func (fileTagDao *FileTagDao) UpdateTagByTagId(tagName, tagId string) (result sql.Result, err error) {
	result, err = fileTagDao.mysql.Exec("UPDATE bee_file_tag SET tag_name=? WHERE tag_id=?", tagName, tagId)

	return
}

func (fileTagDao *FileTagDao) DeleteRowByTagId(tagId string) (result sql.Result, err error) {
	result, err = fileTagDao.mysql.Exec("DELETE FROM bee_file_tag WHERE tag_id=?", tagId)

	return
}

func (fileTagDao *FileTagDao) SelectTagsByFileId(fileId string) (result []*model.BeeFileTag, err error) {
	err = fileTagDao.mysql.Select(result, "SELECT tag_name FROM bee_file_tag WHERE file_id=?", fileId)
	return
}

func (fileTagDao *FileTagDao) SelectTagsByFileIds(fileIds []string, userId int) (result []*model.BeeFileTag, err error) {
	query, args, err := sqlx.In("SELECT tag_id,file_id,tag_name FROM bee_file_tag WHERE user_id=? AND file_id IN (?)", userId, fileIds)

	if err != nil {
		return nil, err
	}

	query = fileTagDao.mysql.Rebind(query)

	err = fileTagDao.mysql.Select(&result, query, args...)
	if err != nil {
		return nil, err
	}

	return
}
