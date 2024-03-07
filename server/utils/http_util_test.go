package utils

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPostBodyWithHeaders(t *testing.T) {
	url := "https://mp.weixin.qq.com/mp/getappmsgext?uin=MjI3NDAxNzUwNw%253D%253D&key=2d3c4771ae605be54eab0500606b9aca72620362a2796655b97d0631fdd28af2de2baa3dd94eeafbeb3d619072485fb1d225f12cdbda36b8a5b152ae829af17053b985d2dc11cad28b11672a439cd73a407d66a5d9fb5394f7840a369729979bd0a5de21c77da749352ebefd7011be4493d2d812b0438c8625dba0af708be37e&pass_ticket=hMlJaboOsNIicJ5TxSN9l4ICQCk7GMwfbJ20MyRDOiNw3GwIlsxIdU1L%25252FJtWuFa%25252FmHP7tZ1Fr4o30a0lIekT5Q%25253D%25253D&wxtoken=&devicetype=Windows%2B10%2Bx64&clientversion=6309092b"
	url = "https://mp.weixin.qq.com/mp/getappmsgext?f=json&mock=&uin=MjI3NDAxNzUwNw%3D%3D&key=4902a9562fb292e17fd63fedf9d6bc6cf02d49765483a2dee83844e600f783ad3314520c37e3a60cab10c9a50e80c1ff6e348776f02beb554cceefa4f237997df5c7330557c5955fd1bd81f838a08372280870bc3032e3d0cbe3058968abf97fe847f05108e6de744f089e7d4717d2a1c6dfa11994ddff746f2768ae7b8dcd6f&pass_ticket=hMlJaboOsNIicJ5TxSN9l4ICQCk7GMwfbJ20MyRDOiM2FtBe7yGn6G34scVxzMrQ5S%2Fw5MsBngvry0zifnq12g%3D%3D&wxtoken=777&devicetype=Windows%26nbsp%3B10%26nbsp%3Bx64&clientversion=6309092b&__biz=MzA4MzEwMzAwOA%3D%3D&enterid=1709579708&appmsg_token=1259_zNOg8egNdaZ2wu1LgPWuA9ukhXxGwqfM2tREgMqg6KT113YtAW0Aoon7Z4XEqmcRT1XatDb7yXz0LtaN&x5=0&f=json"
	headers := make(map[string]string)
	headers["user-agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x6309092b) XWEB/8555 Flue"
	reqBody := "r=0.12149454670049198&__biz=MzA4MzEwMzAwOA%3D%3D&appmsg_type=9&mid=2649631499&sn=6210c6deea32ddea7ffa728d870ba3f4&idx=1&scene=90&title=%25E9%2594%2580%25E9%2587%258F%25E7%258B%2582%25E6%25B6%25A8%25E8%25B6%2585500%2525%25EF%25BC%258C%25E5%25A4%259A%25E5%259C%25B0%25E5%25B7%25B2%25E6%2596%25AD%25E8%25B4%25A7%25EF%25BC%2581&ct=1709560869&abtest_cookie=&devicetype=Windows%2010%20x64&version=6309092b&is_need_ticket=0&is_need_ad=0&comment_id=3354915842291056650&is_need_reward=0&both_ad=0&reward_uin_count=0&send_time=&msg_daily_idx=1&is_original=0&is_only_read=1&req_id=05031QAmi4cBwKQSkBuLPmPe&pass_ticket=hMlJaboOsNIicJ5TxSN9l4ICQCk7GMwfbJ20MyRDOiM2FtBe7yGn6G34scVxzMrQ5S%2Fw5MsBngvry0zifnq12g%3D%3D&is_temp_url=0&item_show_type=0&tmp_version=1&more_read_type=0&appmsg_like_type=2&related_video_sn=&related_video_num=5&vid=&is_pay_subscribe=0&pay_subscribe_uin_count=0&has_red_packet_cover=0&album_id=1296223588617486300&album_video_num=5&cur_album_id=undefined&is_public_related_video=NaN&encode_info_by_base64=undefined&exptype=&export_key=n_ChQIAhIQK3ieNkHCsmfL6G2%252F0DvqqRLgAQIE97dBBAEAAAAAAEM1OaRPAFYAAAAOpnltbLcz9gKNyK89dVj01Wov%252BTv57vrVzWhDnJ2c5koPOBSjF4WWRoj2E67eWrkgRTYx47vle3%252BsnMvmP%252BRlHcMWiFHYRu%252F%252F2GqWEQqf0%252FuUNibJIG7EKfqgeiNDZb2oPI0uE1ajabkUzwQWrTU3hHIoSdBlGnCGlfNYrJUhQepchZbDaH9jxpD7PTEHpYbj46693LiaJqG7H%252B4oLU7Hjnk8ez0LKtCCtIdGPhdny8cDO2k7sQMHwnkjbu7%252BTfjdxDoZdJ2sWH96&export_key_extinfo=&business_type=0"

	// appmsg_type,is_temp_url,is_only_read
	reqBody = "r=0.02006583521972516&appmsg_type=9&mid=2649631392&sn=23539f1125a1c1f97ad2272d88c02b3d&idx=1&scene=126&title=%25E4%25BB%25B7%25E6%25A0%25BC%25E8%2585%25B0%25E6%2596%25A9%25EF%25BC%2581%25E7%25AA%2581%25E5%258F%2591%25E5%25A4%25A7%25E8%25B7%25B3%25E6%25B0%25B4&ct=1708955096&abtest_cookie=&devicetype=Windows%2B10%2Bx64&version=6309092b&is_need_ad=0&both_ad=0&send_time=&msg_daily_idx=1&pass_ticket=hMlJaboOsNIicJ5TxSN9l4ICQCk7GMwfbJ20MyRDOiMGkXUTPZNkxj%25252FHmttBUgqloFA1WAKNdeDU9UoGxdfdJQ%25253D%25253D&is_temp_url=0&item_show_type=0&tmp_version=1&pos_type_list=&vid_list=%5B%5D&exportkey=n_ChQIAhIQYWeby2rGgKvQT4os8OoAjhLgAQIE97dBBAEAAAAAAHD2Fb96z0sAAAAOpnltbLcz9gKNyK89dVj0PTDyZ2BDxmasm8TnA4L%252BcLkfgQmw1VDgOXNWNmkSM2%252FfFnqxC6aByoOjVXroePsvdIFQXXvLyzhcwsf6SXfRR9amo4L3ru2WfmR5EcepcpN1ADXTnrC6aIU7iMPDU%252BAIP2rWY8xguhvwmC%252B6zyG9foZSDu6SKJQN4kH%252FwtAds7Dj1wsLoqgVrpWppmTmJHWUiZIux97aCldi27mNZ4zZeJWg9QHyfHXNtN6BgrAGzMEt4q3%252Bosv1850g&waid=&is_care_mode=0&is_teenager_mode=0&ad_device_info="
	reqBody = "r=0.9671959261244314&__biz=MzA4MzEwMzAwOA%3D%3D&appmsg_type=9&mid=2649631392&sn=23539f1125a1c1f97ad2272d88c02b3d&idx=1&scene=126&title=%25E4%25BB%25B7%25E6%25A0%25BC%25E8%2585%25B0%25E6%2596%25A9%25EF%25BC%2581%25E7%25AA%2581%25E5%258F%2591%25E5%25A4%25A7%25E8%25B7%25B3%25E6%25B0%25B4&ct=1708955096&abtest_cookie=&devicetype=Windows%2B10%2Bx64&version=6309092b&is_need_ticket=0&is_need_ad=0&comment_id=3344752717986185216&is_need_reward=0&both_ad=0&reward_uin_count=0&send_time=&msg_daily_idx=1&is_original=0&is_only_read=1&req_id=&pass_ticket=hMlJaboOsNIicJ5TxSN9l4ICQCk7GMwfbJ20MyRDOiMGkXUTPZNkxj%25252FHmttBUgqloFA1WAKNdeDU9UoGxdfdJQ%25253D%25253D&is_temp_url=0&item_show_type=0&tmp_version=1&more_read_type=0&appmsg_like_type=2&related_video_sn=&related_video_num=5&vid=&is_pay_subscribe=0&pay_subscribe_uin_count=0&has_red_packet_cover=0&album_id=1296223588617486300&album_video_num=5&cur_album_id=undefined&is_public_related_video=NaN&encode_info_by_base64=undefined&exptype=&export_key=n_ChQIAhIQYWeby2rGgKvQT4os8OoAjhLgAQIE97dBBAEAAAAAAHD2Fb96z0sAAAAOpnltbLcz9gKNyK89dVj0PTDyZ2BDxmasm8TnA4L%252BcLkfgQmw1VDgOXNWNmkSM2%252FfFnqxC6aByoOjVXroePsvdIFQXXvLyzhcwsf6SXfRR9amo4L3ru2WfmR5EcepcpN1ADXTnrC6aIU7iMPDU%252BAIP2rWY8xguhvwmC%252B6zyG9foZSDu6SKJQN4kH%252FwtAds7Dj1wsLoqgVrpWppmTmJHWUiZIux97aCldi27mNZ4zZeJWg9QHyfHXNtN6BgrAGzMEt4q3%252Bosv1850g&export_key_extinfo=&business_type=0"
	gotStatusCode, gotBody, _ := PostWithHeaders(url, reqBody, headers)
	fmt.Println(gotStatusCode)
	fmt.Println(string(gotBody))
}

func TestGet(t *testing.T) {
	url := "https://mp.weixin.qq.com/cgi-bin/appmsgpublish?sub=list&search_field=null&begin=0&count=5&query=&fakeid=MzkwMzUyMjE3OQ%3D%3D&type=101_1&free_publish_type=1&sub_action=list_ex&token=2067135357&lang=zh_CN&f=json&ajax=1"
	url = "https://mp.weixin.qq.com/cgi-bin/appmsgpublish?sub=list&search_field=null&begin=0&count=5&query=&fakeid=MzkwMzUyMjE3OQ%3D%3D&type=101_1&free_publish_type=1&sub_action=list_ex&token=2067135357&lang=zh_CN&f=json&ajax=1"
	headers := make(map[string]string)
	headers["user-agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"

	cookieList := make([]*http.Cookie, 0)
	ck1 := &http.Cookie{Name: "slave_user", Value: "gh_18ab92cf562c"}
	cookieList = append(cookieList, ck1)
	ck2 := &http.Cookie{Name: "slave_sid", Value: "UlVSdWpYREVTaURVQ1p0X1BGS2Ewek9lWGJuc04zWW1WcnViWFUwbnBHcnl6QTFBd2dmYnhIZlUyazZQMlJ5NnFWUVNOWk9UbU9LQlFmOHB4QzN3RW45dFpYZDZrTnNER3dFUWZLRVNBRTNlcl9mNHVvdHh0SWZibjBNSGdiWGUxNDU1aEpRSTVmODBNbmpn"}
	cookieList = append(cookieList, ck2)
	ck3 := &http.Cookie{Name: "bizuin", Value: "3273422465"}
	cookieList = append(cookieList, ck3)
	ck4 := &http.Cookie{Name: "data_bizuin", Value: "3273422465"}
	cookieList = append(cookieList, ck4)
	ck5 := &http.Cookie{Name: "data_ticket", Value: "Q9tsi771B/7EWITNI1KwGNaJL4AuIUosvLl1oNh20h7a3d4WodKJ/pE+nS9MiwWK"}
	cookieList = append(cookieList, ck5)
	ck6 := &http.Cookie{Name: "slave_bizuin", Value: "3273422465"}
	cookieList = append(cookieList, ck6)
	ck7 := &http.Cookie{Name: "rand_info", Value: "CAESIBI+5EwAwk/1fmAyq2tLjsqh4cT1yvPdwcFKwbH6gAEh"}
	cookieList = append(cookieList, ck7)

	gotBody, _ := GetWithCookie(url, cookieList)

	fmt.Println(gotBody)
}
