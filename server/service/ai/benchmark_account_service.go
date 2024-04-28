package ai

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BenchmarkAccountService struct {
}

var BenchmarkAccountServiceApp = new(BenchmarkAccountService)

func (exa *BenchmarkAccountService) CheckBizId(bizId string) bool {
	var count int64
	global.GVA_DB.Model(&ai.BenchmarkAccount{}).Where("account_id = ?", bizId).Count(&count)
	return count > 0
}
func (exa *BenchmarkAccountService) UpdateBenchmarkAccount(e ai.BenchmarkAccount) error {
	err := global.GVA_DB.Save(e).Error
	return err
}

//@function: CreatePortal
//@description: 创建门户
//@param: e model.Portal
//@return: err error

func (exa *BenchmarkAccountService) CreateBenchmarkAccount(e ai.BenchmarkAccount) (err error) {
	// 拿到微信公众号id
	articleLink := e.ArticleLink
	link, _ := url.Parse(articleLink)
	values, _ := url.ParseQuery(link.RawQuery)
	accountIds := values["__biz"]
	if len(accountIds) == 0 {
		return err
	}
	accountId := accountIds[0]
	e.BASEMODEL = BaseModel()
	e.AccountId = accountId
	err = global.GVA_DB.Create(&e).Error
	if err != nil {
		return err
	}

	// 异步抓取微信公众号文章
	go func() {
		temp := e

		wxToken := &ai.WxToken{}
		er := global.GVA_DB.Model(&ai.WxToken{}).Where("1=1").Last(&wxToken).Error
		if er != nil {
			fmt.Println(er)
			return
		}

		// 爬取微信公众号文章
		articleList := exa.SpiderOfficialAccount(wxToken, temp)

		if len(articleList) > 0 {
			err = global.GVA_DB.CreateInBatches(&articleList, 1000).Error
			if err != nil {
				fmt.Println(err)
			}
		}

		// 再获取点赞数，评论数
		exa.updateMoreInfo(wxToken, temp, articleList)

	}()

	return err
}

