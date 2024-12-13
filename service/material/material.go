package material

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/SKAT/models"
)

func addMaterial(c *gin.Context) {
	var material models.Material

	// Чтение данных из запроса
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Чтение файла из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Открытие файла
	data, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer data.Close()

	// Чтение содержимого файла
	fileBytes, err := ioutil.ReadAll(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
		return
	}

	material.File = fileBytes

	mat := map[string]interface{}{
		"id":      material.ID,
		"name":    material.Name,
		"desc":    material.Desc,
		"file":    material.File,
		"user_id": material.UserID,
	}
	if err := models.AddMaterial(mat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save material to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material added successfully"})
}

func getMaterial(c *gin.Context) {
	idStr := c.Param("id") // Получаем ID из параметров URL

	// Преобразуем строку в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	material, err := models.GetMaterial(id) // Приводим к uint, если это необходимо
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve material"})
		return
	}

	c.JSON(http.StatusOK, material)
}

// func (a *Article) GetAll() ([]*models.Article, error) {
// 	var (
// 		articles, cacheArticles []*models.Article
// 	)

// 	cache := cache_service.Article{
// 		TagID: a.TagID,
// 		State: a.State,

// 		PageNum:  a.PageNum,
// 		PageSize: a.PageSize,
// 	}
// 	key := cache.GetArticlesKey()
// 	if gredis.Exists(key) {
// 		data, err := gredis.Get(key)
// 		if err != nil {
// 			logging.Info(err)
// 		} else {
// 			json.Unmarshal(data, &cacheArticles)
// 			return cacheArticles, nil
// 		}
// 	}

// 	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
// 	if err != nil {
// 		return nil, err
// 	}

// 	gredis.Set(key, articles, 3600)
// 	return articles, nil
// }

// func (a *Article) Delete() error {
// 	return models.DeleteArticle(a.ID)
// }
