package utils

import (
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Article struct {
	Number  string
	Title   string
	Content string
}

func extractArticles(filePath string) ([]Article, error) {
	f, r, err := pdf.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var articles []Article
	var currentArticle Article
	numberPattern := regexp.MustCompile(`^(\d \d \d)$`) // 正则表达式匹配形式为"0 0 1"的序号
	inArticle := false

	for pageNum := 1; pageNum <= r.NumPage(); pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		rows, _ := page.GetTextByRow()
		reverse(rows) // 反转行的顺序

		for _, row := range rows {
			line := ""
			for _, word := range row.Content {
				line += word.S + " "
			}
			line = strings.TrimSpace(line)

			if numberPattern.MatchString(line) {
				if currentArticle.Number != "" {
					articles = append(articles, currentArticle)
					currentArticle = Article{}
				}
				currentArticle.Number = line
				inArticle = false // Reset to capture next line as title
			} else if !inArticle {
				currentArticle.Title = line
				inArticle = true
			} else {
				currentArticle.Content += line + "\n"
			}
		}
	}

	if currentArticle.Number != "" {
		articles = append(articles, currentArticle)
	}

	return articles, nil
}

// 反转函数，用于反转行数组
func reverse(rows pdf.Rows) {
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
}
