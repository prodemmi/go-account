package pages

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"
	"aman/internal/domains/app"

	"github.com/a-h/templ"
)

// AppForm renders the Persian create/edit app form. Every setting an app can
// carry is exposed here, with a side column for the settings guide.
func AppForm(data AppFormData, isEdit bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		title := "برنامه جدید"
		subtitle := "برنامه تازه‌ای را به سرویس احراز هویت وصل کنید."
		action := "/admin/apps"
		if isEdit {
			title = ""
			subtitle = ""
			action = "/admin/apps/" + data.ID + "/update"
		}

		hostedHidden := " hidden"
		defaultHidden := ""
		if data.AuthPageEnabled {
			hostedHidden = ""
			defaultHidden = " hidden"
		}

		// The back button comes first, then the title/subtitle, and the save
		// actions are the last item of the header.
		w.Write([]byte(`<section class="page-head form-head">
	<a href="/admin/apps" class="btn btn-ghost btn-back">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M9 5l7 7-7 7"/></svg>
		بازگشت
	</a>
	<div class="form-head-title">
		<h2>`))
		pages.WriteEscaped(w, title)
		w.Write([]byte(`</h2>
		<p>`))
		pages.WriteEscaped(w, subtitle)
		w.Write([]byte(`</p>
	</div>
	<div class="form-head-actions">
		<a href="/admin/apps" class="btn btn-ghost">انصراف</a>
		<button type="submit" form="app-form" class="btn btn-primary">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12.5l4.5 4.5L19 7.5"/></svg>
			ذخیره برنامه
		</button>
	</div>
</section>`))

		if data.Error != "" || data.Success != "" {
			if err := flash(data.Error, data.Success).Render(ctx, w); err != nil {
				return err
			}
		}

		w.Write([]byte(`<form method="POST" action="`))
		pages.WriteEscaped(w, action)
		w.Write([]byte(`" class="form-layout" id="app-form">
	<input type="hidden" name="id" value="`))
		pages.WriteEscaped(w, data.ID)
		if isEdit {
			w.Write([]byte(`" />
	<input type="hidden" name="auth_page_copy_custom" value="1" />`))
		} else {
			w.Write([]byte(`" />`))
		}
		w.Write([]byte(`

	<div class="form-main">

		<section class="card form-section" data-guide-anchor="app-basics">
			<div class="section-head">
				<div>
					<h3>اطلاعات برنامه</h3>
					<p>نام نمایشی، آدرس اصلی و مقصد بازگشت کاربر پس از ورود.</p>
				</div>
			</div>

			<div class="field-grid">
				<div class="field col-2">
					<label for="name">نام برنامه</label>
					<input type="text" id="name" name="name" value="`))
		pages.WriteEscaped(w, data.Name)
		w.Write([]byte(`" placeholder="فوتومی" required />
					<p class="hint">این نام در صفحه ورود، فهرست برنامه‌ها و نشست‌ها نمایش داده می‌شود.</p>
				</div>
				<div class="field col-2">
					<label for="url">آدرس اصلی برنامه</label>
					<input type="url" id="url" name="url" value="`))
		pages.WriteEscaped(w, data.Url)
		w.Write([]byte(`" placeholder="https://fotomi.net/" dir="ltr" class="ltr" />
					<p class="hint">مبنای اعتبارسنجی بازگشت خودکار به برنامه است؛ فقط آدرس‌هایی با همین دامنه پذیرفته می‌شوند.</p>
				</div>
				<div class="field col-2">
					<label for="callback_url">آدرس بازگشت (callback)</label>
					<input type="url" id="callback_url" name="callback_url" value="`))
		pages.WriteEscaped(w, data.CallbackURL)
		w.Write([]byte(`" placeholder="https://fotomi.net/auth/callback" dir="ltr" class="ltr" required />
					<p class="hint">پس از ورود موفق، کاربر با یک تیکت یکبار مصرف (OTT) به این آدرس بازمی‌گردد.</p>
				</div>
				<div class="field col-2" data-guide-anchor="webhook">
					<label for="webhook_url">آدرس وب‌هوک <span class="label-optional">اختیاری</span></label>
					<input type="url" id="webhook_url" name="webhook_url" value="`))
		pages.WriteEscaped(w, data.WebhookURL)
		w.Write([]byte(`" placeholder="https://fotomi.net/webhooks/cas" dir="ltr" class="ltr" />
					<p class="hint">پس از هر رویداد (ارسال کد، ورود، خروج، ابطال نشست) یک درخواست JSON به این آدرس ارسال می‌شود.</p>
				</div>
			</div>
		</section>

		<section class="card form-section" data-guide-anchor="otp">
			<div class="section-head">
				<div>
					<h3>کد یکبار مصرف</h3>
					<p>کانال ارسال، طول و ترکیب کد، و مدت اعتبار کد و تیکت ورود.</p>
				</div>
			</div>

			<div class="field-grid">
				<div class="field col-2">
					<label for="otp_channel">روش ارسال کد</label>
					<select id="otp_channel" name="otp_channel">`))
		writeOption(w, app.OTPChannelSMS, "پیامک (موبایل)", data.OTPChannel)
		writeOption(w, app.OTPChannelEmail, "ایمیل", data.OTPChannel)
		writeOption(w, app.OTPChannelBoth, "موبایل و ایمیل", data.OTPChannel)
		w.Write([]byte(`</select>
					<p class="hint">با انتخاب «موبایل و ایمیل»، نوع ورودی صفحه ورود بر اساس نوع شناسه تشخیص داده می‌شود.</p>
				</div>
				<div class="field">
					<label for="otp_length">طول کد</label>
					<div class="input-suffix">
						<input type="number" id="otp_length" name="otp_length" value="`))
		pages.WriteEscaped(w, itostr(data.OTPLength))
		w.Write([]byte(`" placeholder="5" min="4" max="32" required />
						<span>رقم/حرف</span>
					</div>
					<p class="hint">بین ۴ تا ۳۲ کاراکتر.</p>
				</div>
				<div class="field">
					<label for="otp_charset">ترکیب کد</label>
					<select id="otp_charset" name="otp_charset">`))
		writeOption(w, app.OTPCharsetNumeric, "فقط عدد", data.OTPCharset)
		writeOption(w, app.OTPCharsetAlpha, "فقط حرف", data.OTPCharset)
		writeOption(w, app.OTPCharsetAlphanumeric, "حرف و عدد", data.OTPCharset)
		w.Write([]byte(`</select>
					<p class="hint">کدهای عددی برای صفحه کلید موبایل ساده‌تر هستند.</p>
				</div>
				<div class="field">
					<label for="otp_expire">اعتبار کد</label>
					<div class="input-suffix">
						<input type="number" id="otp_expire" name="otp_expire" value="`))
		pages.WriteEscaped(w, itostr(data.OTPExpire))
		w.Write([]byte(`" placeholder="120" min="1" />
						<span>ثانیه</span>
					</div>
				</div>
				<div class="field">
					<label for="ott_expire">اعتبار تیکت یکبار مصرف</label>
					<div class="input-suffix">
						<input type="number" id="ott_expire" name="ott_expire" value="`))
		pages.WriteEscaped(w, itostr(data.OTTExpire))
		w.Write([]byte(`" placeholder="60" min="1" />
						<span>ثانیه</span>
					</div>
					<p class="hint">کوتاه نگه دارید؛ تیکت فقط یک‌بار قابل استفاده است.</p>
				</div>
			</div>
		</section>

		<section class="card form-section" data-guide-anchor="sessions">
			<div class="section-head">
				<div>
					<h3>نشست و کوکی</h3>
					<p>مدت اعتبار نشست و تنظیمات کوکی اختصاصی همین برنامه.</p>
				</div>
			</div>

			<div class="field-grid">
				<div class="field" data-disabled-when="#infinite_session">
					<label for="session_timeout">مدت نشست</label>
					<div class="input-suffix">
						<input type="number" id="session_timeout" name="session_timeout" value="`))
		pages.WriteEscaped(w, itostr(data.SessionTimeout))
		w.Write([]byte(`" placeholder="3600" min="1" />
						<span>ثانیه</span>
					</div>
					<p class="hint">با فعال بودن نشست همیشگی نادیده گرفته می‌شود.</p>
				</div>
				<div class="field">
					<label for="session_cookie_name">نام کوکی نشست</label>
					<input type="text" id="session_cookie_name" name="session_cookie_name" value="`))
		pages.WriteEscaped(w, data.SessionCookieName)
		w.Write([]byte(`" placeholder="cas_cookie" dir="ltr" />
					<p class="hint">اگر خالی بماند، مقدار پیش‌فرض <code>cas_cookie</code> استفاده می‌شود.</p>
				</div>
				<div class="field">
					<label for="session_cookie_same_site">SameSite کوکی</label>
					<select id="session_cookie_same_site" name="session_cookie_same_site">`))
		writeOption(w, app.SameSiteLax, "Lax (پیش‌فرض)", data.SessionCookieSameSite)
		writeOption(w, app.SameSiteStrict, "Strict", data.SessionCookieSameSite)
		writeOption(w, app.SameSiteNone, "None", data.SessionCookieSameSite)
		writeOption(w, app.SameSiteDefault, "پیش‌فرض مرورگر", data.SessionCookieSameSite)
		w.Write([]byte(`</select>
					<p class="hint">برای ورود از دامنه‌های مختلف معمولاً <code>None</code> همراه با Secure لازم است.</p>
				</div>
			</div>

			<label class="switch-row" for="infinite_session">
				<span class="switch-copy">
					<strong>نشست همیشگی</strong>
					<small>نشست بدون تاریخ انقضای عملی ساخته می‌شود و مدت نشست نادیده گرفته می‌شود.</small>
				</span>
				<span class="switch">
					<input type="checkbox" id="infinite_session" name="infinite_session" value="true"`))
		if data.InfiniteSession {
			w.Write([]byte(` checked`))
		}
		w.Write([]byte(` />
					<span class="switch-track" aria-hidden="true"></span>
				</span>
			</label>

			<label class="switch-row" for="session_cookie_http_only">
				<span class="switch-copy">
					<strong>کوکی HttpOnly</strong>
					<small>دسترسی جاوااسکریپت صفحه به کوکی نشست را می‌بندد. برای امنیت بیشتر روشن بماند.</small>
				</span>
				<span class="switch">
					<input type="checkbox" id="session_cookie_http_only" name="session_cookie_http_only" value="true"`))
		if data.SessionCookieHTTPOnly {
			w.Write([]byte(` checked`))
		}
		w.Write([]byte(` />
					<span class="switch-track" aria-hidden="true"></span>
				</span>
			</label>

			<label class="switch-row" for="session_cookie_secure">
				<span class="switch-copy">
					<strong>کوکی Secure</strong>
					<small>کوکی فقط روی اتصال HTTPS ارسال می‌شود. در محیط توسعه بدون HTTPS خاموش بماند.</small>
				</span>
				<span class="switch">
					<input type="checkbox" id="session_cookie_secure" name="session_cookie_secure" value="true"`))
		if data.SessionCookieSecure {
			w.Write([]byte(` checked`))
		}
		w.Write([]byte(` />
					<span class="switch-track" aria-hidden="true"></span>
				</span>
			</label>

			<label class="switch-row" for="redirect_back">
				<span class="switch-copy">
					<strong>بازگشت به آخرین آدرس برنامه</strong>
					<small>پس از ورود، اگر آدرس مبدأ متعلق به دامنه همین برنامه باشد کاربر به همان صفحه بازمی‌گردد.</small>
				</span>
				<span class="switch">
					<input type="checkbox" id="redirect_back" name="redirect_back" value="true"`))
		if data.RedirectBack {
			w.Write([]byte(` checked`))
		}
		w.Write([]byte(` />
					<span class="switch-track" aria-hidden="true"></span>
				</span>
			</label>
		</section>

		<section class="card form-section" data-guide-anchor="geo-restrictions">
			<div class="section-head">
				<div>
					<h3>محدودیت جغرافیایی</h3>
					<p>در صورت تعیین فهرست، درخواست ورود فقط از کشورهای انتخاب‌شده پذیرفته می‌شود.</p>
				</div>
			</div>
			<div class="field">
				<label for="allowed_countries">کشورهای مجاز (کد ISO دوحرفی)</label>
				<input type="text" id="allowed_countries" name="allowed_countries" value="`))
		pages.WriteEscaped(w, data.AllowedCountries)
		w.Write([]byte(`" placeholder="CA,US" dir="ltr" class="ltr" />
				<p class="hint">با ویرگول جدا کنید؛ خالی یعنی بدون محدودیت. برای محدودیت و تکمیل مکان، فایل محلی MaxMind GeoLite2 City را با <code>GEOIP_DATABASE_PATH</code> پیکربندی کنید؛ در نبود پایگاه داده، ورود محدودشده رد می‌شود.</p>
			</div>
		</section>

		<section class="card form-section" data-guide-anchor="auth-page">
			<div class="section-head">
				<div>
					<h3>صفحه ورود</h3>
					<p>انتخاب جریان احراز هویت و تنظیم صفحه‌های مرورگر CAS.</p>
				</div>
			</div>

			<label class="switch-row" for="custom_auth_enabled">
				<span class="switch-copy">
					<strong>حالت احراز هویت سفارشی از طریق API</strong>
					<small>رابط ورود را خودتان می‌سازید؛ همه مسیرهای مرورگری <code>/auth/**</code> برای این برنامه ۴۰۴ می‌دهند. از APIهای <code>/v1/auth/*</code> استفاده کنید؛ تنظیم صفحه میزبانی‌شده در این حالت کاربرد ندارد.</small>
				</span>
				<span class="switch">
					<input type="checkbox" id="custom_auth_enabled" name="custom_auth_enabled" value="true"`))
		if data.CustomAuthEnabled {
			w.Write([]byte(` checked`))
		}
		w.Write([]byte(` />
					<span class="switch-track" aria-hidden="true"></span>
				</span>
			</label>

			<label class="switch-row" for="auth_page_enabled">
				<span class="switch-copy">
					<strong>استفاده از صفحه ورود میزبانی‌شده در برنامه</strong>
					<small>در صورت فعال بودن، کاربر برای ورود به آدرس سفارشی زیر فرستاده می‌شود و تنظیمات ظاهری صفحه پیش‌فرض پنهان می‌شوند.</small>
				</span>
				<span class="switch">
					<input type="checkbox" id="auth_page_enabled" name="auth_page_enabled" value="true"`))
		if data.AuthPageEnabled {
			w.Write([]byte(` checked`))
		}
		w.Write([]byte(` />
					<span class="switch-track" aria-hidden="true"></span>
				</span>
			</label>

			<div class="field" data-show-when="#auth_page_enabled"`))
		w.Write([]byte(hostedHidden))
		w.Write([]byte(`>
				<label for="auth_page_url">آدرس صفحه ورود سفارشی</label>
				<input type="url" id="auth_page_url" name="auth_page_url" value="`))
		pages.WriteEscaped(w, data.AuthPageURL)
		w.Write([]byte(`" placeholder="https://fotomi.net/auth/login" />
				<p class="hint">آدرس صفحه‌ای در برنامه شما که فرم ورود را نمایش می‌دهد. CAS پارامترهای لازم را به آن اضافه می‌کند.</p>
			</div>

			<div data-hide-when="#auth_page_enabled" data-guide-anchor="branding"`))
		w.Write([]byte(defaultHidden))
		w.Write([]byte(`>
				<div class="subsection-head">
					<h4>ظاهر صفحه پیش‌فرض</h4>
					<p>چیدمان، فونت، رنگ‌ها و لوگوی صفحه‌ای که CAS برای این برنامه نمایش می‌دهد.</p>
				</div>

				<div class="field-grid">
					<div class="field">
						<label for="auth_page_layout">چیدمان</label>
						<select id="auth_page_layout" name="auth_page_layout">`))
		writeOption(w, app.LayoutCard, "کارت مرکزی", data.AuthPageLayout)
		writeOption(w, app.LayoutSplit, "دو ستونه", data.AuthPageLayout)
		writeOption(w, app.LayoutMinimal, "مینیمال", data.AuthPageLayout)
		w.Write([]byte(`</select>
					</div>
					<div class="field">
						<label for="auth_page_font">فونت</label>
						<select id="auth_page_font" name="auth_page_font">`))
		writeOption(w, app.FontVazirmatn, "وزیرمتن (Vazirmatn)", data.AuthPageFont)
		writeOption(w, app.FontIRANSans, "ایران‌سنس (IRANSansX)", data.AuthPageFont)
		writeOption(w, app.FontEstedad, "استعداد (Estedad)", data.AuthPageFont)
		writeOption(w, app.FontSahel, "ساحل (Sahel)", data.AuthPageFont)
		writeOption(w, app.FontInter, "اینتر (Inter - English)", data.AuthPageFont)
		writeOption(w, app.FontManrope, "منروپ (Manrope - English)", data.AuthPageFont)
		writeOption(w, app.FontSystem, "فونت سیستم", data.AuthPageFont)
		w.Write([]byte(`</select>
					</div>
					<div class="field">
						<label for="auth_page_text_align">تراز متن</label>
						<select id="auth_page_text_align" name="auth_page_text_align">`))
		writeOption(w, app.AlignCenter, "وسط‌چین", data.AuthPageTextAlign)
		writeOption(w, app.AlignStart, "راست‌چین (شروع)", data.AuthPageTextAlign)
		writeOption(w, app.AlignEnd, "چپ‌چین (پایان)", data.AuthPageTextAlign)
		w.Write([]byte(`</select>
					</div>
					<div class="field">
						<label for="auth_page_theme">پوسته</label>
						<select id="auth_page_theme" name="auth_page_theme">`))
		writeOption(w, app.ThemeLight, "روشن", data.AuthPageTheme)
		writeOption(w, app.ThemeDark, "تیره", data.AuthPageTheme)
		w.Write([]byte(`</select>
					</div>
					<div class="field">
						<label for="auth_page_primary_color">رنگ اصلی</label>
						<div class="color-field">
							<input type="color" id="auth_page_primary_color" name="auth_page_primary_color" value="`))
		pages.WriteEscaped(w, data.AuthPagePrimaryColor)
		w.Write([]byte(`" />
							<span class="color-value ltr">`))
		pages.WriteEscaped(w, data.AuthPagePrimaryColor)
		w.Write([]byte(`</span>
						</div>
					</div>
					<div class="field">
						<label for="auth_page_background_color">رنگ پس‌زمینه صفحه</label>
						<div class="color-field">
							<input type="color" id="auth_page_background_color" name="auth_page_background_color" value="`))
		pages.WriteEscaped(w, data.AuthPageBackgroundColor)
		w.Write([]byte(`" />
							<span class="color-value ltr">`))
		pages.WriteEscaped(w, data.AuthPageBackgroundColor)
		w.Write([]byte(`</span>
						</div>
					</div>
					<div class="field">
						<label for="auth_page_card_background_color">رنگ پس‌زمینه کارت ورود</label>
						<div class="color-field">
							<input type="color" id="auth_page_card_background_color" name="auth_page_card_background_color" value="`))
		pages.WriteEscaped(w, data.AuthPageCardBackgroundColor)
		w.Write([]byte(`" />
							<span class="color-value ltr">`))
		pages.WriteEscaped(w, data.AuthPageCardBackgroundColor)
		w.Write([]byte(`</span>
						</div>
						<p class="hint">رنگ سطح کارتی که فرم ورود داخل آن قرار می‌گیرد.</p>
					</div>
					<div class="field">
						<label for="auth_page_text_color">رنگ متن اصلی</label>
						<div class="color-field">
							<input type="color" id="auth_page_text_color" name="auth_page_text_color" value="`))
		pages.WriteEscaped(w, data.AuthPageTextColor)
		w.Write([]byte(`" />
							<span class="color-value ltr">`))
		pages.WriteEscaped(w, data.AuthPageTextColor)
		w.Write([]byte(`</span>
						</div>
					</div>
					<div class="field">
						<label for="auth_page_muted_text_color">رنگ متن ثانویه</label>
						<div class="color-field">
							<input type="color" id="auth_page_muted_text_color" name="auth_page_muted_text_color" value="`))
		pages.WriteEscaped(w, data.AuthPageMutedTextColor)
		w.Write([]byte(`" />
							<span class="color-value ltr">`))
		pages.WriteEscaped(w, data.AuthPageMutedTextColor)
		w.Write([]byte(`</span>
						</div>
					</div>
					<div class="field">
						<label for="auth_page_pin_background_color">رنگ پس‌زمینه خانه‌های کد</label>
						<div class="color-field">
							<input type="color" id="auth_page_pin_background_color" name="auth_page_pin_background_color" value="`))
		pages.WriteEscaped(w, data.AuthPagePinBackgroundColor)
		w.Write([]byte(`" />
							<span class="color-value ltr">`))
		pages.WriteEscaped(w, data.AuthPagePinBackgroundColor)
		w.Write([]byte(`</span>
						</div>
						<p class="hint">پس‌زمینه جعبه‌های وارد کردن کد یکبار مصرف.</p>
					</div>
					<div class="field">
						<label for="auth_page_pin_text_color">رنگ متن خانه‌های کد</label>
						<div class="color-field">
							<input type="color" id="auth_page_pin_text_color" name="auth_page_pin_text_color" value="`))
		pages.WriteEscaped(w, data.AuthPagePinTextColor)
		w.Write([]byte(`" />
							<span class="color-value ltr">`))
		pages.WriteEscaped(w, data.AuthPagePinTextColor)
		w.Write([]byte(`</span>
						</div>
						<p class="hint">رنگ ارقام یا حروف داخل خانه‌های کد.</p>
					</div>
					<div class="field col-2">
						<label for="auth_page_logo">آدرس لوگو</label>
						<input type="url" id="auth_page_logo" name="auth_page_logo" value="`))
		pages.WriteEscaped(w, data.AuthPageLogo)
		w.Write([]byte(`" placeholder="https://fotomi.net/logo.png" dir="ltr" class="ltr" />
						<p class="hint">اگر خالی بماند، حرف اول نام برنامه به‌جای لوگو نمایش داده می‌شود.</p>
					</div>
					<div class="field">
						<label for="auth_page_logo_mask">شکل لوگو</label>
						<select id="auth_page_logo_mask" name="auth_page_logo_mask">`))
		writeOption(w, app.LogoMaskNone, "بدون ماسک", data.AuthPageLogoMask)
		writeOption(w, app.LogoMaskCircle, "دایره", data.AuthPageLogoMask)
		writeOption(w, app.LogoMaskRounded, "گوشه‌گرد", data.AuthPageLogoMask)
		writeOption(w, app.LogoMaskSquare, "مربع", data.AuthPageLogoMask)
		w.Write([]byte(`</select>
					</div>
					<div class="field">
						<label for="auth_page_logo_position">جای لوگو</label>
						<select id="auth_page_logo_position" name="auth_page_logo_position">`))
		writeOption(w, app.LogoPositionTop, "بالای کارت (بیرون)", data.AuthPageLogoPosition)
		writeOption(w, app.LogoPositionCenter, "بالای فرم (داخل کارت)", data.AuthPageLogoPosition)
		writeOption(w, app.LogoPositionStart, "ابتدای عنوان (راست)", data.AuthPageLogoPosition)
		writeOption(w, app.LogoPositionEnd, "انتهای عنوان (چپ)", data.AuthPageLogoPosition)
		writeOption(w, app.LogoPositionBottom, "پایین کارت", data.AuthPageLogoPosition)
		w.Write([]byte(`</select>
					</div>
				</div>

				<label class="switch-row" for="auth_page_captcha_enabled" data-guide-anchor="captcha">
					<span class="switch-copy">
						<strong>کپچا</strong>
						<small>پیش از ارسال کد، ارائه‌دهنده انتخاب‌شده بررسی می‌شود. کلیدهای هر برنامه جداگانه نگهداری می‌شوند.</small>
					</span>
					<span class="switch">
						<input type="checkbox" id="auth_page_captcha_enabled" name="auth_page_captcha_enabled" value="true"`))
		if data.AuthPageCaptchaEnabled {
			w.Write([]byte(` checked`))
		}
		w.Write([]byte(` />
						<span class="switch-track" aria-hidden="true"></span>
					</span>
				</label>

				<div class="field-grid captcha-fields" data-show-when="#auth_page_captcha_enabled" data-guide-anchor="captcha">
					<div class="field">
						<label for="auth_page_captcha_provider">ارائه‌دهنده کپچا</label>
						<select id="auth_page_captcha_provider" name="auth_page_captcha_provider">`))
		writeOption(w, app.CaptchaProviderARCaptcha, "ArCaptcha", data.AuthPageCaptchaProvider)
		writeOption(w, app.CaptchaProviderGoogle, "Google reCAPTCHA v2", data.AuthPageCaptchaProvider)
		w.Write([]byte(`</select>
						<p class="hint">برای reCAPTCHA از کلیدهای نوع checkbox (v2) استفاده کنید.</p>
					</div>
					<div class="field">
						<label for="auth_page_captcha_site_key">Site Key</label>
						<input type="text" id="auth_page_captcha_site_key" name="auth_page_captcha_site_key" value="`))
		pages.WriteEscaped(w, data.AuthPageCaptchaSiteKey)
		w.Write([]byte(`" placeholder="site key" dir="ltr" class="ltr" />
						<p class="hint">کلید عمومی که در ویجت کپچا در صفحه ورود استفاده می‌شود.</p>
					</div>
					<div class="field">
						<label for="auth_page_captcha_secret_key">Secret Key</label>
						<input type="text" id="auth_page_captcha_secret_key" name="auth_page_captcha_secret_key" value="`))
		pages.WriteEscaped(w, data.AuthPageCaptchaSecretKey)
		w.Write([]byte(`" placeholder="secret key" dir="ltr" class="ltr" />
						<p class="hint">کلید خصوصی برای تایید سمت سرور؛ فقط در پنل مدیریت نگهداری می‌شود.</p>
					</div>
				</div>

				<div class="subsection-head">
					<h4>بنرهای صفحه ورود</h4>
					<p>هر بنر برای مدت مشخصی نمایش داده می‌شود و سپس بنر بعدی جای آن را می‌گیرد. ترتیب را با کشیدن دسته هر کارت تغییر دهید.</p>
				</div>

				<div class="field-grid">
					<div class="field">
						<label for="auth_page_banner_mode">نمایش زمان بنر</label>
						<select id="auth_page_banner_mode" name="auth_page_banner_mode">`))
		writeOption(w, app.BannerModeNone, "بدون نمایش زمان", data.AuthPageBannerMode)
		writeOption(w, app.BannerModeProgress, "نوار پیشرفت", data.AuthPageBannerMode)
		writeOption(w, app.BannerModeTime, "زمان باقی‌مانده", data.AuthPageBannerMode)
		w.Write([]byte(`</select>
						<p class="hint">نحوه نمایش مدت باقی‌مانده بنر جاری روی صفحه ورود.</p>
					</div>
				</div>

				<div class="banner-editor" data-banner-list>`))

		if len(data.AuthPageBanners) == 0 {
			data.AuthPageBanners = []app.Banner{{Duration: 6}}
		}
		for i, banner := range data.AuthPageBanners {
			if err := bannerRow(i, banner).Render(ctx, w); err != nil {
				return err
			}
		}

		w.Write([]byte(`</div>
				<div class="banner-actions">
					<button type="button" class="btn btn-ghost" data-banner-add>
						<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M12 5v14M5 12h14"/></svg>
						افزودن بنر
					</button>
				</div>
				<template data-banner-template>`))
		if err := bannerRow(0, app.Banner{Duration: 6}).Render(ctx, w); err != nil {
			return err
		}
		w.Write([]byte(`</template>

				<div class="subsection-head">
					<h4>متن‌های صفحه ورود</h4>
					<p>عبارت‌ها را با زبان برند خودتان بنویسید. خالی گذاشتن هر مورد در ویرایش، آن متن را از صفحه ورود حذف می‌کند.</p>
				</div>

				<div class="field-grid">
					<div class="field">
						<label for="auth_page_title">عنوان سند مرورگر</label>
						<input type="text" id="auth_page_title" name="auth_page_title" value="`))
		pages.WriteEscaped(w, data.AuthPageTitle)
		w.Write([]byte(`" placeholder="ورود به {app}" />
						<p class="hint">در تب مرورگر نمایش داده می‌شود. عبارت <code>{app}</code> با نام برنامه جایگزین می‌شود.</p>
					</div>
					<div class="field">
						<label for="auth_page_login_title">عنوان صفحه ورود</label>
						<input type="text" id="auth_page_login_title" name="auth_page_login_title" value="`))
		pages.WriteEscaped(w, data.AuthPageLoginTitle)
		w.Write([]byte(`" placeholder="ورود به {app}" />
					</div>
					<div class="field">
						<label for="auth_page_login_subtitle">زیرعنوان صفحه ورود</label>
						<input type="text" id="auth_page_login_subtitle" name="auth_page_login_subtitle" value="`))
		pages.WriteEscaped(w, data.AuthPageLoginSubtitle)
		w.Write([]byte(`" placeholder="برای دریافت کد ورود، شناسه خود را وارد کنید." />
					</div>
					<div class="field">
						<label for="auth_page_identifier_label">فیلد شناسه</label>
						<input type="text" id="auth_page_identifier_label" name="auth_page_identifier_label" value="`))
		pages.WriteEscaped(w, data.AuthPageIdentifierLabel)
		w.Write([]byte(`" placeholder="ایمیل یا شماره موبایل" />
					</div>
					<div class="field">
						<label for="auth_page_send_button">متن دکمه ارسال کد</label>
						<input type="text" id="auth_page_send_button" name="auth_page_send_button" value="`))
		pages.WriteEscaped(w, data.AuthPageSendButton)
		w.Write([]byte(`" placeholder="دریافت کد ورود" />
					</div>
					<div class="field">
						<label for="auth_page_verify_title">عنوان مرحله تایید کد</label>
						<input type="text" id="auth_page_verify_title" name="auth_page_verify_title" value="`))
		pages.WriteEscaped(w, data.AuthPageVerifyTitle)
		w.Write([]byte(`" placeholder="کد تایید را وارد کنید" />
					</div>
					<div class="field">
						<label for="auth_page_verify_subtitle">زیرعنوان مرحله تایید</label>
						<input type="text" id="auth_page_verify_subtitle" name="auth_page_verify_subtitle" value="`))
		pages.WriteEscaped(w, data.AuthPageVerifySubtitle)
		w.Write([]byte(`" placeholder="کد ارسال‌شده به {identifier} را وارد کنید." />
					</div>
					<div class="field">
						<label for="auth_page_verify_button">متن دکمه تایید</label>
						<input type="text" id="auth_page_verify_button" name="auth_page_verify_button" value="`))
		pages.WriteEscaped(w, data.AuthPageVerifyButton)
		w.Write([]byte(`" placeholder="تایید و ورود" />
					</div>
					<div class="field">
						<label for="auth_page_retry_text">متن ارسال دوباره</label>
						<input type="text" id="auth_page_retry_text" name="auth_page_retry_text" value="`))
		pages.WriteEscaped(w, data.AuthPageRetryText)
		w.Write([]byte(`" placeholder="ارسال دوباره کد" />
					</div>
					<div class="field">
						<label for="auth_page_back_text">متن بازگشت</label>
						<input type="text" id="auth_page_back_text" name="auth_page_back_text" value="`))
		pages.WriteEscaped(w, data.AuthPageBackText)
		w.Write([]byte(`" placeholder="تغییر ایمیل یا شماره موبایل" />
					</div>
					<div class="field">
						<label for="auth_page_footer_text">متن پانوشت صفحه</label>
						<input type="text" id="auth_page_footer_text" name="auth_page_footer_text" value="`))
		pages.WriteEscaped(w, data.AuthPageFooterText)
		w.Write([]byte(`" placeholder="ورود امن با Aman CAS" />
					</div>
				</div>
			</div>
		</section>
	</div>

	<aside class="form-side">
		<section class="card side-card">
			<h3>راهنمای تنظیمات</h3>
			<a href="/admin/guide" class="btn btn-ghost btn-block" target="_blank" rel="noopener noreferrer">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M4.5 5.5A2 2 0 0 1 6.5 3.5H19v17H6.5a2 2 0 0 1-2-2z"/><path d="M8 7.5h7M8 11h7M8 14.5h4"/></svg>
				مطالعه راهنما
			</a>
			<a href="/admin/developers" class="btn btn-ghost btn-block" target="_blank" rel="noopener noreferrer">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M8.5 8.5L5 12l3.5 3.5M15.5 8.5L19 12l-3.5 3.5M13 6l-2 12"/></svg>
				مستندات API
			</a>
		</section>
	</aside>

	<script>
		(function () {
			document.querySelectorAll('#app-form .field, #app-form .switch-row').forEach(function (field) {
				var section = field.closest('[data-guide-anchor]');
				if (!section || field.querySelector('.field-guide-link')) return;
				var isSwitch = field.classList.contains('switch-row');
				var label = isSwitch ? field.querySelector('.switch-copy strong') : field.querySelector('label');
				if (!label) return;

				var link = document.createElement('a');
				link.className = 'field-guide-link';
				link.href = '/admin/guide#' + encodeURIComponent(section.getAttribute('data-guide-anchor'));
				link.target = '_blank';
				link.rel = 'noopener noreferrer';
				link.textContent = 'راهنما ↗';
				link.setAttribute('aria-label', 'راهنمای ' + label.textContent.trim());

				if (isSwitch) {
					var row = document.createElement('div');
					row.className = 'switch-guide-row';
					field.parentNode.insertBefore(row, field);
					row.appendChild(field);
					row.appendChild(link);
				} else {
					var labelRow = document.createElement('div');
					labelRow.className = 'field-label-row';
					field.insertBefore(labelRow, label);
					labelRow.appendChild(label);
					labelRow.appendChild(link);
				}
			});
		})();
	</script>
</form>`))
		return nil
	})
}

