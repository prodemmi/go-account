package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"aman/internal/adapters/inbound/http/pages"
	adminpages "aman/internal/adapters/inbound/http/pages/admin"
	authpages "aman/internal/adapters/inbound/http/pages/auth"
	"aman/internal/domains/app"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

var voidTags = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

func check(name, doc string) {
	z := html.NewTokenizer(strings.NewReader(doc))
	var stack []string
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		tok := z.Token()
		switch tt {
		case html.StartTagToken:
			if !voidTags[tok.Data] {
				stack = append(stack, tok.Data)
			}
		case html.EndTagToken:
			if len(stack) == 0 {
				fmt.Printf("%s: stray </%s>\n", name, tok.Data)
				continue
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if top != tok.Data {
				fmt.Printf("%s: </%s> closes <%s>\n", name, tok.Data, top)
			}
		}
	}
	if len(stack) > 0 {
		fmt.Printf("%s: unclosed %v\n", name, stack)
	}
	for _, bad := range []string{"`))", "<%!", "fmt.Sprintf"} {
		if strings.Contains(doc, bad) {
			fmt.Printf("%s: leaked %q\n", name, bad)
		}
	}
	fmt.Printf("%s: %d bytes ok\n", name, len(doc))
}

func render(c interface {
	Render(context.Context, interface{ Write([]byte) (int, error) }) error
}) string {
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		fmt.Println("render error:", err)
		os.Exit(1)
	}
	return buf.String()
}

type writer struct{ bytes.Buffer }

func (w *writer) Write(p []byte) (int, error) { return w.Buffer.Write(p) }

func main() {
	data := adminpages.AppFormData{
		Name:                       "فوتومی",
		Url:                        "https://example.com",
		CallbackURL:                "https://example.com/cb",
		OTPChannel:                 "sms",
		OTPLength:                  5,
		OTPExpire:                  120,
		AuthPageEnabled:            true,
		AuthPagePrimaryColor:       "#2563eb",
		AuthPageBackgroundColor:    "#eef2f7",
		AuthPageCardBackgroundColor: "#ffffff",
		AuthPagePinBackgroundColor: "#f8fafc",
		AuthPagePinTextColor:       "#0f172a",
		AuthPageTheme:              "light",
		AuthPageFont:               "vazirmatn",
		AuthPageLayout:             "split",
		AuthPageLogoPosition:       "top",
		AuthPageBannerMode:         "progress",
		AuthPageLoginTitle:         "ورود به {app}",
		AuthPageLoginSubtitle:      "شماره موبایل خود را وارد کنید.",
		AuthPageIdentifierLabel:    "شماره موبایل",
		AuthPageSendButton:         "دریافت کد",
		AuthPageVerifyTitle:        "کد تایید را وارد کنید",
		AuthPageVerifySubtitle:     "کد ارسال‌شده به {identifier} را وارد کنید.",
		AuthPageVerifyButton:       "تایید و ورود",
		AuthPageRetryText:          "ارسال دوباره کد",
		AuthPageBackText:           "تغییر ایمیل یا شماره موبایل",
		AuthPageFooterText:         "ورود امن با Aman CAS",
		AuthPageBanners:            []app.Banner{{URL: "https://example.com/b.svg", Text: "خوش آمدید", Duration: 8}},
	}
	check("AppForm create", render(adminpages.AppForm(data, false)))
	check("AppForm edit", render(adminpages.AppForm(data, true)))

	authData := pages.AuthData{
		App:             app.App{ID: uuid.New(), Name: "فوتومی"},
		AppID:           uuid.New().String(),
		PageTitle:       "ورود",
		DocumentTitle:   "ورود به فوتومی",
		PrimaryColor:    "#2563eb",
		BackgroundColor: "#eef2f7",
		Theme:           "light",
		Font:            "vazirmatn",
		Layout:          "split",
		LogoPosition:    "top",
		BannerMode:      "progress",
		AllowMobile:     true,
		CopyCustom:      true,
		Step:            "login",
		Title:           "ورود به فوتومی",
		Subtitle:        "شماره موبایل خود را وارد کنید.",
		IdentifierLabel: "شماره موبایل",
		SendButton:      "دریافت کد",
		VerifyTitle:     "کد تایید را وارد کنید",
		VerifySubtitle:  "کد ارسال‌شده به {identifier} را وارد کنید.",
		VerifyButton:    "تایید و ورود",
		RetryText:       "ارسال دوباره کد",
		BackText:        "تغییر ایمیل یا شماره موبایل",
		FooterText:      "ورود امن با Aman CAS",
		OTPLength:       5,
		OTPInputType:    "tel",
		OTPCodeInputType: "tel",
		OTPExpire:       120,
		Banners:         []app.Banner{{URL: "https://example.com/b.svg", Text: "خوش آمدید", Duration: 8}},
	}
	check("Auth login", render(authpages.Login(authData)))
}
