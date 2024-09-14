package httpx

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

var NoRedirectClient *resty.Client
var RestyClient = NewRestyClient()
var UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 TitansX/20.0.1.old KNB/1.0 iOS/17.0 meituangroup/com.meituan.imeituan/12.14.205 meituangroup/12.14.205 App/10110/12.14.205 iPhone/iPhone16,1 WKWebView"
var DefaultTimeout = time.Second * 30

func init() {
	NoRedirectClient = resty.New().SetRedirectPolicy(
		resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}),
	)
	NoRedirectClient.SetHeader("User-Agent", UserAgent)
}

func NewRestyClient() *resty.Client {
	return resty.New().
		SetHeader("User-Agent", UserAgent).
		SetRetryCount(3).
		SetTimeout(DefaultTimeout)
}
