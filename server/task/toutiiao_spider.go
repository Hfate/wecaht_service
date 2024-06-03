package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	ai2 "github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timeutil"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strings"
	"time"
)

func CollectToutiaoArticle() {
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	subCollector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	collector.SetRequestTimeout(time.Second * 60)
	subCollector.SetRequestTimeout(time.Second * 60)

	// 请求发起时回调,一般用来设置请求头等
	collector.OnRequest(func(request *colly.Request) {
	})

	resultMap := make(map[string]*ai.Article)
	// 请求完成后回调
	collector.OnResponse(func(response *colly.Response) {
		respStr := string(response.Body)
		feedResp := &FeedResp{}
		utils.JsonStrToStruct(respStr, feedResp)

		if len(feedResp.Data) > 0 {
			for _, item := range feedResp.Data {
				feedConet := &FeedContent{}
				utils.JsonStrToStruct(item.Content, feedConet)
				if feedConet.ReadCount > 100000 {

					topic := ""
					if feedConet.UserInfo.UserAuthInfo != "" {
						userAuthInfo := &UserAuthInfo{}
						utils.JsonStrToStruct(feedConet.UserInfo.UserAuthInfo, userAuthInfo)
						authInfo := userAuthInfo.AuthInfo
						authInfo = strings.ReplaceAll(authInfo, "领域创作者", "")
						authInfo = strings.ReplaceAll(authInfo, "优质", "")
						topic = authInfo
					}

					publishTime := int64(feedConet.PublishTime) * 1000
					curTime := timeutil.GetCurTime()
					sevenDayBefore := timeutil.AddDays(curTime, -7)
					if publishTime < sevenDayBefore {
						continue
					}

					resultMap[feedConet.ArticleUrl] = &ai.Article{
						PortalName:  "今日头条",
						PublishTime: timeutil.Format(publishTime, timeutil.DateTimeLong),
						ReadNum:     feedConet.ReadCount,
						LikeNum:     feedConet.LikeCount,
						Link:        feedConet.ArticleUrl,
						Topic:       topic,
						Title:       feedConet.Title,
					}

					subCollector.Visit(feedConet.ArticleUrl)
				}
			}
		}

	})

	// 请求发起时回调,一般用来设置请求头等
	subCollector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Cookie", "__ac_nonce=0665c7753007e1cb7a608; _ga_QEHZPBE5HH=GS1.1.1717331872.5.1.1717335891.0.0.0; __ac_signature=_02B4Z6wo00f01RGuJyQAAIDDrCEMzB4yQDkRiiOAACI6ygmk--UIXy5rYsqsj77ZsKNZ1BUBXNjXmFRHWTmse5fovWgCM2ZCcr7uLGHU11lH73Ic9X7qRlSym7LE36oEZb12cjvPbtdyGBKFd1; __ac_referer=__ac_blank; msToken=CmnfAvSofZoACxNZrC_KTPsiwKryd4S0Yq-ProVEXgDhKVI4Bc31IuLl34lxYGu_VO_XkrNAYQtKNmaMzfrVGvDQc2kYbk42cobWHsMN; ttwid=1%7CjJhV-ZQBbp-rRPiTpQX02ojxAnJzbXA3xow5e5Q7iTA%7C1717335892%7Cb3e97916cc5efc25a7706234d6b4810ae649450704a42f890f87091b690f756d")
	})

	// 请求完成后回调
	subCollector.OnResponse(func(response *colly.Response) {

	})

	collectSize := 0
	subCollector.OnHTML("div.article-content", func(element *colly.HTMLElement) {

		url := element.Request.URL.String()

		author := element.ChildText(".name")
		content := element.ChildText(".syl-article-base")
		//
		author = strings.TrimSpace(author)
		content = strings.TrimSpace(content)

		groupUrl := strings.ReplaceAll(url, "article", "group")
		groupUrl = strings.ReplaceAll(groupUrl, "www.", "")
		item := resultMap[groupUrl]

		if item == nil {
			return
		}

		item.Content = content
		item.AuthorName = author
		item.Link = url

		item.BASEMODEL = ai2.BaseModel()

		global.GVA_DB.Create(item)

		curTime := timeutil.GetCurTime()
		startOfDay := timeutil.GetDateStartTime(curTime)
		publishTime, _ := timeutil.StrToTimeStamp(item.PublishTime, timeutil.DateTimeLong)

		if publishTime > startOfDay {
			PublishArticle(item)
		}

		collectSize++

	})

	//请求发生错误回调
	collector.OnError(func(response *colly.Response, err error) {
		fmt.Printf("发生错误了:%v", err)
	})

	//encodedParam := url.QueryEscape(hotspot.Headline)

	err := collector.Visit("https://lf.snssdk.com/api/news/feed/v200/")
	if err != nil {
		global.GVA_LOG.Error("collectArticle", zap.Error(err))
	}

	fmt.Println("头条爬取完成【" + cast.ToString(collectSize) + "】")
}

