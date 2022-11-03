package bookmark

import (
	"github.com/labbs/alfred/pkg/exception"
)

type Bookmark struct {
	Id          string `gorm:"primaryKey" json:"id"`
	Name        string `json:"name,omitempty"`
	Url         string `json:"url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
	Tags        []*Tag `gorm:"many2many:bookmark_tags" json:"tags,omitempty"`

	UserId string `gorm:"index" json:"-"`
}

type Tag struct {
	Id   string `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`

	Bookmarks []*Bookmark `gorm:"many2many:bookmark_tags;" json:"bookmarks,omitempty"`

	UserId string `gorm:"index" json:"-"`
}

type BookmarkRepository interface {
	GetAllBookmarks(userId string) ([]Bookmark, *exception.AppError)
	GetBookmarkById(userId string, id string) (Bookmark, *exception.AppError)
	CreateBookmark(b Bookmark) *exception.AppError
	UpdateBookmark(b Bookmark) *exception.AppError
	DeleteBookmark(id string, userId string) *exception.AppError
	GetBookmarkByTag(userId string, tag string) ([]Bookmark, *exception.AppError)
	GetBookmarkByTags(userId string, tags []string) ([]Bookmark, *exception.AppError)
	FindBookmarkByWord(userId string, word string) ([]Bookmark, *exception.AppError)
	GetTags(userId string) ([]Tag, *exception.AppError)
	GetUniqueTags(userId string) ([]Tag, *exception.AppError)
	DeleteTag(id string, userId string) *exception.AppError
	GetTagByName(userId string, name string) (Tag, *exception.AppError)
	CreateTag(t Tag) *exception.AppError
	DeleteUnusedTag(b Bookmark, t Tag) *exception.AppError
}
