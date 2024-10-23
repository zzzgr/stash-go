package mt_util

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tidwall/gjson"
	"time"
)

func SleepIfNeed4(name string, advance float64) time.Duration {
	var getTimeRes MtGetTimeResponse

	// 发出请求并获取响应
	client := resty.New()
	start := time.Now()
	res, err := client.R().
		EnableTrace().
		SetResult(&getTimeRes).
		Get("https://cube.meituan.com/ipromotion/cube/toc/component/base/getServerCurrentTime")
	end := time.Now()

	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}

	// 检查HTTP状态码
	if res.IsError() {
		fmt.Printf("HTTP Error: %s\n", res.Status())
		return 0
	}

	// 计算RTT
	rtt := end.Sub(start)
	serverProcessingTime := res.Request.TraceInfo().ServerTime
	networkLatency := (rtt - serverProcessingTime) / 2

	// 计算当前的服务器时间
	now := time.UnixMilli(getTimeRes.Data + networkLatency.Milliseconds())

	// 计算下一个整分钟的时间
	nextMinute := now.Add(time.Minute)
	nextMinute = nextMinute.Truncate(time.Minute)
	nextMinute = nextMinute.Add(time.Millisecond * -time.Duration(advance))
	timeToNextMinute := nextMinute.Sub(now).Milliseconds()

	// 计算本地时间和服务器时间的差异
	diff := float64(time.Now().UnixMilli() - now.UnixMilli())

	// 计算需要休眠的时间
	sleepDuration := time.Duration(timeToNextMinute) * time.Millisecond
	log.Infof("%s [时差: %f] 等待%.3f秒", name, diff, float64(timeToNextMinute)/1000)
	time.Sleep(sleepDuration)
	return time.Duration(-diff) * time.Millisecond
}

func GetUserId(cookie string) string {
	res, _ := resty.New().R().
		SetHeaders(map[string]string{
			"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 TitansX/20.0.1.old KNB/1.0 iOS/17.0 meituangroup/com.meituan.imeituan/12.14.205 meituangroup/12.14.205 App/10110/12.14.205 iPhone/iPhone16,1 WKWebView",
			"Cookie":     cookie,
		}).
		Get("https://open.meituan.com/user/v1/info/auditting?fields=auditAvatarUrl%2CauditUsername")

	return gjson.Get(res.String(), "user.id").String()
}