func PublishArticle(article *ai.Article) {
	topic := article.Topic

	accountList := make([]*ai.OfficialAccount, 0)
	global.GVA_DB.Model(&ai.OfficialAccount{}).Where("topic=?", topic).Where("is_publish=0").Limit(1).Find(&accountList)

	if len(accountList) == 0 {
		topic = "热点"
		global.GVA_DB.Model(&ai.OfficialAccount{}).Where("topic=?", topic).Where("is_publish=0").Limit(1).Find(&accountList)
	}

	if len(accountList) > 0 {
		account := accountList[0]
		batchId := timeutil.GetCurDate() + account.AppId

		global.GVA_DB.Where("batch_id=?", batchId).Delete(&ai.DailyArticle{})

		aiDailyArticle := ai.DailyArticle{
			Title:             article.Title,
			PortalName:        article.PortalName,
			Topic:             article.Topic,
			AuthorName:        article.AuthorName,
			Tags:              article.Tags,
			Content:           article.Content,
			BatchId:           batchId,
			ReadNum:           article.ReadNum,
			LikeNum:           article.LikeNum,
			CommentNum:        article.CommentNum,
			Link:              article.Link,
			UseTimes:          0,
			HotspotId:         article.HotspotId,
			TargetAccountId:   account.AppId,
			TargetAccountName: account.AccountName,
		}
		aiDailyArticle.BASEMODEL = ai2.BaseModel()

		article.UseTimes = 1
		global.GVA_DB.Save(article)

		global.GVA_DB.Model(&ai.DailyArticle{}).Create(&aiDailyArticle)

		account.IsPublish = 1
		global.GVA_DB.Save(account)

		if topic == "热点" {
			account.Topic = "时事"
		}

		// 改写文章
		ai2.AIArticleServiceApp.GenerateArticle(account)
	}

}

type UserAuthInfo struct {
	Thread struct {
		AuthType string `json:"auth_type"`
		AuthInfo string `json:"auth_info"`
	} `json:"thread"`
	AuthType  string `json:"auth_type"`
	AuthInfo  string `json:"auth_info"`
	OtherAuth struct {
		Interest string `json:"interest"`
	} `json:"other_auth"`
}