func (exa *BenchmarkAccountService) updateMoreInfo(wxToken *ai.WxToken, e ai.BenchmarkAccount, list []*ai.Article) {
	//key := e.Key
	//pass_ticket := wxToken.PassTicket

	//wechatApi := fmt.Sprintf("https://mp.weixin.qq.com/mp/getappmsgext?uin=MjI3NDAxNzUwNw%3D%3D&key=%s&pass_ticket=%s", key, pass_ticket)
	wechatApi := ""
	headers := make(map[string]string)
	headers["user-agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x6309092b) XWEB/8555 Flue"
	for _, item := range list {

		time.Sleep(10 * time.Second)
		articleLink := item.Link

		link, _ := url.Parse(articleLink)
		values, _ := url.ParseQuery(link.RawQuery)

		bizId := values["__biz"][0]
		bizId = strings.ReplaceAll(bizId, "=", "%3D")
		mid := values["mid"][0]
		sn := values["sn"][0]
		idx := values["idx"][0]

		reqBody := fmt.Sprintf("__biz=%s&mid=%s&sn=%s&idx=%s&is_only_read=1&appmsg_type=9", bizId, mid, sn, idx)
		gotStatusCode, gotBody, err := utils.PostWithHeaders(wechatApi, reqBody, headers)
		if gotStatusCode != 200 || err != nil {
			fmt.Println(gotStatusCode, err)
			continue
		}

		var wechatArticleResp WechatArticleResp

		// 将 JSON 字符串解析到结构体中
		err = json.Unmarshal(gotBody, &wechatArticleResp)
		if err != nil {
			fmt.Println("解析 JSON 失败:", err)
			continue
		}

		item.ReadNum = wechatArticleResp.Appmsgstat.ReadNum
		item.LikeNum = wechatArticleResp.Appmsgstat.OldLikeNum

		err = global.GVA_DB.Save(&item).Error
		if err != nil {
			fmt.Println(err)
		}

	}
}

func (exa *BenchmarkAccountService) SpiderOfficialAccount(wxToken *ai.WxToken, e ai.BenchmarkAccount) []*ai.Article {
	tarGetNum := e.InitNum
	accountId := e.AccountId
	token := wxToken.Token
	accountId = strings.ReplaceAll(accountId, "=", "%3D")

	pages := tarGetNum / 5

	curPage := 0

	wechatArticleList := make([]*ai.Article, 0)

	for curPage < pages {
		begin := cast.ToString(5 * curPage)

		time.Sleep(60 * time.Second)

		pageUrl := fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/appmsgpublish?sub=list&search_field=null&begin=%s&count=5&query=&fakeid=%s&type=101_1&free_publish_type=1&sub_action=list_ex&token=%s&lang=zh_CN&f=json&ajax=1", begin, accountId, token)

		cookieList := make([]*http.Cookie, 0)
		ck1 := &http.Cookie{Name: "slave_user", Value: wxToken.SlaveUser}
		cookieList = append(cookieList, ck1)
		ck2 := &http.Cookie{Name: "slave_sid", Value: wxToken.SlaveSid}
		cookieList = append(cookieList, ck2)
		ck3 := &http.Cookie{Name: "bizuin", Value: wxToken.BizUin}
		cookieList = append(cookieList, ck3)
		ck4 := &http.Cookie{Name: "data_bizuin", Value: wxToken.BizUin}
		cookieList = append(cookieList, ck4)
		ck5 := &http.Cookie{Name: "data_ticket", Value: wxToken.DataTicket}
		cookieList = append(cookieList, ck5)
		ck6 := &http.Cookie{Name: "slave_bizuin", Value: wxToken.BizUin}
		cookieList = append(cookieList, ck6)
		ck7 := &http.Cookie{Name: "rand_info", Value: wxToken.RandInfo}
		cookieList = append(cookieList, ck7)

		gotBody, err := utils.GetWithCookie(pageUrl, cookieList)
		if err != nil {
			global.GVA_LOG.Error("SpiderOfficialAccount", zap.Error(err))
		}

		subList := getArticle(e.AccountName, gotBody)

		if len(subList) > 0 {
			wechatArticleList = append(wechatArticleList, subList...)
		}

		curPage++
	}

	return wechatArticleList
}

func getArticle(portalName, resp string) []*ai.Article {
	var officialAccountsResp OfficialAccountsResp

	// 将 JSON 字符串解析到结构体中
	err := json.Unmarshal([]byte(resp), &officialAccountsResp)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return []*ai.Article{}
	}

	var publishPageResp PublishPageResp
	err = json.Unmarshal([]byte(officialAccountsResp.PublishPage), &publishPageResp)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return []*ai.Article{}
	}

	result := make([]*ai.Article, 0)
	for _, item := range publishPageResp.PublishList {
		var publishInfo PublishInfo

		err = json.Unmarshal([]byte(item.PublishInfo), &publishInfo)
		if err != nil {
			fmt.Println("解析 JSON 失败:", err)
			return []*ai.Article{}
		}

		for _, appMsg := range publishInfo.AppMsgEx {
			tags := ""
			for _, tag := range appMsg.AppMsgAlbumInfos {
				tags += tag.Title + ","
			}

			result = append(result, &ai.Article{
				BASEMODEL:   global.BASEMODEL{ID: cast.ToString(utils.GenID()), CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:       appMsg.Title,
				AuthorName:  appMsg.AuthorName,
				PortalName:  portalName,
				Link:        appMsg.Link,
				PublishTime: cast.ToString(appMsg.UpdateTime),
				Tags:        tags,
			})
		}
	}

	return result
}

//@function: DeleteFileChunk
//@description: 删除对标账号
//@param: e model.BenchmarkAccount
//@return: err error

func (exa *BenchmarkAccountService) DeleteBenchmarkAccount(e ai.BenchmarkAccount) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

//@function: GetBenchmarkAccount
//@description: 获取对标账号信息
//@param: id uint
//@return: customer model.BenchmarkAccount, err error

func (exa *BenchmarkAccountService) GetBenchmarkAccount(id uint64) (portal ai.BenchmarkAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&portal).Error
	return
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: GetBenchmarkAccountList
// @description: 分页获取对标账号列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *BenchmarkAccountService) GetBenchmarkAccountList(sysUserAuthorityID uint, info aiReq.BenchmarkAccountSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	var a system.SysAuthority
	a.AuthorityId = sysUserAuthorityID
	auth, err := systemService.AuthorityServiceApp.GetAuthorityInfo(a)
	if err != nil {
		return
	}
	var dataId []uint
	for _, v := range auth.DataAuthorityId {
		dataId = append(dataId, v.AuthorityId)
	}
	var benchmarkAccountList []ai.BenchmarkAccount

	db := global.GVA_DB.Model(&ai.BenchmarkAccount{})

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.AccountName != "" {
		db = db.Where("account_name LIKE ?", "%"+info.AccountName+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return benchmarkAccountList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&benchmarkAccountList).Error
	}
	return benchmarkAccountList, total, err
}

