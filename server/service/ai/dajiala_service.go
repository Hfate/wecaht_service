package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strings"
	"time"
)

type DajialaService struct {
}

var filterWord = []string{"日报", "人民", "官网", "官方", "公众号", "微信", "国际",
	"研究所", "招聘", "本地宝", "活动", "公益", "媒体", "火锅", "咨询", "论坛", "行业",
	"中国", "国家", "论坛", "行业", "公司", "中心", "新闻", "新华", "虎嗅", "腾讯", "央视",
	"网站", "集团", "团队", "深圳", "北京", "上海", "广州", "杭州", "成都", "复旦", "长安街",
	"武汉", "长沙", "重庆", "西安", "南京", "厦门", "泉州", "广东", "湖南", "柳叶刀", "发布", "融媒",
	"电视"}

var DajialaServiceApp = new(DajialaService)

func (receiver *DajialaService) SpiderReadNum(articleList []ai.Article) {
	for _, item := range articleList {
		readNum, err := receiver.ArticleReadNum(item.Link)
		if err != nil {
			global.GVA_LOG.Error("获取阅读数失败", zap.Error(err))
			continue
		}

		item.ReadNum = readNum.Data.Read
		item.LikeNum = readNum.Data.Zan

		// 更新文章阅读量
		err = ArticleServiceApp.Update(item)
		if err != nil {
			global.GVA_LOG.Error("更新文章阅读量失败", zap.Error(err))
			continue
		}

		global.GVA_LOG.Info("热点文章阅读量获取完成",
			zap.String("title", item.Title), zap.String("link", item.Link))

		time.Sleep(2 * time.Second)
	}

	global.GVA_LOG.Info("热点文章阅读量获取完成")

}

func (receiver *DajialaService) SpiderWechatHotArticle() {
	// 获取当前公众号的topic列表
	topicList := OfficialAccountServiceApp.FindTopicList()

	// 获取topic models
	topicModels, err := TopicServiceApp.FindByTopics(topicList)
	if err != nil {
		global.GVA_LOG.Error("获取主题失败", zap.Error(err))
		return
	}

	articleSize := 0
	topicSize := len(topicModels)

	articleList := make([]ai.Article, 0)

	// 获取每个公众号的榜单
	for _, topic := range topicModels {

		if topic.IndustryId == 0 {
			continue
		}

		wpNameList, err := receiver.GetRankList(topic.IndustryId)
		if err != nil {
			global.GVA_LOG.Error("获取榜单失败", zap.Error(err))
			return
		}

		time.Sleep(5 * time.Second)
		for _, wpName := range wpNameList {
			hotLinks, err := receiver.PostHistory(wpName)
			if err != nil {
				global.GVA_LOG.Error("获取历史数据失败", zap.Error(err))
				return
			}

			time.Sleep(5 * time.Second)

			// 获取文章详情
			for _, link := range hotLinks {
				articleDetail, err := receiver.ArticleDetail(link)
				if err != nil {
					global.GVA_LOG.Error("获取文章详情失败", zap.Error(err))
					continue
				}

				url := strings.ReplaceAll(articleDetail.Url, "amp;", "")

				article := ai.Article{
					AuthorName:  articleDetail.Author,
					Link:        url,
					Title:       articleDetail.Title,
					Content:     articleDetail.Content,
					PublishTime: articleDetail.Pubtime,
					PortalName:  articleDetail.Name,
					Topic:       topic.Topic,
					IsHot:       true,
				}

				article.BASEMODEL = global.BASEMODEL{ID: cast.ToString(utils2.GenID()), CreatedAt: time.Now(), UpdatedAt: time.Now()}

				err = ArticleServiceApp.Create(article)
				if err != nil {
					global.GVA_LOG.Error("保存文章失败", zap.Error(err))
					continue
				}

				checkBizId := BenchmarkAccountServiceApp.CheckBizId(articleDetail.Biz)
				if !checkBizId {
					// 创建对标账号
					wechatBenchMarkAccount := ai.BenchmarkAccount{
						AccountName: articleDetail.Name,
						AccountId:   articleDetail.Biz,
						Topic:       topic.Topic,
						InitNum:     10,
						ArticleLink: articleDetail.Url,
					}
					err = BenchmarkAccountServiceApp.CreateBenchmarkAccount(wechatBenchMarkAccount)
					if err != nil {
						global.GVA_LOG.Error("保存对标账号失败", zap.Error(err))
					}
				}

				articleSize++

				articleList = append(articleList, article)

				time.Sleep(5 * time.Second)
			}
		}
	}

	global.GVA_LOG.Info("热点文章获取完成",
		zap.Int("文章数量", articleSize), zap.Int("公众号数量", topicSize))

	receiver.SpiderReadNum(articleList)

}