func bannerRow(index int, banner app.Banner) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		w.Write([]byte(`<div class="banner-row" data-banner-row>
	<span class="drag-handle" draggable="true" data-drag-handle title="برای جابه‌جایی بکشید" aria-label="جابه‌جایی بنر">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M8 7h.01M8 12h.01M8 17h.01M16 7h.01M16 12h.01M16 17h.01"/></svg>
	</span>
	<div class="banner-row-grid">
		<div class="field">
			<div class="field-label-row"><label>آدرس تصویر</label>`))
		writeGuideLink(w, "branding", "آدرس تصویر")
		w.Write([]byte(`</div>
			<input type="url" name="banner_url" value="`))
		pages.WriteEscaped(w, banner.URL)
		w.Write([]byte(`" placeholder="https://fotomi.net/banner.jpg" dir="ltr" />
		</div>
		<div class="field">
			<div class="field-label-row"><label>متن بنر</label>`))
		writeGuideLink(w, "branding", "متن بنر")
		w.Write([]byte(`</div>
			<input type="text" name="banner_text" value="`))
		pages.WriteEscaped(w, banner.Text)
		w.Write([]byte(`" placeholder="به فوتومی خوش آمدید" />
		</div>
		<div class="field">
			<div class="field-label-row"><label>لینک (اختیاری)</label>`))
		writeGuideLink(w, "branding", "لینک بنر")
		w.Write([]byte(`</div>
			<input type="url" name="banner_link" value="`))
		pages.WriteEscaped(w, banner.Link)
		w.Write([]byte(`" placeholder="https://fotomi.net/offers" dir="ltr" />
		</div>
		<div class="field">
			<div class="field-label-row"><label>مدت نمایش</label>`))
		writeGuideLink(w, "branding", "مدت نمایش بنر")
		w.Write([]byte(`</div>
			<div class="input-suffix">
				<input type="number" name="banner_duration" value="`))
		pages.WriteEscaped(w, itostr(banner.Duration))
		w.Write([]byte(`" min="1" />
				<span>ثانیه</span>
			</div>
		</div>
		<button type="button" class="icon-action icon-action-danger banner-remove" data-banner-remove title="حذف بنر" aria-label="حذف بنر">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M5 7h14M10 7V5h4v2M7 7l.7 12h8.6L17 7"/><path d="M10.5 10.5v5M13.5 10.5v5"/></svg>
		</button>
	</div>
</div>`))
		return nil
	})
}

