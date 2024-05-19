package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/task"

	"github.com/robfig/cron/v3"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func Timer() {
	//task.SpiderWechatHotArticle()

	//task.ReadPdf(global.GVA_DB)

	// spec 定时任务详细配置参考 https://pkg.go.dev/github.com/robfig/cron?utm_source=godoc
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		//_, err = global.GVA_Timer.AddTaskByFunc("定时任务标识", "@hourly", func() {
		//	err = task.PortalSpider(global.GVA_DB)
		//	if err != nil {
		//		fmt.Println("timer error:", err)
		//	}
		//}, "定时爬取门户数据", option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}

		_, err = global.GVA_Timer.AddTaskByFunc("定时任务标识", "@every 5m", func() {
			err = task.HotspotSpider(global.GVA_DB)
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时收集热点", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("定时重置对标账号标记", "@every 5h", func() {
			task.ResetSpider(global.GVA_DB)
		}, "定时重置对标账号标记", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("爬取百家号", "@every 5m", func() {
			task.SpiderHotArticle()
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "爬取百家号", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("定时爬取微信公众号", "@hourly", func() {
			err = task.WechatSpider(global.GVA_DB)
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时爬取微信公众号", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 每天凌晨四点执行生成文章的任务
		//_, err = global.GVA_Timer.AddTaskByFunc("GenerateArticle", "@daily", func() {
		//	task.GenerateArticle()
		//	if err != nil {
		//		fmt.Println("timer error:", err)
		//	}
		//}, "定时生成文章", option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}

		// 每天凌晨执行生成文章的任务
		_, err = global.GVA_Timer.AddTaskByFunc("SpiderWechatHotArticle", "@daily", func() {
			task.SpiderWechatHotArticle()
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时爬取爆文", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

	}()

}