// GetRankList  获取当日榜单
func (receiver *DajialaService) GetRankList(industryId int) ([]string, error) {
	params := make(map[string]string)
	params["key"] = global.GVA_CONFIG.Dajiala.Key
	params["industry_id"] = cast.ToString(industryId)
	statusCode, resp, err := utils2.Get(global.GVA_CONFIG.Dajiala.GetRankUrl, params)
	if statusCode != 200 {
		global.GVA_LOG.Error("获取榜单失败", zap.Error(err))
		return nil, err
	}

	var rankRespList RankRespList
	err = utils2.JsonStrToStruct(string(resp), &rankRespList)
	if err != nil {
		global.GVA_LOG.Error("获取榜单失败", zap.Error(err))
		return nil, err
	}

	result := make([]string, 0)

	size := 0
	// 找到5个
	resultSet := make(map[string]struct{})
	for _, rank := range rankRespList.Data.Data {
		if !receiver.filterMpName(rank.MpName) {
			continue
		}

		if _, ok := resultSet[rank.Wxid]; ok {
			continue
		}

		result = append(result, rank.Wxid)
		resultSet[rank.Wxid] = struct{}{}
		if size > 3 {
			break
		}
		size++
	}

	return result, nil
}

func (receiver *DajialaService) filterMpName(mpName string) bool {
	for _, word := range filterWord {
		if strings.Contains(mpName, word) {
			return false
		}
	}
	return true

}

// PostHistory 获取公众号历史发文列表
func (receiver *DajialaService) PostHistory(name string) ([]string, error) {

	params := make(map[string]string)
	params["key"] = global.GVA_CONFIG.Dajiala.Key
	params["name"] = name
	statusCode, resp, err := utils2.Get(global.GVA_CONFIG.Dajiala.PostHistoryUrl, params)
	if statusCode != 200 {
		global.GVA_LOG.Error("获取公众号历史发文列表失败", zap.Error(err))
		return []string{}, err
	}

	var postHistoryResp PostHistoryResp
	err = utils2.JsonStrToStruct(string(resp), &postHistoryResp)
	if err != nil {
		global.GVA_LOG.Error("获取公众号历史发文列表失败", zap.Error(err))
		return []string{}, err
	}

	result := make([]string, 0)

	// 爬取第一篇
	if len(postHistoryResp.Data) == 0 {
		return result, nil
	}

	checkUrl := ArticleServiceApp.CheckUrl(postHistoryResp.Data[0].Url)
	if checkUrl {
		return result, nil
	}

	result = append(result, postHistoryResp.Data[0].Url)

	return result, nil
}

// ArticleDetail 获取文章详情
func (receiver *DajialaService) ArticleDetail(url string) (*WechatArticleDetailResp, error) {
	params := make(map[string]string)
	params["key"] = global.GVA_CONFIG.Dajiala.Key
	params["url"] = url
	statusCode, resp, err := utils2.Get(global.GVA_CONFIG.Dajiala.ArticleDetailUrl, params)

	if statusCode != 200 {
		global.GVA_LOG.Error("获取文章详情失败", zap.Error(err))
		return nil, err
	}

	var articleDetailResp ArticleDetailResp
	err = utils2.JsonStrToStruct(string(resp), &articleDetailResp)
	if err != nil {
		global.GVA_LOG.Error("获取文章详情失败", zap.Error(err))
		return nil, err
	}

	return articleDetailResp.Data, nil
}

func (receiver *DajialaService) ArticleReadNum(url string) (ArticleReadNum, error) {
	params := make(map[string]string)
	params["key"] = global.GVA_CONFIG.Dajiala.Key
	params["url"] = url
	statusCode, resp, err := utils2.Get(global.GVA_CONFIG.Dajiala.ReadAndZanUrl, params)

	if statusCode != 200 {
		global.GVA_LOG.Error("获取文章详情失败", zap.Error(err))
		return ArticleReadNum{}, err
	}

	var articleReadNum ArticleReadNum
	err = utils2.JsonStrToStruct(string(resp), &articleReadNum)
	if err != nil {
		global.GVA_LOG.Error("获取文章阅读数目失败", zap.Error(err))
		return ArticleReadNum{}, err
	}

	return articleReadNum, nil
}

type ArticleReadNum struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Read    int `json:"read"`
		Zan     int `json:"zan"`
		Looking int `json:"looking"`
	} `json:"data"`
}