type ContentInfo struct {
	Content string `json:"content"`
}

type Article struct {
	Link             string `json:"link"`
	Title            string `json:"title"`
	Author           string `json:"author"`
	OfficialAccounts string `json:"officialAccounts"`
	PublishTime      string `json:"publishTime"`
	Content          string `json:"content"`
	Tags             string `json:"tags"`
}

type OfficialAccountsResp struct {
	BaseResp struct {
		ErrMsg string `json:"err_msg"`
		Ret    int    `json:"ret"`
	} `json:"base_resp"`
	PublishPage string `json:"publish_page"`
}

type PublishPageResp struct {
	TotalCount    int `json:"total_count"`
	PublishCount  int `json:"publish_count"`
	MasssendCount int `json:"masssend_count"`
	PublishList   []struct {
		PublishType int    `json:"publish_type"`
		PublishInfo string `json:"publish_info"`
	} `json:"publish_list"`
}

type PublishPage struct {
	TotalCount    int `json:"total_count"`
	PublishCount  int `json:"publish_count"`
	MasssendCount int `json:"masssend_count"`
	PublishList   []struct {
		PublishType int         `json:"publish_type"`
		PublishInfo PublishInfo `json:"publish_info"`
	} `json:"publish_list"`
}

type PublishInfo struct {
	Type       int          `json:"type"`
	MsgID      int          `json:"msgid"`
	SentInfo   SentInfo     `json:"sent_info"`
	SentStatus SentStatus   `json:"sent_status"`
	SentResult SentResult   `json:"sent_result"`
	AppMsgInfo []AppMsgInfo `json:"appmsg_info"`
	AppMsgEx   []AppMsgEx   `json:"appmsgex"`
}

type SentInfo struct {
	Time        int  `json:"time"`
	FuncFlag    int  `json:"func_flag"`
	IsSendAll   bool `json:"is_send_all"`
	IsPublished int  `json:"is_published"`
}

type SentStatus struct {
	Total       int `json:"total"`
	Succ        int `json:"succ"`
	Fail        int `json:"fail"`
	Progress    int `json:"progress"`
	UserProtect int `json:"userprotect"`
}

type SentResult struct {
	MsgStatus       int    `json:"msg_status"`
	RefuseReason    string `json:"refuse_reason"`
	RejectIndexList []int  `json:"reject_index_list"`
	UpdateTime      int    `json:"update_time"`
}

type AppMsgInfo struct {
	ShareType      int   `json:"share_type"`
	AppMsgID       int   `json:"appmsgid"`
	VoteID         []int `json:"vote_id"`
	SuperVoteID    []int `json:"super_vote_id"`
	SmartProduct   int   `json:"smart_product"`
	AppMsgLikeType int   `json:"appmsg_like_type"`
	ItemIdx        int   `json:"itemidx"`
	IsPaySubscribe int   `json:"is_pay_subscribe"`
	IsFromTransfer int   `json:"is_from_transfer"`
	OpenFansMsg    int   `json:"open_fansmsg"`
}

type AppMsgEx struct {
	Aid                   string `json:"aid"`
	Title                 string `json:"title"`
	Cover                 string `json:"cover"`
	Link                  string `json:"link"`
	Digest                string `json:"digest"`
	UpdateTime            int    `json:"update_time"`
	AppMsgID              int    `json:"appmsgid"`
	ItemIdx               int    `json:"itemidx"`
	ItemShowType          int    `json:"item_show_type"`
	AuthorName            string `json:"author_name"`
	CreateTime            int    `json:"create_time"`
	HasRedPacketCover     int    `json:"has_red_packet_cover"`
	AlbumID               string `json:"album_id"`
	Checking              int    `json:"checking"`
	MediaDuration         string `json:"media_duration"`
	MediaAPIPublishStatus int    `json:"mediaapi_publish_status"`
	CopyrightType         int    `json:"copyright_type"`
	AppMsgAlbumInfos      []struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		AlbumID   int    `json:"album_id"`
		TagSource int    `json:"tagSource"`
	} `json:"appmsg_album_infos"`
}

