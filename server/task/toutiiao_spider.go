package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	ai2 "github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timeutil"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"time"
)

func collectToutiaoArticle() []ai.Article {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	collector.SetRequestTimeout(time.Second * 60)

	subCollector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	// 请求发起时回调,一般用来设置请求头等
	collector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", global.GVA_CONFIG.QianFan.Cookie)
		fmt.Println("----> 开始请求了")
	})

	// 请求发起时回调,一般用来设置请求头等
	subCollector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", global.GVA_CONFIG.QianFan.Cookie)
		fmt.Println(request.URL.Path + "----> 开始请求了")
	})

	// 请求完成后回调
	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
	})

	// 请求完成后回调
	subCollector.OnResponse(func(response *colly.Response) {
		global.GVA_LOG.Info(response.Request.URL.Path+"----> 开始返回了",
			zap.Int("status_code", response.StatusCode),
			zap.Int("content_length", len(response.Body)))
	})

	result := make([]ai.Article, 0)

	collectNum := 0

	// 定义一个回调函数，处理页面响应
	collector.OnHTML("h3", func(e *colly.HTMLElement) {
		articleUrl := e.ChildAttr("a", "href")

		if strings.Contains(articleUrl, "baijiahao") && collectNum <= 3 {
			// 解析URL
			parsedURL, err := url.Parse(articleUrl)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return
			}

			// 更改URL的协议为http
			parsedURL.Scheme = "http"

			// 解析查询参数
			queryParams := parsedURL.Query()

			// 删除特定的查询参数wfr
			queryParams.Del("wfr")

			// 更新URL的查询参数
			parsedURL.RawQuery = queryParams.Encode()

			// 访问文章URL
			subCollector.Visit(parsedURL.String())

			time.Sleep(3 * time.Second)

			collectNum++
		}

	})

	//请求发生错误回调
	collector.OnError(func(response *colly.Response, err error) {
		fmt.Printf("发生错误了:%v", err)
	})

	// 提取标题
	subCollector.OnHTML("div.EaCvy", func(element *colly.HTMLElement) {
		link := element.Request.URL.String()

		title := element.ChildText(".sKHSJ")
		author := element.ChildText("._2gGWi")
		publishTime := element.ChildText("._2sjh9")
		content := element.ChildText("._18p7x")
		//
		title = strings.TrimSpace(title) // 移除多余的空格
		author = strings.TrimSpace(author)
		publishTime = strings.TrimSpace(publishTime)
		content = strings.TrimSpace(content)

		topic := "时事"

		readNum := 0

		item := ai.Article{
			Title:       title,
			Link:        link,
			Content:     content,
			AuthorName:  author,
			PublishTime: publishTime,
			Topic:       topic,
			PortalName:  "百家号",
			ReadNum:     readNum,
		}

		publishTimeInt, _ := timeutil.StrToTimeStamp(publishTime, "2006-01-02 15:04:05")
		// 发布时间需大于今年
		if publishTimeInt < timeutil.GetYearStartTime(int64(time.Now().Year())) {
			return
		}

		item.BASEMODEL = ai2.BaseModel()

		// 将文章添加到结果切片中
		result = append(result, item)
	})

	//encodedParam := url.QueryEscape(hotspot.Headline)

	err := collector.Visit("https://www.toutiao.com/article/7373941840820568610")
	if err != nil {
		global.GVA_LOG.Error("collectArticle", zap.Error(err))
	}

	return result
}