type WechatArticleDetailResp struct {
	Url            string        `json:"url"`
	Hashid         string        `json:"hashid"`
	Title          string        `json:"title"`
	Content        string        `json:"content"`
	PureText       string        `json:"pure_text"`
	GhId           string        `json:"gh_id"`
	Name           string        `json:"name"`
	Index          string        `json:"index"`
	Pubtime        string        `json:"pubtime"`
	Biz            string        `json:"biz"`
	Desc           string        `json:"desc"`
	Wxid           string        `json:"wxid"`
	ArticleHeadImg string        `json:"article_head_img"`
	MpHeadImg      string        `json:"mp_head_img"`
	CreateTime     string        `json:"create_time"`
	Signature      string        `json:"signature"`
	MsgDailyIdx    int           `json:"msg_daily_idx"`
	Author         string        `json:"author"`
	SourceUrl      string        `json:"source_url"`
	Copyright      int           `json:"copyright"`
	VideoIds       []interface{} `json:"video_ids"`
	VoiceInAppmsg  []interface{} `json:"voice_in_appmsg"`
	Type           string        `json:"type"`
	Servicetype    string        `json:"servicetype"`
	IpWording      struct {
		CountryName  string `json:"country_name"`
		CountryId    string `json:"country_id"`
		ProvinceName string `json:"province_name"`
	} `json:"ip_wording"`
}

type ArticleDetailResp struct {
	Code        int                      `json:"code"`
	Msg         string                   `json:"msg"`
	Data        *WechatArticleDetailResp `json:"data"`
	CostMoney   float64                  `json:"cost_money"`
	RemainMoney float64                  `json:"remain_money"`
}

type PostHistoryResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Position      int    `json:"position"`
		Url           string `json:"url"`
		PostTime      int    `json:"post_time"`
		PostTimeStr   string `json:"post_time_str"`
		CoverUrl      string `json:"cover_url"`
		Original      int    `json:"original"`
		ItemShowType  int    `json:"item_show_type"`
		Digest        string `json:"digest"`
		Title         string `json:"title"`
		PrePostTime   int    `json:"pre_post_time"`
		Appmsgid      int64  `json:"appmsgid"`
		MsgStatus     int    `json:"msg_status"`
		MsgFailReason string `json:"msg_fail_reason"`
		SendToFansNum int    `json:"send_to_fans_num"`
		UpdateTime    int    `json:"update_time"`
		IsDeleted     string `json:"is_deleted"`
		Types         int    `json:"types"`
		PicCdnUrl2351 string `json:"pic_cdn_url_235_1"`
		PicCdnUrl169  string `json:"pic_cdn_url_16_9"`
		PicCdnUrl11   string `json:"pic_cdn_url_1_1"`
	} `json:"data"`
	TotalNum           int     `json:"total_num"`
	TotalPage          int     `json:"total_page"`
	PublishCount       int     `json:"publish_count"`
	MasssendCount      int     `json:"masssend_count"`
	NowPage            int     `json:"now_page"`
	NowPageArticlesNum int     `json:"now_page_articles_num"`
	MpNickname         string  `json:"mp_nickname"`
	MpWxid             string  `json:"mp_wxid"`
	MpGhid             string  `json:"mp_ghid"`
	HeadImg            string  `json:"head_img"`
	CostMoney          float64 `json:"cost_money"`
	RemainMoney        float64 `json:"remain_money"`
}

type RankRespList struct {
	Msg       string `json:"msg"`
	ErrorCode int    `json:"error_code"`
	Data      struct {
		Data []struct {
			Rank          int         `json:"rank"`
			Id            int         `json:"id"`
			MpName        string      `json:"mp_name"`
			Wxid          string      `json:"wxid"`
			RankAid       int         `json:"rank_aid"`
			AvgTopReadnum interface{} `json:"avg_top_readnum"`
			AvgReadnum    interface{} `json:"avg_readnum"`
			TotalRead     interface{} `json:"total_read"`
			PostTotal     int         `json:"post_total"`
			ZanTotal      int         `json:"zan_total"`
			RealIndex     float64     `json:"real_index"`
			OriginalIndex float64     `json:"original_index"`
			DajialaIndex  float64     `json:"dajiala_index"`
			Avatar        string      `json:"avatar"`
			IsFavorite    int         `json:"is_favorite"`
		} `json:"data"`
		Total int `json:"total"`
	} `json:"data"`
}

type RankResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *RankRespList
}

type ArticleDetailReq struct {
	Url string
}
