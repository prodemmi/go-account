package pages

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// Developers renders the API reference for app developers. It combines the
// browser flow used by the branded auth page with the JSON endpoints used by
// apps that host their own login form, in a Swagger-like layout.
func Developers() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		w.Write([]byte(`<section class="page-head">	<div>
		<h2>توسعه‌دهندگان</h2>
	</div>
	<span class="badge badge-info">نسخه ۱.۰</span>
</section>

<div class="docs-layout">
	<aside class="card docs-nav" aria-label="فهرست مستندات">
		<span class="docs-nav-title">در این صفحه</span>
		<a href="#basics">مبانی</a>
		<a href="#auth-modes">حالت‌های احراز هویت</a>
		<a href="#browser-flow">جریان مرورگر</a>
		<a href="#browser-api">مسیرهای صفحه ورود</a>
		<a href="#json-api">API برنامه‌ها</a>
		<a href="#cookies">کوکی و نشست</a>
		<a href="#verify-ott">مبادله تیکت</a>
		<a href="#user-sessions">نشست‌های کاربر</a>
		<a href="#hosted-page">صفحه ورود اختصاصی</a>
		<a href="#errors">خطاها</a>
	</aside>

	<div class="docs-body">

		<section class="card docs-section" id="basics">
			<h3>مبانی</h3>
			<p>آدرس پایه سرویس را از دامنه استقرار خود بگیرید. در نمونه‌ها <code>https://cas.example.com</code> استفاده شده است.</p>
			<dl class="docs-defs">
				<div><dt>شناسه برنامه</dt><dd>همه درخواست‌های API با هدر <code>APP-ID</code> شناسایی می‌شوند. این مقدار همان شناسه‌ای است که در پنل مدیریت می‌بینید.</dd></div>
				<div><dt>نشست</dt><dd>پس از تایید کد، کوکی نشست روی مرورگر تنظیم می‌شود. در درخواست‌های سمت سرور می‌توانید توکن را در هدر <code>API-TOKEN</code> بفرستید.</dd></div>
				<div><dt>پاسخ‌ها</dt><dd>پاسخ‌های موفق JSON با ساختار <code>{"status":"success","data":{...}}</code> هستند.</dd></div>
				<div><dt>خطاها</dt><dd>درخواست نامعتبر یا برنامه ناشناس به صفحه ۴۰۴ هدایت می‌شود. خطاهای منطقی جریان مرورگر با پارامتر <code>error</code> و پیام فارسی برمی‌گردند.</dd></div>
			</dl>
		</section>

		<section class="card docs-section" id="auth-modes">
			<h3>حالت‌های احراز هویت</h3>
			<p>حالت هر برنامه در تنظیمات آن مشخص می‌شود. <code>custom_auth_enabled</code> به‌طور پیش‌فرض خاموش است.</p>
			<h4>مدیریت‌شده (Managed)</h4>
			<p>با <code>custom_auth_enabled=false</code>، Aman مسیر مرورگری <code>/auth/check</code> و صفحه‌های ورود و تایید کد را مدیریت می‌کند. تنظیم «صفحه ورود میزبانی‌شده» فقط مسیر ورود این حالت را به URL ثبت‌شده برنامه هدایت می‌کند؛ با حالت API سفارشی یکی نیست.</p>
			<h4>سفارشی (Custom)</h4>
			<p>با <code>custom_auth_enabled=true</code>، درخواست‌های این برنامه به همه مسیرهای <code>/auth/**</code> با HTTP 404 پاسخ می‌گیرند؛ مسیرهای بدون شناسه برنامه نیز ۴۰۴ هستند. رابط کسب‌وکار باید خودش کد ورود را نمایش دهد و بک‌اند کسب‌وکار باید APIهای زیر را با هدر <code>APP-ID</code> فراخوانی کند. مسیرهای <code>/v1/auth/*</code> فعال می‌مانند.</p>
			<ol class="docs-list">
				<li><code>POST /v1/auth/send-otp</code> با <code>{"identifier":"…"}</code> کد را می‌فرستد و <code>otp_hash</code> رمزنگاری‌شده و احرازاصالت‌شده با <code>APP_SECRET</code> را برمی‌گرداند. اگر کپچا روشن است، توکن ویجت انتخاب‌شده را در <code>captcha_token</code> بفرستید. توکن OTP، برنامه، شناسه، User-Agent و انقضا را به هم متصل می‌کند؛ رکورد OTP سمت سرور همچنان مرجع یک‌بارمصرف بودن کد است.</li>
				<li><code>POST /v1/auth/verify-otp</code> با <code>identifier</code>، <code>otp_code</code> و همان <code>otp_hash</code> کد را تایید می‌کند و نشست معلق، OTT یک‌بارمصرف و کوکی Aman را صادر می‌کند.</li>
				<li>بک‌اند برنامه، <code>ott</code> را با <code>POST /v1/auth/verify-ott</code> مبادله می‌کند. پاسخ شامل توکن نشست و هویت تاییدشده است؛ برنامه باید نشست و کوکی مخصوص دامنه خودش را در سمت سرور ایجاد کند.</li>
			</ol>
			<p>توکن نشست را در کد مرورگر یا URL قرار ندهید. بک‌اند باید <code>User-Agent</code> کاربر را در درخواست‌های ارسال، تایید و بررسی نشست یکسان نگه دارد؛ توکن به همان مرورگر مقید است. این سرور CORS را فعال نمی‌کند؛ برای رابطی روی دامنه دیگر، API را از بک‌اند خود فراخوانی کنید، نه مستقیماً از JavaScript مرورگر.</p>
		</section>

		<section class="card docs-section" id="browser-flow">
			<h3>جریان مرورگر</h3>
			<p>برای برنامه‌هایی که از صفحه ورود پیش‌فرض CAS استفاده می‌کنند، کاربر را به مسیر <code>/auth/check</code> بفرستید. اگر نشست معتبر باشد، کاربر با یک OTT تازه به برنامه برمی‌گردد؛ در غیر این صورت صفحه ورود برند برنامه نمایش داده می‌شود.</p>
			<div class="flow">
				<span class="flow-node">برنامه</span>
				<span class="flow-arrow">←</span>
				<span class="flow-node">/auth/check</span>
				<span class="flow-arrow">←</span>
				<span class="flow-node">صفحه ورود</span>
				<span class="flow-arrow">←</span>
				<span class="flow-node">تایید کد</span>
				<span class="flow-arrow">←</span>
				<span class="flow-node">callback + OTT</span>
			</div>
			<div class="code-block">
				<div class="code-head"><span>شروع ورود</span><button type="button" class="code-copy" data-copy>کپی</button></div>
				<pre><code>GET /auth/check?app_id=APP_ID&amp;redirect=https://app.example.com/dashboard</code></pre>
			</div>
			<p>پارامتر <code>redirect</code> اختیاری است و فقط وقتی پذیرفته می‌شود که با دامنه «آدرس اصلی» برنامه یکی باشد.</p>
		</section>

		<section class="card docs-section" id="browser-api">
			<h3>مسیرهای صفحه ورود</h3>
			<p>این مسیرها فرم‌های HTML را پردازش می‌کنند. صفحه ورود به‌صورت SPA کار می‌کند: اگر درخواست با هدر <code>X-Aman-SPA: 1</code> یا <code>Accept: application/json</code> ارسال شود، پاسخ JSON است و بدون بارگذاری دوباره صفحه بین مراحل جابه‌جا می‌شوید. بدون جاوااسکریپت همان فرم‌ها با ریدایرکت کار می‌کنند.</p>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-get">GET</span>
					<code class="endpoint-path">/auth/login</code>
				</div>
				<p>نمایش صفحه ورود برند برنامه. پارامترها: <code>app_id</code> (الزامی)، <code>redirect</code> (اختیاری).</p>
			</article>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-post">POST</span>
					<code class="endpoint-path">/auth/send-otp</code>
				</div>
				<p>ارسال کد یکبار مصرف. بدنه <code>application/x-www-form-urlencoded</code>:</p>
				<table class="params-table">
					<thead><tr><th>پارامتر</th><th>نوع</th><th>توضیح</th></tr></thead>
					<tbody>
						<tr><td><code>app_id</code></td><td>string</td><td>شناسه برنامه، الزامی.</td></tr>
						<tr><td><code>identifier</code></td><td>string</td><td>ایمیل یا شماره موبایل، الزامی.</td></tr>
						<tr><td><code>redirect</code></td><td>string</td><td>آدرس بازگشت اختیاری روی دامنه برنامه.</td></tr>
						<tr><td><code>arcaptcha-token</code> / <code>g-recaptcha-response</code></td><td>string</td><td>توکن ارائه‌دهنده انتخاب‌شده، فقط وقتی کپچا برای برنامه روشن است.</td></tr>
						<tr><td><code>retry</code></td><td>string</td><td>مقدار <code>1</code> برای ارسال دوباره از صفحه تایید؛ کپچا همچنان الزامی است.</td></tr>
					</tbody>
				</table>
				<p>پاسخ SPA:</p>
				<div class="code-block">
					<div class="code-head"><span>JSON</span><button type="button" class="code-copy" data-copy>کپی</button></div>
					<pre><code>{
  "ok": true,
  "step": "verify",
  "identifier": "09123456789",
  "otp_hash": "…",
  "length": 5,
  "expires_at": 1730000000,
  "redirect": ""
}</code></pre>
				</div>
			</article>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-post">POST</span>
					<code class="endpoint-path">/auth/verify-otp</code>
				</div>
				<p>تایید کد و صدور نشست. پارامترها: <code>app_id</code>، <code>identifier</code>، <code>otp_hash</code>، <code>otp_code</code>، <code>expires_at</code> و <code>redirect</code>. مقدار URL برای <code>expires_at</code> فقط شمارش‌گر رابط است؛ اعتبار امضا و انقضای رکورد سمت سرور بررسی می‌شود. در حالت SPA پاسخ شامل <code>{"ok":true,"redirect":"…"}</code> است و کوکی نشست تنظیم می‌شود.</p>
			</article>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-get">GET</span>
					<span class="method method-post">POST</span>
					<code class="endpoint-path">/auth/logout</code>
				</div>
				<p>ابطال نشست جاری و بازگشت به آدرس اصلی برنامه. پارامتر <code>app_id</code> الزامی است.</p>
			</article>
		</section>

		<section class="card docs-section" id="json-api">
			<h3>API برنامه‌ها</h3>
			<p>اگر فرم ورود را در برنامه خود می‌سازید، از مسیرهای زیر استفاده کنید. همه آن‌ها هدر <code>APP-ID</code> می‌خواهند و بدنه JSON می‌گیرند.</p>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-post">POST</span>
					<code class="endpoint-path">/v1/auth/send-otp</code>
				</div>
				<div class="code-block">
					<div class="code-head"><span>درخواست</span><button type="button" class="code-copy" data-copy>کپی</button></div>
					<pre><code>curl -X POST https://cas.example.com/v1/auth/send-otp \
  -H "APP-ID: APP_ID" \
  -H "Content-Type: application/json" \
  -d '{"identifier":"09123456789","captcha_token":"…"}'</code></pre>
				</div>
				<p><code>captcha_token</code> فقط وقتی تنظیم کپچای برنامه روشن است لازم است؛ مقدار آن را از فیلد پاسخ ارائه‌دهنده انتخاب‌شده بگیرید.</p>
				<div class="code-block">
					<div class="code-head"><span>پاسخ</span><button type="button" class="code-copy" data-copy>کپی</button></div>
					<pre><code>{
  "status": "success",
  "data": {
    "app_id": "…",
    "identifier": "09123456789",
    "channel": "sms",
    "otp_hash": "…",
    "length": 5,
    "expires_at": "2026-01-01T12:02:00Z"
  }
}</code></pre>
				</div>
			</article>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-post">POST</span>
					<code class="endpoint-path">/v1/auth/verify-otp</code>
				</div>
				<div class="code-block">
					<div class="code-head"><span>درخواست</span><button type="button" class="code-copy" data-copy>کپی</button></div>
					<pre><code>curl -X POST https://cas.example.com/v1/auth/verify-otp \
  -H "APP-ID: APP_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "09123456789",
    "otp_code": "12345",
    "otp_hash": "…"
  }'</code></pre>
				</div>
				<p>در پاسخ، <code>ott</code> تیکت یکبار مصرف و <code>redirect_url</code> آدرس بازگشت است. کوکی نشست هم روی پاسخ تنظیم می‌شود.</p>
			</article>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-get">GET</span>
					<code class="endpoint-path">/v1/auth/check</code>
				</div>
				<p>مسیر سازگاری قدیمی: در حالت مدیریت‌شده، نشست را بررسی و مرورگر را با OTT تازه به callback هدایت می‌کند یا به صفحه ورود می‌فرستد. در حالت سفارشی، نتیجه JSON برمی‌گرداند تا به صفحه غیرفعال‌شده <code>/auth/login</code> هدایت نشود.</p>
			</article>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-get">GET</span>
					<code class="endpoint-path">/v1/auth/check-login</code>
				</div>
				<p>بررسی JSON نشست برای بک‌اند برنامه. هدرهای <code>APP-ID</code> و یکی از کوکی نشست یا <code>API-TOKEN</code> را بفرستید. پارامتر <code>redirect</code> اختیاری است و فقط URL مجاز همان برنامه پذیرفته می‌شود.</p>
				<div class="code-block">
					<div class="code-head"><span>درخواست</span><button type="button" class="code-copy" data-copy>کپی</button></div>
					<pre><code>curl -i "https://cas.example.com/v1/auth/check-login?redirect=https://app.example.com/dashboard" \
  -H "APP-ID: APP_ID" \
  -H "API-TOKEN: SESSION_TOKEN"</code></pre>
				</div>
				<div class="code-block">
					<div class="code-head"><span>پاسخ در صورت ورود</span><button type="button" class="code-copy" data-copy>کپی</button></div>
					<pre><code>{
  "status": "success",
  "data": {
    "authenticated": true,
    "ott": "ONE_TIME_TICKET",
    "redirect_url": "https://app.example.com/callback?ott=ONE_TIME_TICKET"
  }
}</code></pre>
				</div>
				<p>بدون نشست معتبر، <code>data.authenticated</code> برابر <code>false</code> است. در صورت معتبر بودن نشست، هر فراخوانی یک OTT تازه صادر می‌کند.</p>
			</article>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-get">GET</span>
					<code class="endpoint-path">/v1/auth/me</code>
				</div>
				<p>اطلاعات نشست جاری. هدر <code>API-TOKEN</code> یا کوکی نشست لازم است.</p>
			</article>

			<article class="endpoint">
				<div class="endpoint-head">
					<span class="method method-delete">DELETE</span>
					<code class="endpoint-path">/v1/auth/logout</code>
				</div>
				<p>ابطال نشست جاری. معادل <code>POST</code> هم پشتیبانی می‌شود.</p>
			</article>
		</section>

		<section class="card docs-section" id="verify-ott">
			<h3>مبادله تیکت یکبار مصرف</h3>
			<p>پس از بازگشت کاربر به callback، پارامتر <code>ott</code> را بگیرید و آن را در بک‌اند خود مبادله کنید. هر تیکت فقط یک‌بار معتبر است.</p>
			<div class="code-block">
				<div class="code-head"><span>درخواست</span><button type="button" class="code-copy" data-copy>کپی</button></div>
				<pre><code>curl -X POST https://cas.example.com/v1/auth/verify-ott \
  -H "APP-ID: APP_ID" \
  -H "Content-Type: application/json" \
  -d '{"ott":"OTT_VALUE"}'</code></pre>
			</div>
			<div class="code-block">
				<div class="code-head"><span>پاسخ</span><button type="button" class="code-copy" data-copy>کپی</button></div>
				<pre><code>{
  "status": "success",
  "data": {
    "session_token": "…",
    "identifier": "09123456789",
    "expires_at": "2026-01-01T13:00:00Z",
    "infinite": false
  }
}</code></pre>
			</div>
		</section>

		<section class="card docs-section" id="user-sessions">
			<h3>نشست‌های کاربر</h3>
			<p>کاربر می‌تواند نشست‌های خودش را ببیند و خاتمه دهد. این مسیرها هم به <code>APP-ID</code> و توکن نشست نیاز دارند.</p>
			<article class="endpoint">
				<div class="endpoint-head"><span class="method method-get">GET</span><code class="endpoint-path">/v1/user/sessions</code></div>
				<p>فهرست نشست‌های کاربر جاری.</p>
			</article>
			<article class="endpoint">
				<div class="endpoint-head"><span class="method method-delete">DELETE</span><code class="endpoint-path">/v1/user/sessions/:id</code></div>
				<p>خاتمه یک نشست مشخص از نشست‌های کاربر.</p>
			</article>
			<article class="endpoint">
				<div class="endpoint-head"><span class="method method-delete">DELETE</span><code class="endpoint-path">/v1/user/sessions</code></div>
				<p>خاتمه همه نشست‌های کاربر.</p>
			</article>
		</section>

		<section class="card docs-section" id="hosted-page">
			<h3>صفحه ورود میزبانی‌شده</h3>
			<p>تنظیم قدیمی «صفحه ورود میزبانی‌شده» در حالت Managed، <code>/auth/login</code> را به URL ثبت‌شده با پارامتر <code>app_id</code> می‌فرستد. برای ساخت یک جریان کاملاً سفارشی، به‌جای آن <code>custom_auth_enabled</code> را روشن کنید و API را از بک‌اند برنامه فراخوانی کنید.</p>
		</section>

		<section class="card docs-section" id="cookies">
			<h3>کوکی و نشست</h3>
			<ul class="docs-list">
				<li>نام کوکی از تنظیم برنامه می‌آید و مقدار پیش‌فرض آن <code>cas_cookie</code> است. APIهای بررسی نشست توکن را از همین کوکی می‌خوانند؛ سپس <code>cas_cookie</code> و <code>cas_session</code> را نیز برای سازگاری قدیمی امتحان می‌کنند.</li>
				<li>کوکی روی پاسخ تایید OTP و مبادله OTT تنظیم می‌شود: مسیر <code>/</code>، انقضای نشست، و گزینه‌های <code>HttpOnly</code>، <code>Secure</code> و <code>SameSite</code> مطابق تنظیمات برنامه. نام پیش‌فرض <code>cas_cookie</code> و SameSite پیش‌فرض <code>Lax</code> است؛ درخواست HTTPS نیز کوکی را <code>Secure</code> می‌کند. برای <code>SameSite=None</code> از HTTPS استفاده کنید. عمر Max-Age نشست همیشگی حداکثر ۱۰ سال است.</li>
				<li>اگر دامنه Aman و callback برنامه زیر یک دامنه والد مشترک باشند، کوکی برای آن والد تنظیم می‌شود. در دامنه‌های نامرتبط، مرورگر کوکی صادرشده از Aman را برای دامنه برنامه نمی‌پذیرد؛ تبادل OTT و کوکی نشست خود برنامه راه ارتباطی است.</li>
				<li><code>GET /v1/auth/check-login</code> کوکی را می‌پذیرد؛ برای فراخوانی سمت سرور می‌توانید همان توکن نشست را در هدر <code>API-TOKEN</code> بفرستید. <code>User-Agent</code> اصلی کاربر را نیز حفظ کنید. توکن را محرمانه و فقط سمت سرور نگه دارید.</li>
			</ul>
		</section>

		<section class="card docs-section" id="errors">
			<h3>خطاها و محدودیت‌ها</h3>
			<ul class="docs-list">
				<li>برنامه ناشناس، درخواست ناقص یا متد نامعتبر به صفحه ۴۰۴ هدایت می‌شود.</li>
				<li>کد منقضی یا نادرست پیام «کد وارد شده معتبر نیست یا منقضی شده است» برمی‌گرداند.</li>
				<li>در حالت SPA خطاها با ساختار <code>{"ok":false,"error":"…"}</code> و کد HTTP 200 برمی‌گردند تا پیام فارسی در همان صفحه نمایش داده شود.</li>
				<li>توکن ArCaptcha یک‌بار مصرف است؛ برای هر ارسال تازه باید کپچا دوباره تایید شود.</li>
				<li>تیکت OTT کوتاه‌عمر و یک‌بار مصرف است؛ آن را فقط در بک‌اند مبادله کنید.</li>
			</ul>
		</section>

	</div>
</div>`))
		return nil
	})
}
