package models

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

func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{})
	return true
}

func ExistArticleById(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)

	return article.ID > 0
}

func GetArticleTotal(maps map[string]interface{}) (count int64) {
	db.Where(maps).Count(&count)
	return
}

func GetArticles(pageNum, pageSize int, maps map[string]interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Association("TagID").Find(&article.Tag)

	return
}

func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:          data["tag_id"].(int),
		Title:          data["title"].(string),
		Desc:           data["desc"].(string),
		Content:        data["content"].(string),
		CreatedBy:      data["created_by"].(string),
		State:          data["state"].(int),
		ConverImageUrl: data["conver_image_url"].(string),
	})

	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(&Article{})
	return true
}
