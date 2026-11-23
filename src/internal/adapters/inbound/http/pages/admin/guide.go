package pages

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// Guide renders the written documentation for the admin panel. It is a static
// page: every section explains a real setting from the app form so an operator
// can finish the integration without leaving the panel.
func Guide() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		w.Write([]byte(`<section class="page-head">
	<div>
		<h2>راهنمای امان</h2>
		<p>از ساخت برنامه تا برندینگ صفحه ورود و امنیت نشست، در چند دقیقه.</p>
	</div>
	<a href="/admin/apps/new" class="btn btn-primary">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M12 5v14M5 12h14"/></svg>
		برنامه جدید
	</a>
</section>

<div class="docs-layout">
	<aside class="card docs-nav" aria-label="فهرست راهنما">
		<span class="docs-nav-title">در این صفحه</span>
		<a href="#quick-start">شروع سریع</a>
		<a href="#app-basics">اطلاعات برنامه</a>
		<a href="#otp">کد یکبار مصرف</a>
		<a href="#sessions">نشست و کوکی</a>
		<a href="#auth-page">صفحه ورود</a>
		<a href="#branding">برندینگ و بنرها</a>
		<a href="#captcha">کپچای ArCaptcha</a>
		<a href="#webhook">وب‌هوک</a>
		<a href="#checklist">چک‌لیست انتشار</a>
	</aside>

	<div class="docs-body">

		<section class="card docs-section" id="quick-start">
			<div class="docs-section-head">
				<span class="docs-step">۱</span>
				<div>
					<h3>شروع سریع</h3>
					<p>سه قدم تا اولین ورود موفق.</p>
				</div>
			</div>
			<ol class="docs-list">
				<li>از صفحه <a href="/admin/apps/new">برنامه جدید</a> یک برنامه بسازید. نام، آدرس اصلی و آدرس بازگشت (callback) را وارد کنید.</li>
				<li>پس از ذخیره، شناسه برنامه در آدرس صفحه ویرایش و در ستون شناسه فهرست برنامه‌ها دیده می‌شود. همین مقدار را به‌عنوان <code>app_id</code> نگه دارید.</li>
				<li>کاربر را برای ورود به این آدرس بفرستید:
					<div class="code-block">
						<div class="code-head"><span>آدرس شروع ورود</span><button type="button" class="code-copy" data-copy>کپی</button></div>
						<pre><code>https://cas.example.com/auth/check?app_id=APP_ID</code></pre>
					</div>
				</li>
			</ol>
			<div class="callout callout-info">
				<strong>نکته</strong>
				<span>اگر برنامه‌ی شما صفحه ورود اختصاصی دارد، گزینه «استفاده از صفحه ورود میزبانی‌شده» را روشن کنید تا کاربر به‌جای صفحه پیش‌فرض CAS به آدرس شما فرستاده شود.</span>
			</div>
		</section>

		<section class="card docs-section" id="app-basics">
			<div class="docs-section-head">
				<span class="docs-step">۲</span>
				<div>
					<h3>اطلاعات برنامه</h3>
					<p>هویت برنامه و مقصد بازگشت کاربر.</p>
				</div>
			</div>
			<dl class="docs-defs">
				<div><dt>نام برنامه</dt><dd>در صفحه ورود، فهرست برنامه‌ها و جزئیات نشست‌ها نمایش داده می‌شود. کوتاه و قابل تشخیص باشد.</dd></div>
				<div><dt>آدرس اصلی</dt><dd>مبنای اعتبارسنجی بازگشت خودکار است. فقط آدرس‌هایی با همین دامنه پذیرفته می‌شوند تا از بازگشت به مقصد ناشناس جلوگیری شود.</dd></div>
				<div><dt>آدرس بازگشت (callback)</dt><dd>پس از تایید کد، کاربر همراه یک تیکت یکبار مصرف (OTT) به این آدرس برمی‌گردد. این تیکت را در بک‌اند خود با <code>POST /v1/auth/verify-ott</code> مبادله کنید.</dd></div>
				<div><dt>آدرس وب‌هوک</dt><dd>اختیاری. رویدادهای ارسال کد، ورود و خروج به‌صورت JSON به این آدرس ارسال می‌شوند.</dd></div>
			</dl>
		</section>

		<section class="card docs-section" id="otp">
			<div class="docs-section-head">
				<span class="docs-step">۳</span>
				<div>
					<h3>کد یکبار مصرف</h3>
					<p>کانال ارسال و طول عمر کد.</p>
				</div>
			</div>
			<dl class="docs-defs">
				<div><dt>روش ارسال کد</dt><dd>پیامک برای شماره موبایل، ایمیل برای آدرس ایمیل، یا هر دو. در حالت «هر دو» نوع شناسه‌ای که کاربر وارد می‌کند تعیین می‌کند کد از کدام کانال برود.</dd></div>
				<div><dt>طول و ترکیب کد</dt><dd>بین ۴ تا ۳۲ کاراکتر. کدهای عددی برای صفحه‌کلید موبایل ساده‌تر هستند؛ کدهای حرف و عدد امنیت بیشتری دارند.</dd></div>
				<div><dt>اعتبار کد</dt><dd>مدت زمان قابل استفاده بودن کد. پس از پایان این زمان، دکمه «ارسال دوباره» در صفحه تایید فعال می‌شود.</dd></div>
				<div><dt>اعتبار تیکت یکبار مصرف</dt><dd>فرصت برنامه شما برای مبادله OTT پس از بازگشت کاربر. کوتاه نگه دارید؛ تیکت فقط یک‌بار قابل استفاده است.</dd></div>
			</dl>
		</section>

		<section class="card docs-section" id="sessions">
			<div class="docs-section-head">
				<span class="docs-step">۴</span>
				<div>
					<h3>نشست و کوکی</h3>
					<p>مدت ورود و رفتار کوکی روی مرورگر کاربر.</p>
				</div>
			</div>
			<dl class="docs-defs">
				<div><dt>مدت نشست</dt><dd>پس از این مدت، نشست منقضی می‌شود و کاربر باید دوباره کد بگیرد. با فعال کردن «نشست همیشگی» این مقدار نادیده گرفته می‌شود.</dd></div>
				<div><dt>نام کوکی نشست</dt><dd>اگر خالی بماند <code>cas_cookie</code> استفاده می‌شود. برای چند برنامه روی یک دامنه، نام‌های متفاوت انتخاب کنید تا نشست‌ها روی هم نیفتند.</dd></div>
				<div><dt>SameSite</dt><dd>برای ورود از دامنه‌های مختلف معمولاً <code>None</code> همراه با Secure لازم است. در حالت تک‌دامنه، <code>Lax</code> امن و کافی است.</dd></div>
				<div><dt>HttpOnly و Secure</dt><dd>HttpOnly دسترسی جاوااسکریپت به کوکی را می‌بندد. Secure کوکی را فقط روی HTTPS ارسال می‌کند؛ در محیط توسعه بدون HTTPS خاموش بماند.</dd></div>
				<div><dt>بازگشت به آخرین آدرس</dt><dd>اگر روشن باشد و آدرس مبدأ متعلق به دامنه برنامه باشد، کاربر پس از ورود به همان صفحه برمی‌گردد، نه فقط به callback.</dd></div>
			</dl>
		</section>

		<section class="card docs-section" id="auth-page">
			<div class="docs-section-head">
				<span class="docs-step">۵</span>
				<div>
					<h3>صفحه ورود: پیش‌فرض یا میزبانی‌شده</h3>
					<p>انتخاب کنید CAS صفحه ورود را نشان دهد یا برنامه شما.</p>
				</div>
			</div>
			<p>به‌صورت پیش‌فرض، CAS صفحه ورود و تایید کد را با برند برنامه شما نمایش می‌دهد. اگر صفحه ورود اختصاصی دارید، گزینه «استفاده از صفحه ورود میزبانی‌شده در برنامه» را روشن کنید و آدرس آن را وارد کنید. CAS پارامترهای لازم را به آدرس شما اضافه می‌کند و تنظیمات ظاهری صفحه پیش‌فرض پنهان می‌شوند.</p>
			<div class="callout callout-warn">
				<strong>توجه</strong>
				<span>در حالت میزبانی‌شده، مسئولیت نمایش فرم، اعتبارسنجی شناسه و نمایش خطاها با برنامه شماست. مسیرهای JSON در صفحه توسعه‌دهندگان آمده‌اند.</span>
			</div>
		</section>

		<section class="card docs-section" id="branding">
			<div class="docs-section-head">
				<span class="docs-step">۶</span>
				<div>
					<h3>برندینگ و بنرها</h3>
					<p>ظاهر صفحه پیش‌فرض را با هویت برنامه هماهنگ کنید.</p>
				</div>
			</div>
			<dl class="docs-defs">
				<div><dt>چیدمان</dt><dd>کارت مرکزی، دو ستونه یا مینیمال. در چیدمان دو ستونه، بنرها در ستون کنار فرم نمایش داده می‌شوند.</dd></div>
				<div><dt>فونت و تراز متن</dt><dd>وزیرمتن، ایران‌سنس یا فونت سیستم. تراز متن روی عنوان‌ها و پانوشت اثر می‌گذارد.</dd></div>
				<div><dt>رنگ‌ها</dt><dd>رنگ اصلی برای دکمه‌ها و تاکیدها، رنگ پس‌زمینه صفحه و رنگ کارت ورود جداگانه تنظیم می‌شوند.</dd></div>
				<div><dt>لوگو</dt><dd>آدرس تصویر، شکل ماسک (دایره، گوشه‌گرد، مربع) و جای آن. اگر آدرس خالی باشد حرف اول نام برنامه نمایش داده می‌شود.</dd></div>
				<div><dt>بنرها</dt><dd>هر بنر تصویر، متن، لینک و مدت نمایش خودش را دارد. ترتیب با کشیدن دسته هر کارت تغییر می‌کند. نمایش زمان بنر را می‌توان به‌صورت نوار پیشرفت، ثانیه شمار یا بدون نمایش انتخاب کرد.</dd></div>
			</dl>
			<p class="docs-note">پیش‌نمایش کنار فرم با هر تغییر به‌روزرسانی می‌شود؛ نیازی به ذخیره برای دیدن نتیجه نیست.</p>
		</section>

		<section class="card docs-section" id="captcha">
			<div class="docs-section-head">
				<span class="docs-step">۷</span>
				<div>
					<h3>کپچای ArCaptcha</h3>
					<p>جلوگیری از درخواست خودکار برای ارسال کد.</p>
				</div>
			</div>
			<p>کپچا برای هر برنامه جداگانه تنظیم می‌شود. در فرم برنامه گزینه «کپچای ArCaptcha» را روشن کنید تا دو فیلد کلید نمایش داده شوند:</p>
			<dl class="docs-defs">
				<div><dt>ArCaptcha Site Key</dt><dd>کلید عمومی که در ویجت کپچای صفحه ورود استفاده می‌شود.</dd></div>
				<div><dt>ArCaptcha Secret Key</dt><dd>کلید خصوصی که سرور برای تایید توکن کپچا به ArCaptcha ارسال می‌کند. این مقدار فقط در پنل ذخیره می‌شود و در صفحه ورود نمایش داده نمی‌شود.</dd></div>
			</dl>
			<p>اگر یکی از کلیدها خالی باشد، ویجت نمایش داده نمی‌شود و ورود همچنان کار می‌کند؛ این کار از قفل شدن کاربران در محیط توسعه جلوگیری می‌کند. تایید کپچا فقط هنگام ارسال اولیه کد انجام می‌شود و ارسال دوباره پس از پایان شمارش معکوس از آن مستثنی است.</p>
		</section>

		<section class="card docs-section" id="webhook">
			<div class="docs-section-head">
				<span class="docs-step">۸</span>
				<div>
					<h3>وب‌هوک</h3>
					<p>رویدادها را در سیستم خودتان دریافت کنید.</p>
				</div>
			</div>
			<p>پس از هر رویداد، یک درخواست <code>POST</code> با بدنه JSON به آدرس وب‌هوک برنامه ارسال می‌شود. رویدادهای ارسال کد، تایید کد، ورود، خروج و ابطال نشست پوشش داده می‌شوند. پاسخ ۲xx را سریع برگردانید؛ عملیات سنگین را در صف خودتان انجام دهید.</p>
		</section>

		<section class="card docs-section" id="checklist">
			<div class="docs-section-head">
				<span class="docs-step">۹</span>
				<div>
					<h3>چک‌لیست انتشار</h3>
					<p>پیش از تولید، این موارد را بررسی کنید.</p>
				</div>
			</div>
			<ul class="docs-checks">
				<li>آدرس اصلی و callback با دامنه واقعی برنامه یکی است.</li>
				<li>کوکی Secure روی HTTPS روشن و در محیط توسعه خاموش است.</li>
				<li>نام کوکی بین برنامه‌های هم‌دامنه یکتاست.</li>
				<li>کلیدهای ArCaptcha در فرم برنامه وارد شده‌اند.</li>
				<li>وب‌هوک به یک آدرس پایدار و پاسخ‌دهنده اشاره می‌کند.</li>
				<li>صفحه ورود در موبایل و دسکتاپ بررسی شده است.</li>
			</ul>
		</section>

	</div>
</div>`))
		return nil
	})
}