func writeGuideLink(w io.Writer, anchor, label string) {
	w.Write([]byte(`<a class="field-guide-link" href="/admin/guide#`))
	pages.WriteEscaped(w, anchor)
	w.Write([]byte(`" target="_blank" rel="noopener noreferrer" aria-label="راهنمای `))
	pages.WriteEscaped(w, label)
	w.Write([]byte(`">راهنما ↗</a>`))
}

func writeOption(w io.Writer, value, label, selected string) {
	w.Write([]byte(`<option value="`))
	pages.WriteEscaped(w, value)
	w.Write([]byte(`"`))
	if value == selected {
		w.Write([]byte(` selected`))
	}
	w.Write([]byte(`>`))
	pages.WriteEscaped(w, label)
	w.Write([]byte(`</option>`))
}

// AppFormData is everything the create/edit form needs. The handler builds it
// from the domain app (edit) or safe defaults (create).
type AppFormData struct {
	ID          string
	Name        string
	Url         string
	CallbackURL string
	WebhookURL  string

	OTPChannel string
	OTPLength  int
	OTPCharset string
	OTPExpire  int
	OTTExpire  int

	SessionTimeout        int
	InfiniteSession       bool
	SessionCookieName     string
	SessionCookieHTTPOnly bool
	SessionCookieSecure   bool
	SessionCookieSameSite string
	RedirectBack          bool

	CustomAuthEnabled           bool
	AllowedCountries            string
	AuthPageEnabled             bool
	AuthPageURL                 string
	AuthPageTitle               string
	AuthPageLayout              string
	AuthPageFont                string
	AuthPageTextAlign           string
	AuthPageLogo                string
	AuthPageLogoMask            string
	AuthPageLogoPosition        string
	AuthPagePrimaryColor        string
	AuthPageBackgroundColor     string
	AuthPageCardBackgroundColor string
	AuthPageTextColor           string
	AuthPageMutedTextColor      string
	AuthPagePinBackgroundColor  string
	AuthPagePinTextColor        string
	AuthPageTheme               string
	AuthPageCaptchaEnabled      bool
	AuthPageCaptchaProvider     string
	AuthPageCaptchaSiteKey      string
	AuthPageCaptchaSecretKey    string
	AuthPageBannerMode          string
	AuthPageBanners             []app.Banner
	AuthPageCopyCustom          bool

	AuthPageLoginTitle      string
	AuthPageLoginSubtitle   string
	AuthPageIdentifierLabel string
	AuthPageSendButton      string
	AuthPageVerifyTitle     string
	AuthPageVerifySubtitle  string
	AuthPageVerifyButton    string
	AuthPageRetryText       string
	AuthPageBackText        string
	AuthPageFooterText      string

	Error   string
	Success string
}
