package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
	"time"
)

func ReadPdf(db *gorm.DB) {
	articleMap, _ := utils.ReadPDFFiles("/Users/changqi.huang/Documents/hao/chengzhang/")

	index := 0

	for title, content := range articleMap {
		if len(content) == 0 || len(content) > 20000 {

			fmt.Println(len(content))
			continue
		}

		title = strings.ReplaceAll(title, ".pdf", "")

		article := &ai.Article{
			Topic:       "个人成长",
			PublishTime: "2024-05-19 02:57:16",
			BASEMODEL:   global.BASEMODEL{ID: cast.ToString(utils.GenID()), CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Title:       title,
			Content:     content,
			Tags:        "pdf",
		}

		err := db.Model(&ai.Article{}).Create(article).Error

		fmt.Println(index, err)
		index++
	}
}
