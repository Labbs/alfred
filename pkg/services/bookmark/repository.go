package bookmark

import (
	"github.com/labbs/alfred/pkg/database"
	"github.com/labbs/alfred/pkg/exception"
)

type BookmarkRepositoryDB struct {
	client database.DbConnection
}

func NewBookmarkRepository() BookmarkRepositoryDB {
	client := database.GetDbConnection()
	return BookmarkRepositoryDB{client: client}
}

func (d BookmarkRepositoryDB) GetAllBookmarks(userId string) ([]Bookmark, *exception.AppError) {
	var b []Bookmark
	r := d.client.DB.
		Preload("Tags").
		Where("user_id = ?", userId).Find(&b)
	if r.Error != nil {
		return []Bookmark{}, exception.NewUnexpectedError("unable to find bookmark(s)", r.Error)
	}
	return b, nil
}

func (d BookmarkRepositoryDB) GetBookmarkById(userId string, id string) (Bookmark, *exception.AppError) {
	b := Bookmark{}
	r := d.client.DB.
		Preload("Tags").
		Where("id = ? and user_id = ?", id, userId).First(&b)
	if r.Error != nil {
		return Bookmark{}, exception.NewUnexpectedError("unable to find bookmark", r.Error)
	}
	return b, nil
}

func (d BookmarkRepositoryDB) CreateBookmark(b Bookmark) *exception.AppError {
	r := d.client.DB.Create(&b)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to create bookmark", r.Error)
	}
	return nil
}

func (d BookmarkRepositoryDB) UpdateBookmark(b Bookmark) *exception.AppError {
	r := d.client.DB.Save(&b)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to update bookmark", r.Error)
	}
	return nil
}

func (d BookmarkRepositoryDB) DeleteBookmark(id string, userId string) *exception.AppError {
	r := d.client.DB.Where("id = ? and user_id = ?", id, userId).Delete(&Bookmark{})
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to delete bookmark", r.Error)
	}
	r = d.client.DB.Where("bookmark_id = ?", id).Delete(&Tag{})
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to delete bookmark tags", r.Error)
	}
	return nil
}

func (d BookmarkRepositoryDB) GetBookmarkByTag(userId string, tag string) ([]Bookmark, *exception.AppError) {
	var b []Bookmark
	r := d.client.DB.
		Preload("Tags").
		Where("user_id = ? and tags.name = ?", userId, tag).Find(&b)
	if r.Error != nil {
		return []Bookmark{}, exception.NewUnexpectedError("unable to find bookmark(s)", r.Error)
	}
	return b, nil
}

func (d BookmarkRepositoryDB) GetBookmarkByTags(userId string, tags []string) ([]Bookmark, *exception.AppError) {
	var b []Bookmark
	r := d.client.DB.
		Preload("Tags").
		Where("user_id = ? and tags.name in ?", userId, tags).Find(&b)
	if r.Error != nil {
		return []Bookmark{}, exception.NewUnexpectedError("unable to find bookmark(s)", r.Error)
	}
	return b, nil
}

func (d BookmarkRepositoryDB) FindBookmarkByWord(userId string, word string) ([]Bookmark, *exception.AppError) {
	var b []Bookmark
	r := d.client.DB.
		Preload("Tags").
		Where("user_id = ? and (title like ? or description like ?)", userId, "%"+word+"%", "%"+word+"%").Find(&b)
	if r.Error != nil {
		return []Bookmark{}, exception.NewUnexpectedError("unable to find bookmark(s)", r.Error)
	}
	return b, nil
}

func (d BookmarkRepositoryDB) GetTags(userId string) ([]Tag, *exception.AppError) {
	var t []Tag
	r := d.client.DB.
		Where("user_id = ?", userId).Find(&t)
	if r.Error != nil {
		return []Tag{}, exception.NewUnexpectedError("unable to find tag(s)", r.Error)
	}
	return t, nil
}

func (d BookmarkRepositoryDB) GetUniqueTags(userId string) ([]Tag, *exception.AppError) {
	var t []Tag
	r := d.client.DB.
		Where("user_id = ?", userId).Group("name").Find(&t)
	if r.Error != nil {
		return []Tag{}, exception.NewUnexpectedError("unable to find tag(s)", r.Error)
	}
	return t, nil
}

func (d BookmarkRepositoryDB) DeleteTag(id string, userId string) *exception.AppError {
	r := d.client.DB.Where("id = ? and user_id = ?", id, userId).Delete(&Tag{})
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to delete tag", r.Error)
	}
	return nil
}