type ToutiaoArticleList struct {
	HasMore bool   `json:"has_more"`
	Message string `json:"message"`
	Data    []struct {
		Abstract       string `json:"abstract"`
		AggrType       int    `json:"aggr_type"`
		ArticleSubType int    `json:"article_sub_type"`
		ArticleType    int    `json:"article_type"`
		ArticleUrl     string `json:"article_url"`
		ArticleVersion int    `json:"article_version"`
		BanComment     bool   `json:"ban_comment"`
		BehotTime      int    `json:"behot_time"`
		BuryCount      int    `json:"bury_count"`
		BuryStyleShow  int    `json:"bury_style_show"`
		CellCtrls      struct {
			CellFlag        int `json:"cell_flag"`
			CellHeight      int `json:"cell_height"`
			CellLayoutStyle int `json:"cell_layout_style"`
		} `json:"cell_ctrls"`
		CellFlag          int    `json:"cell_flag"`
		CellLayoutStyle   int    `json:"cell_layout_style"`
		CellType          int    `json:"cell_type"`
		CommentCount      int    `json:"comment_count"`
		CommonRawData     string `json:"common_raw_data"`
		ContentDecoration string `json:"content_decoration"`
		ControlMeta       struct {
			Modify struct {
				Hide       bool   `json:"hide"`
				Name       string `json:"name"`
				Permission bool   `json:"permission"`
				Tips       string `json:"tips"`
			} `json:"modify"`
			Remove struct {
				Hide       bool   `json:"hide"`
				Name       string `json:"name"`
				Permission bool   `json:"permission"`
				Tips       string `json:"tips"`
			} `json:"remove"`
			Share struct {
				Hide       bool   `json:"hide"`
				Name       string `json:"name"`
				Permission bool   `json:"permission"`
				Tips       string `json:"tips"`
			} `json:"share"`
		} `json:"control_meta"`
		Cursor         int64  `json:"cursor"`
		DataType       int    `json:"data_type"`
		DiggCount      int    `json:"digg_count"`
		DisplayUrl     string `json:"display_url"`
		ForumExtraData string `json:"forum_extra_data"`
		ForwardInfo    struct {
			ForwardCount int `json:"forward_count"`
		} `json:"forward_info"`
		GallaryImageCount int    `json:"gallary_image_count"`
		GroupFlags        int    `json:"group_flags"`
		GroupId           string `json:"group_id"`
		GroupSource       int    `json:"group_source"`
		GroupType         int    `json:"group_type"`
		HasImage          bool   `json:"has_image"`
		HasM3U8Video      bool   `json:"has_m3u8_video"`
		HasMp4Video       bool   `json:"has_mp4_video"`
		HasVideo          bool   `json:"has_video"`
		Hot               int    `json:"hot"`
		Id                string `json:"id"`
		ImageList         []struct {
			Height  int    `json:"height"`
			Uri     string `json:"uri"`
			Url     string `json:"url"`
			UrlList []struct {
				Url string `json:"url"`
			} `json:"url_list"`
			Width int `json:"width"`
		} `json:"image_list"`
		IsOriginal bool `json:"is_original"`
		ItemCell   struct {
			ActionCtrl struct {
				ActionBar struct {
					ActionSettingList []struct {
						ActionType   int `json:"actionType"`
						StyleSetting struct {
							IconKey         string `json:"iconKey"`
							LayoutDirection int    `json:"layoutDirection"`
							Text            string `json:"text"`
						} `json:"styleSetting"`
					} `json:"actionSettingList"`
				} `json:"actionBar"`
				BanBury     bool `json:"banBury"`
				BanComment  bool `json:"banComment"`
				BanDigg     bool `json:"banDigg"`
				ControlMeta struct {
					Modify struct {
						Permission bool   `json:"permission"`
						Tips       string `json:"tips"`
					} `json:"modify"`
					Remove struct {
						Permission bool   `json:"permission"`
						Tips       string `json:"tips"`
					} `json:"remove"`
					Share struct {
						Permission bool   `json:"permission"`
						Tips       string `json:"tips"`
					} `json:"share"`
				} `json:"controlMeta"`
			} `json:"actionCtrl"`
			ArticleBase struct {
				GidStr     string `json:"gidStr"`
				ItemStatus int    `json:"itemStatus"`
			} `json:"articleBase"`
			ArticleClassification struct {
				AggrType           int  `json:"aggrType"`
				ArticleSubType     int  `json:"articleSubType"`
				ArticleType        int  `json:"articleType"`
				BizID              int  `json:"bizID"`
				BizTag             int  `json:"bizTag"`
				GroupSource        int  `json:"groupSource"`
				IsForAudioPlaylist bool `json:"isForAudioPlaylist"`
				IsOriginal         bool `json:"isOriginal"`
				IsSubject          bool `json:"isSubject"`
				Level              int  `json:"level"`
			} `json:"articleClassification"`
			CellCtrl struct {
				BuryStyleShow   int    `json:"buryStyleShow"`
				CellFlag        int    `json:"cellFlag"`
				CellLayoutStyle int    `json:"cellLayoutStyle"`
				CellType        int    `json:"cellType"`
				CellUIType      string `json:"cellUIType"`
				GroupFlags      int    `json:"groupFlags"`
			} `json:"cellCtrl"`
			Extra struct {
				Ping string `json:"ping"`
			} `json:"extra"`
			ImageList struct {
			} `json:"imageList"`
			ItemCounter struct {
				CommentCount    int `json:"commentCount"`
				DiggCount       int `json:"diggCount"`
				ForwardCount    int `json:"forwardCount"`
				ReadCount       int `json:"readCount"`
				RepinCount      int `json:"repinCount"`
				ShareCount      int `json:"shareCount"`
				ShowCount       int `json:"showCount"`
				TextCount       int `json:"textCount"`
				VideoWatchCount int `json:"videoWatchCount"`
				BuryCount       int `json:"buryCount,omitempty"`
			} `json:"itemCounter"`
			LocationInfo struct {
				PublishLocInfo string `json:"publishLocInfo"`
			} `json:"locationInfo"`
			ShareInfo struct {
				ShareControl struct {
					IsHighQuality bool `json:"isHighQuality"`
				} `json:"shareControl,omitempty"`
				ShareURL string `json:"shareURL"`
			} `json:"shareInfo"`
			TagInfo struct {
			} `json:"tagInfo"`
			VideoInfo struct {
			} `json:"videoInfo"`
		} `json:"itemCell"`
		ItemCellDebug  interface{} `json:"itemCellDebug"`
		ItemId         string      `json:"item_id"`
		ItemIdStr      string      `json:"item_id_str"`
		ItemVersion    int         `json:"item_version"`
		LabelStyle     int         `json:"label_style"`
		LargeImageList []struct {
			Height  int    `json:"height"`
			Uri     string `json:"uri"`
			Url     string `json:"url"`
			UrlList []struct {
				Url string `json:"url"`
			} `json:"url_list"`
			Width int `json:"width"`
		} `json:"large_image_list"`
		Level     int `json:"level"`
		LikeCount int `json:"like_count"`
		LogPb     struct {
			CellLayoutStyle string `json:"cell_layout_style"`
			GroupIdStr      string `json:"group_id_str"`
			GroupSource     string `json:"group_source"`
			ImprId          string `json:"impr_id"`
			IsFollowing     string `json:"is_following"`
			IsYaowen        string `json:"is_yaowen"`
		} `json:"log_pb"`
		LynxServer struct {
		} `json:"lynx_server"`
		MiddleImage struct {
			Height  int    `json:"height"`
			Uri     string `json:"uri"`
			Url     string `json:"url"`
			UrlList []struct {
				Url string `json:"url"`
			} `json:"url_list"`
			Width int `json:"width"`
		} `json:"middle_image"`
		NatantLevel int    `json:"natant_level"`
		PreloadWeb  int    `json:"preload_web"`
		PublishTime int    `json:"publish_time"`
		ReadCount   int    `json:"read_count"`
		RebackFlag  int    `json:"reback_flag"`
		RepinCount  int    `json:"repin_count"`
		RepinTime   int    `json:"repin_time"`
		ReqId       string `json:"req_id"`
		ShareUrl    string `json:"share_url"`
		ShowMore    struct {
			Title string `json:"title"`
			Url   string `json:"url"`
		} `json:"show_more"`
		Source         string  `json:"source"`
		SubjectGroupId int     `json:"subject_group_id"`
		Tag            string  `json:"tag,omitempty"`
		TagId          float64 `json:"tag_id"`
		Tip            int     `json:"tip"`
		Title          string  `json:"title"`
		Url            string  `json:"url"`
		UserBury       int     `json:"user_bury"`
		UserDigg       int     `json:"user_digg"`
		UserInfo       struct {
			AvatarUrl       string `json:"avatar_url"`
			Description     string `json:"description"`
			Follow          bool   `json:"follow"`
			Name            string `json:"name"`
			UserAuthInfo    string `json:"user_auth_info"`
			UserId          string `json:"user_id"`
			UserVerified    bool   `json:"user_verified"`
			VerifiedContent string `json:"verified_content"`
		} `json:"user_info"`
		UserLike        int `json:"user_like"`
		UserRepin       int `json:"user_repin"`
		UserRepinTime   int `json:"user_repin_time"`
		VideoDetailInfo struct {
			DetailVideoLargeImage struct {
				Height  int    `json:"height"`
				Uri     string `json:"uri"`
				Url     string `json:"url"`
				UrlList []struct {
					Url string `json:"url"`
				} `json:"url_list"`
				Width int `json:"width"`
			} `json:"detail_video_large_image"`
			DirectPlay       int    `json:"direct_play"`
			GroupFlags       int    `json:"group_flags"`
			ShowPgcSubscribe int    `json:"show_pgc_subscribe"`
			VideoId          string `json:"video_id"`
		} `json:"video_detail_info"`
		VideoDuration int    `json:"video_duration"`
		VideoId       string `json:"video_id,omitempty"`
		VideoStyle    int    `json:"video_style"`
	} `json:"data"`
	Next struct {
		MaxBehotTime int `json:"max_behot_time"`
	} `json:"next"`
	Offset int `json:"offset"`
}