type WechatArticleResp struct {
	AdvertisementInfo  []interface{} `json:"advertisement_info"`
	Appid              string        `json:"appid"`
	AppmsgAlbumExtinfo struct {
		NextArticleLink  string        `json:"next_article_link"`
		NextArticleTitle string        `json:"next_article_title"`
		PayHeadImgs      []interface{} `json:"pay_head_imgs"`
		PreArticleLink   string        `json:"pre_article_link"`
		PreArticleTitle  string        `json:"pre_article_title"`
		PreviewTitles    []interface{} `json:"preview_titles"`
	} `json:"appmsg_album_extinfo"`
	AppmsgAlbumVideos []interface{} `json:"appmsg_album_videos"`
	Appmsgact         struct {
		FavoriteBefore int `json:"favorite_before"`
		FollowBefore   int `json:"follow_before"`
		OldLikedBefore int `json:"old_liked_before"`
		PayBefore      int `json:"pay_before"`
		RewardBefore   int `json:"reward_before"`
		SeenBefore     int `json:"seen_before"`
		ShareBefore    int `json:"share_before"`
	} `json:"appmsgact"`
	Appmsgstat struct {
		FriendLikeNum   int  `json:"friend_like_num"`
		IsLogin         bool `json:"is_login"`
		LikeDisabled    bool `json:"like_disabled"`
		LikeNum         int  `json:"like_num"`
		Liked           bool `json:"liked"`
		OldLikeNum      int  `json:"old_like_num"`
		OldLiked        bool `json:"old_liked"`
		OldLikedBefore  int  `json:"old_liked_before"`
		Prompted        int  `json:"prompted"`
		ReadNum         int  `json:"read_num"`
		RealReadNum     int  `json:"real_read_num"`
		Ret             int  `json:"ret"`
		ShareNum        int  `json:"share_num"`
		Show            bool `json:"show"`
		ShowGray        int  `json:"show_gray"`
		ShowLike        int  `json:"show_like"`
		ShowLikeGray    int  `json:"show_like_gray"`
		ShowOldLike     int  `json:"show_old_like"`
		ShowOldLikeGray int  `json:"show_old_like_gray"`
		ShowRead        int  `json:"show_read"`
		Style           int  `json:"style"`
		Version         int  `json:"version"`
		VideoPv         int  `json:"video_pv"`
		VideoUv         int  `json:"video_uv"`
	} `json:"appmsgstat"`
	BaseResp struct {
		ExportkeyToken string `json:"exportkey_token"`
		Ret            int    `json:"ret"`
	} `json:"base_resp"`
	BizfileRet                      int `json:"bizfile_ret"`
	CloseRelatedArticle             int `json:"close_related_article"`
	DistanceToGetRelatedArticleData int `json:"distance_to_get_related_article_data"`
	FavoriteFlag                    struct {
		Show     int `json:"show"`
		ShowGray int `json:"show_gray"`
	} `json:"favorite_flag"`
	FriendSubscribeCount int           `json:"friend_subscribe_count"`
	HitBizrecommend      int           `json:"hit_bizrecommend"`
	IsFans               int           `json:"is_fans"`
	LinkComponentList    []interface{} `json:"link_component_list"`
	MoreReadList         []interface{} `json:"more_read_list"`
	OriginalArticleCount int           `json:"original_article_count"`
	PublicTagInfo        struct {
		Tags []interface{} `json:"tags"`
	} `json:"public_tag_info"`
	RelatedArticleFastClose int           `json:"related_article_fast_close"`
	RelatedArticleUnderAd   int           `json:"related_article_under_ad"`
	RelatedTagVideo         []interface{} `json:"related_tag_video"`
	RewardHeadImgInfos      []interface{} `json:"reward_head_img_infos"`
	RewardHeadImgs          []interface{} `json:"reward_head_imgs"`
	SecControl              struct {
		AdViolationMiddlePage int `json:"ad_violation_middle_page"`
	} `json:"sec_control"`
	ShareFlag struct {
		Show     int `json:"show"`
		ShowGray int `json:"show_gray"`
	} `json:"share_flag"`
	ShowBizBanner      int    `json:"show_biz_banner"`
	ShowRelatedArticle int    `json:"show_related_article"`
	TestFlag           int    `json:"test_flag"`
	VideoContinueFlag  int    `json:"video_continue_flag"`
	WapExportToken     string `json:"wap_export_token"`
}
