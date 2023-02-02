package models

import "gorm.io/gorm"

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title          string `json:"title"`
	Desc           string `json:"desc"`
	Content        string `json:"content"`
	CreatedBy      string `json:"created_by"`
	ModifiedBy     string `json:"modified_by"`
	State          int    `json:"state"`
	ConverImageUrl string `json:"conver_iamge_url"`
}

func CleanAllArticle() error {
	if err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{}).Error; err != nil {
		return err
	}
	return nil
}

func ExistArticleById(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetArticleTotal(maps map[string]interface{}) (int64, error) {
	var count int64
	if err := db.Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetArticles(pageNum, pageSize int, maps map[string]interface{}) ([]*Article, error) {
	var articles []*Article
	err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articles, nil
}

func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	err = db.Model(&article).Association("TagID").Find(&article.Tag)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func AddArticle(data map[string]interface{}) error {
	article := &Article{
		TagID:          data["tag_id"].(int),
		Title:          data["title"].(string),
		Desc:           data["desc"].(string),
		Content:        data["content"].(string),
		CreatedBy:      data["created_by"].(string),
		State:          data["state"].(int),
		ConverImageUrl: data["conver_image_url"].(string),
	}
	if err := db.Create(article).Error; err != nil {
		return err
	}

	return nil
}

func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(&Article{}).Error; err != nil {
		return err
	}
	return nil
}