type FeedContent struct {
	ReadCount     int         `json:"read_count"`
	BehotTime     int         `json:"behot_time"`
	DetailContent string      `json:"detail_content"`
	Tip           int         `json:"tip"`
	GroupIdStr    string      `json:"group_id_str"`
	SmallImage    interface{} `json:"small_image"`
	ForwardInfo   struct {
		ForwardCount int `json:"forward_count"`
	} `json:"forward_info"`
	ShowPortrait    bool   `json:"show_portrait"`
	TagId           int64  `json:"tag_id"`
	ShareCount      int    `json:"share_count"`
	ItemId          int64  `json:"item_id"`
	Cursor          int64  `json:"cursor"`
	Url             string `json:"url"`
	SourceOpenUrl   string `json:"source_open_url"`
	BanDanmaku      bool   `json:"ban_danmaku"`
	GroupSource     int    `json:"group_source"`
	SourceIconStyle int    `json:"source_icon_style"`
	CommonRawData   struct {
	} `json:"common_raw_data"`
	UserInfo struct {
		AvatarUrl       string `json:"avatar_url"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		UserId          int64  `json:"user_id"`
		UserVerified    bool   `json:"user_verified"`
		VerifiedContent string `json:"verified_content"`
		Follow          bool   `json:"follow"`
		FollowerCount   int    `json:"follower_count"`
		Schema          string `json:"schema"`
		LiveInfoType    int    `json:"live_info_type"`
		UserAuthInfo    string `json:"user_auth_info"`
	} `json:"user_info"`
	AggrType              int  `json:"aggr_type"`
	HasM3U8Video          bool `json:"has_m3u8_video"`
	BanComment            int  `json:"ban_comment"`
	NeedClientImprRecycle int  `json:"need_client_impr_recycle"`
	CommentCount          int  `json:"comment_count"`
	ItemCell              struct {
		ImageList struct {
		} `json:"imageList"`
		ShareInfo struct {
			CoverImage struct {
				Uri       string   `json:"uri"`
				Height    int      `json:"height"`
				ImageType int      `json:"imageType"`
				Url       string   `json:"url"`
				Width     int      `json:"width"`
				UrlList   []string `json:"urlList"`
			} `json:"coverImage"`
			ShareType struct {
				Wx    int `json:"wx"`
				Qq    int `json:"qq"`
				Qzone int `json:"qzone"`
				Pyq   int `json:"pyq"`
			} `json:"shareType"`
			WeixinCoverImage struct {
				Width     int      `json:"width"`
				UrlList   []string `json:"urlList"`
				Uri       string   `json:"uri"`
				Height    int      `json:"height"`
				ImageType int      `json:"imageType"`
				Url       string   `json:"url"`
			} `json:"weixinCoverImage"`
			TokenType int    `json:"tokenType"`
			ShareURL  string `json:"shareURL"`
			Title     string `json:"title"`
		} `json:"shareInfo"`
		CellCtrl struct {
			CellLayoutStyle  int    `json:"cellLayoutStyle"`
			CellType         int    `json:"cellType"`
			GroupFlags       int    `json:"groupFlags"`
			BuryStyleShow    int    `json:"buryStyleShow"`
			TopShortTitle    string `json:"topShortTitle"`
			ArticleImageType int    `json:"articleImageType"`
			CellFlag         int    `json:"cellFlag"`
		} `json:"cellCtrl"`
		VideoInfo struct {
		} `json:"videoInfo"`
		ActionCtrl struct {
			BanBury               bool `json:"banBury"`
			BanDigg               bool `json:"banDigg"`
			BanDanmaku            bool `json:"banDanmaku"`
			PreloadWeb            int  `json:"preloadWeb"`
			NeedClientImprRecycle bool `json:"needClientImprRecycle"`
			IgnoreWebTransform    bool `json:"ignoreWebTransform"`
			ActionList            []struct {
				Action int `json:"action"`
			} `json:"actionList"`
			BanComment  bool `json:"banComment"`
			ShowDislike bool `json:"showDislike"`
			FilterWord  []struct {
				Id         string `json:"id"`
				Name       string `json:"name"`
				IsSelected bool   `json:"isSelected"`
			} `json:"filterWord"`
			ActionBar struct {
				ActionSettingList []struct {
					ActionType   int `json:"actionType"`
					StyleSetting struct {
						Text            string `json:"text"`
						IconKey         string `json:"iconKey"`
						LayoutDirection int    `json:"layoutDirection"`
					} `json:"styleSetting"`
				} `json:"actionSettingList"`
			} `json:"actionBar"`
		} `json:"actionCtrl"`
		ItemCounter struct {
			VideoWatchCount int `json:"videoWatchCount"`
			TextCount       int `json:"textCount"`
			CommentCount    int `json:"commentCount"`
			DiggCount       int `json:"diggCount"`
			ReadCount       int `json:"readCount"`
			ShareCount      int `json:"shareCount"`
			ForwardCount    int `json:"forwardCount"`
			ShowCount       int `json:"showCount"`
			RepinCount      int `json:"repinCount"`
		} `json:"itemCounter"`
		TagInfo struct {
			Label string `json:"label"`
		} `json:"tagInfo"`
		ArticleClassification struct {
			ArticleType        int  `json:"articleType"`
			GroupSource        int  `json:"groupSource"`
			AggrType           int  `json:"aggrType"`
			Level              int  `json:"level"`
			IsSubject          bool `json:"isSubject"`
			BizTag             int  `json:"bizTag"`
			IsForAudioPlaylist bool `json:"isForAudioPlaylist"`
			ArticleSubType     int  `json:"articleSubType"`
			BizID              int  `json:"bizID"`
			IsStick            bool `json:"isStick"`
		} `json:"articleClassification"`
		ArticleBase struct {
			ItemStatus int    `json:"itemStatus"`
			GidStr     string `json:"gidStr"`
		} `json:"articleBase"`
		ThreadCustom struct {
		} `json:"threadCustom"`
		Extra struct {
			Ping                   string `json:"ping"`
			IsThreadWaterfallTuwen string `json:"is_thread_waterfall_tuwen"`
		} `json:"extra"`
	} `json:"itemCell"`
	ItemVersion   int    `json:"item_version"`
	AllowDownload bool   `json:"allow_download"`
	ArticleUrl    string `json:"article_url"`
	FilterWords   []struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		IsSelected bool   `json:"is_selected"`
	} `json:"filter_words"`
	BuryCount   int  `json:"bury_count"`
	ShowDislike bool `json:"show_dislike"`
	IsKeyVideo  bool `json:"is_key_video"`
	ShareInfo   struct {
		ShareUrl   string `json:"share_url"`
		Title      string `json:"title"`
		CoverImage struct {
			Url     string `json:"url"`
			Width   int    `json:"width"`
			UrlList []struct {
				Url string `json:"url"`
			} `json:"url_list"`
			Uri    string `json:"uri"`
			Height int    `json:"height"`
		} `json:"cover_image"`
		ShareType struct {
			Wx    int `json:"wx"`
			Qq    int `json:"qq"`
			Qzone int `json:"qzone"`
			Pyq   int `json:"pyq"`
		} `json:"share_type"`
		WeixinCoverImage struct {
			Url     string `json:"url"`
			Width   int    `json:"width"`
			UrlList []struct {
				Url string `json:"url"`
			} `json:"url_list"`
			Uri    string `json:"uri"`
			Height int    `json:"height"`
		} `json:"weixin_cover_image"`
		TokenType  int `json:"token_type"`
		OnSuppress int `json:"on_suppress"`
	} `json:"share_info"`
	HasMp4Video       int    `json:"has_mp4_video"`
	DiggCount         int    `json:"digg_count"`
	LikeCount         int    `json:"like_count"`
	FeedTitle         string `json:"feed_title"`
	ContentDecoration string `json:"content_decoration"`
	ArticleVersion    int    `json:"article_version"`
	UserVerified      int    `json:"user_verified"`
	MediaName         string `json:"media_name"`
	RepinCount        int    `json:"repin_count"`
	CellType          int    `json:"cell_type"`
	Hot               int    `json:"hot"`
	ActionList        []struct {
		Action int    `json:"action"`
		Desc   string `json:"desc"`
		Extra  struct {
		} `json:"extra"`
	} `json:"action_list"`
	BanImmersive int `json:"ban_immersive"`
	VideoStyle   int `json:"video_style"`
	OptionalData struct {
		BanDownload       string `json:"ban_download"`
		YunyuType         string `json:"yunyu_type"`
		KeynewsExpireTime string `json:"keynews_expire_time"`
	} `json:"optional_data"`
	ArticleType  int `json:"article_type"`
	UgcRecommend struct {
		Activity string `json:"activity"`
		Reason   string `json:"reason"`
	} `json:"ugc_recommend"`
	ShareType          int    `json:"share_type"`
	XiRelated          bool   `json:"xi_related"`
	IgnoreWebTransform int    `json:"ignore_web_transform"`
	PublishTime        int    `json:"publish_time"`
	DisplayUrl         string `json:"display_url"`
	Level              int    `json:"level"`
	ShareUrl           string `json:"share_url"`
	LabelStyle         int    `json:"label_style"`
	ArticleSubType     int    `json:"article_sub_type"`
	Label              string `json:"label"`
	IsStick            bool   `json:"is_stick"`
	VerifiedContent    string `json:"verified_content"`
	GroupId            int64  `json:"group_id"`
	CellFlag           int    `json:"cell_flag"`
	MediaInfo          struct {
		AvatarUrl       string `json:"avatar_url"`
		Name            string `json:"name"`
		UserVerified    bool   `json:"user_verified"`
		MediaId         int64  `json:"media_id"`
		UserId          int64  `json:"user_id"`
		VerifiedContent string `json:"verified_content"`
		IsStarUser      bool   `json:"is_star_user"`
		RecommendReason string `json:"recommend_reason"`
		RecommendType   int    `json:"recommend_type"`
		Follow          bool   `json:"follow"`
	} `json:"media_info"`
	StickStyle      int         `json:"stick_style"`
	Abstract        string      `json:"abstract"`
	Rid             string      `json:"rid"`
	InteractionData string      `json:"interaction_data"`
	CellLayoutStyle int         `json:"cell_layout_style"`
	RawAdData       interface{} `json:"raw_ad_data"`
	UserRepin       int         `json:"user_repin"`
	Source          string      `json:"source"`
	LogPb           struct {
		GroupSource     string `json:"group_source"`
		UiStyle         string `json:"ui_style"`
		LogpbGroupId    string `json:"logpb_group_id"`
		CellLayoutStyle string `json:"cell_layout_style"`
		IsFollowing     string `json:"is_following"`
		IsYaowen        string `json:"is_yaowen"`
		ImprId          string `json:"impr_id"`
		AuthorId        string `json:"author_id"`
		SentinelType    string `json:"sentinel_type"`
		DItemDataSource string `json:"d_item_data_source"`
		IsPortraitShown string `json:"is_portrait_shown"`
	} `json:"log_pb"`
	ContentHash         string `json:"content_hash"`
	Title               string `json:"title"`
	BuryStyleShow       int    `json:"bury_style_show"`
	ShowPortraitArticle bool   `json:"show_portrait_article"`
	HasVideo            bool   `json:"has_video"`
	IsSubject           bool   `json:"is_subject"`
	ShowMaxLine         int    `json:"show_max_line"`
}

type FeedResp struct {
	Data []struct {
		Content string `json:"content"`
		Code    string `json:"code"`
	} `json:"data"`
	SubEntranceList   []interface{} `json:"sub_entrance_list"`
	Errno             int           `json:"errno"`
	Message           string        `json:"message"`
	TotalNumber       int           `json:"total_number"`
	HasMore           bool          `json:"has_more"`
	LoginStatus       int           `json:"login_status"`
	ShowEtStatus      int           `json:"show_et_status"`
	PostContentHint   string        `json:"post_content_hint"`
	HasMoreToRefresh  bool          `json:"has_more_to_refresh"`
	ActionToLastStick int           `json:"action_to_last_stick"`
	SubEntranceStyle  int           `json:"sub_entrance_style"`
	FeedFlag          int           `json:"feed_flag"`
	Tips              struct {
		Type            string `json:"type"`
		DisplayDuration int    `json:"display_duration"`
		DisplayInfo     string `json:"display_info"`
		DisplayTemplate string `json:"display_template"`
		OpenUrl         string `json:"open_url"`
		WebUrl          string `json:"web_url"`
		DownloadUrl     string `json:"download_url"`
		AppName         string `json:"app_name"`
		PackageName     string `json:"package_name"`
		StreamRespCnt   int    `json:"stream_resp_cnt"`
	} `json:"tips"`
	FollowRecommendTips  string `json:"follow_recommend_tips"`
	HideTopcellCount     int    `json:"hide_topcell_count"`
	IsUseBytedanceStream bool   `json:"is_use_bytedance_stream"`
	GetOfflinePool       bool   `json:"get_offline_pool"`
	ApiBaseInfo          struct {
		InfoType       int    `json:"info_type"`
		AppExtraParams string `json:"app_extra_params"`
	} `json:"api_base_info"`
	ShowLastRead      bool `json:"show_last_read"`
	LastResponseExtra struct {
		Data string `json:"data"`
	} `json:"last_response_extra"`
	Offset int    `json:"offset"`
	Tail   string `json:"tail"`
	Extra  string `json:"extra"`
}
