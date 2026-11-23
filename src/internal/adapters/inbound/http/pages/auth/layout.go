package pages

import (
	"context"
	"io"
	"strconv"
	"strings"

	"aman/internal/adapters/inbound/http/pages"

	"github.com/a-h/templ"
)

// Layout is the base HTML layout for the CAS auth pages. Colors, theme, font,
// alignment, layout variant and banners come from the app configuration.
func Layout(data pages.AuthData, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		theme := data.Theme
		if theme == "" {
			theme = "light"
		}
		primary := data.PrimaryColor
		if primary == "" {
			primary = "#3b82f6"
		}
		background := data.BackgroundColor
		if background == "" {
			background = "#f8fafc"
		}
		pinBackground := data.PinBackgroundColor
		pinText := data.PinTextColor
		font := data.Font
		if font == "" {
			font = "vazirmatn"
		}
		align := data.TextAlign
		if align == "" {
			align = "center"
		}
		layout := data.Layout
		if layout == "" {
			layout = "card"
		}
		bannerMode := data.BannerMode
		if bannerMode == "" {
			bannerMode = "none"
		}
		documentTitle := data.DocumentTitle
		if documentTitle == "" {
			documentTitle = data.PageTitle
		}

		w.Write([]byte(`<!DOCTYPE html>
<html lang="fa" dir="rtl" data-theme="`))
		pages.WriteEscaped(w, theme)
		w.Write([]byte(`" data-font="`))
		pages.WriteEscaped(w, font)
		w.Write([]byte(`">
<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	<title>`))
		pages.WriteEscaped(w, documentTitle)
		w.Write([]byte(`</title>
	<link rel="stylesheet" href="/static/css/auth.css?v=10" />`))
		if data.CaptchaEnabled {
			w.Write([]byte(`
	<script src="https://nwidget.arcaptcha.ir/1/api.js" async defer></script>`))
		}
		w.Write([]byte(`
</head>
<body class="layout-`))
		pages.WriteEscaped(w, layout)
		w.Write([]byte(`" data-banner-mode="`))
		pages.WriteEscaped(w, bannerMode)
		w.Write([]byte(`" style="background:`))
		pages.WriteEscaped(w, background)
		w.Write([]byte(`; --primary:`))
		pages.WriteEscaped(w, primary)
		w.Write([]byte(`; --auth-bg:`))
		pages.WriteEscaped(w, background)
		w.Write([]byte(`; --text-align:`))
		pages.WriteEscaped(w, align)
		if data.CardBackgroundColor != "" {
			w.Write([]byte(`; --auth-card:`))
			pages.WriteEscaped(w, data.CardBackgroundColor)
		}
		if pinBackground != "" {
			w.Write([]byte(`; --pin-bg:`))
			pages.WriteEscaped(w, pinBackground)
		}
		if pinText != "" {
			w.Write([]byte(`; --pin-text:`))
			pages.WriteEscaped(w, pinText)
		}
		w.Write([]byte(";\">\n"))

		w.Write([]byte(`<div class="auth-page">
	<div class="banner-region">`))
		if err := banners(data).Render(ctx, w); err != nil {
			return err
		}
		w.Write([]byte(`</div>
	<div class="auth-content">`))
		if err := content.Render(ctx, w); err != nil {
			return err
		}
		w.Write([]byte(`</div>
</div>
<div class="toast-stack" id="toast-stack" aria-live="polite">`))
		if data.Error != "" {
			if err := toast("error", data.Error).Render(ctx, w); err != nil {
				return err
			}
		}
		if data.Success != "" {
			if err := toast("success", data.Success).Render(ctx, w); err != nil {
				return err
			}
		}
		w.Write([]byte(`</div>
<script>
	(function () {
		// Toasts
		var stack = document.getElementById('toast-stack');
		function dismiss(node) {
			if (!node) return;
			node.classList.add('is-leaving');
			window.setTimeout(function () { node.remove(); }, 260);
		}
		function showToast(kind, message) {
			if (!stack || !message) return;
			var node = document.createElement('div');
			node.className = 'toast toast-' + kind;
			node.setAttribute('role', 'alert');
			var span = document.createElement('span');
			span.textContent = message;
			var close = document.createElement('button');
			close.type = 'button';
			close.className = 'toast-close';
			close.setAttribute('aria-label', 'بستن');
			close.textContent = '×';
			close.addEventListener('click', function () { dismiss(node); });
			node.appendChild(span);
			node.appendChild(close);
			stack.appendChild(node);
			window.setTimeout(function () { dismiss(node); }, 6500);
		}
		if (stack) {
			stack.querySelectorAll('.toast').forEach(function (node) {
				window.setTimeout(function () { dismiss(node); }, 6500);
				var close = node.querySelector('.toast-close');
				if (close) close.addEventListener('click', function () { dismiss(node); });
			});
		}
		window.showToast = showToast;

		// Images fade in after they finish downloading; until then their frame
		// shows a shimmering skeleton instead of an empty hole.
		function watchMedia() {
			document.querySelectorAll('[data-skeleton] img').forEach(function (img) {
				var frame = img.closest('[data-skeleton]');
				if (!frame) return;
				function done() { frame.classList.add('is-loaded'); }
				function broken() { frame.classList.add('is-loaded', 'is-broken'); }
				if (img.complete) {
					if (img.naturalWidth > 0) done(); else broken();
				} else {
					img.addEventListener('load', done);
					img.addEventListener('error', broken);
				}
			});
		}
		watchMedia();

		function toFa(value) {
			return String(value).replace(/[0-9]/g, function (digit) {
				return '۰۱۲۳۴۵۶۷۸۹'[Number(digit)];
			});
		}
		function pad2(value) {
			return (value < 10 ? '0' : '') + value;
		}

		// Banner carousel. Each banner owns its duration, and the selected mode
		// adds either a progress bar or a live remaining-seconds label.
		var bannerNodes = Array.prototype.slice.call(document.querySelectorAll('[data-banner]'));
		var bannerMode = document.body.getAttribute('data-banner-mode') || 'none';
		var bannerTimer = null;
		var bannerTick = null;

		function showBanner(index) {
			if (!bannerNodes.length) return;
			window.clearTimeout(bannerTimer);
			window.clearInterval(bannerTick);
			bannerNodes.forEach(function (node) { node.classList.remove('is-active'); });
			var node = bannerNodes[index];
			node.classList.add('is-active');
			var duration = parseInt(node.getAttribute('data-duration'), 10) || 6;

			if (bannerMode === 'progress') {
				var fill = node.querySelector('.banner-progress-fill');
				if (fill) {
					fill.style.animation = 'none';
					void fill.offsetWidth;
					fill.style.animation = '';
				}
			}

			if (bannerMode === 'time') {
				var label = node.querySelector('[data-banner-time]');
				var end = Date.now() + duration * 1000;
				function tick() {
					var left = Math.max(0, Math.ceil((end - Date.now()) / 1000));
					if (label) label.textContent = toFa(left) + ' ثانیه';
					if (left <= 0) window.clearInterval(bannerTick);
				}
				tick();
				bannerTick = window.setInterval(tick, 1000);
			}

			if (bannerNodes.length > 1) {
				bannerTimer = window.setTimeout(function () {
					showBanner((index + 1) % bannerNodes.length);
				}, duration * 1000);
			}
		}
		if (bannerNodes.length) showBanner(0);

		// PIN-style OTP input: boxes keep the value, one hidden field submits it.
		function clearPins(pin) {
			if (!pin) return;
			var hidden = pin.querySelector('[data-pin-value]');
			if (hidden) hidden.value = '';
			pin.querySelectorAll('.pin-input').forEach(function (input) { input.value = ''; });
		}

		document.querySelectorAll('[data-pin]').forEach(function (pin) {
			var inputs = Array.prototype.slice.call(pin.querySelectorAll('.pin-input'));
			var hidden = pin.querySelector('[data-pin-value]');
			if (!inputs.length || !hidden) return;

			function sync() {
				hidden.value = inputs.map(function (input) { return input.value; }).join('');
			}
			inputs.forEach(function (input, index) {
				input.addEventListener('input', function () {
					var value = input.value.replace(/\s+/g, '');
					if (value.length > 1) {
						var chars = value.split('');
						for (var i = 0; i < chars.length && index + i < inputs.length; i++) {
							inputs[index + i].value = chars[i];
						}
						var next = Math.min(index + chars.length, inputs.length - 1);
						inputs[next].focus();
					} else if (input.value && index < inputs.length - 1) {
						inputs[index + 1].focus();
					}
					sync();
				});
				input.addEventListener('keydown', function (event) {
					if (event.key === 'Backspace' && !input.value && index > 0) {
						inputs[index - 1].focus();
					}
					if (event.key === 'ArrowLeft' && index < inputs.length - 1) {
						event.preventDefault();
						inputs[index + 1].focus();
					}
					if (event.key === 'ArrowRight' && index > 0) {
						event.preventDefault();
						inputs[index - 1].focus();
					}
				});
				input.addEventListener('paste', function (event) {
					var text = (event.clipboardData || window.clipboardData).getData('text') || '';
					text = text.replace(/\s+/g, '');
					if (!text) return;
					event.preventDefault();
					for (var i = 0; i < text.length && index + i < inputs.length; i++) {
						inputs[index + i].value = text[i];
					}
					var next = Math.min(index + text.length, inputs.length - 1);
					inputs[next].focus();
					sync();
				});
			});
		});

		// Step switching and the resend countdown. The retry button stays
		// disabled until the current code expires.
		var verifyStep = document.querySelector('[data-auth-step="verify"]');
		var retryButton = document.querySelector('[data-retry-button]');
		var retryCountdown = document.querySelector('[data-retry-countdown]');
		var countdownTimer = null;
		var expiresAt = verifyStep ? parseInt(verifyStep.getAttribute('data-expires'), 10) || 0 : 0;

		function showStep(name, push) {
			document.querySelectorAll('[data-auth-step]').forEach(function (node) {
				node.hidden = node.getAttribute('data-auth-step') !== name;
			});
			var step = document.querySelector('[data-auth-step="' + name + '"]');
			if (!step) return;
			var focusable = step.querySelector('input:not([type="hidden"]):not([disabled])');
			if (focusable) focusable.focus();
			if (push && window.history && window.history.pushState) {
				var target = stepRoute(name);
				if (target) window.history.pushState({ step: name }, '', target);
			}
		}

		// The two auth steps are real routes. JavaScript advances between them
		// with pushState and restores the right one on browser back/forward, so
		// the page behaves like an SPA without losing deep links.
		function stepRoute(name) {
			if (name === 'verify') {
				var form = document.querySelector('[data-auth-verify]');
				if (!form) return null;
				var params = new URLSearchParams();
				params.set('app_id', appID);
				['identifier', 'otp_hash', 'expires_at', 'redirect'].forEach(function (key) {
					var field = form.elements[key];
					if (field && field.value) params.set(key, field.value);
				});
				return '/auth/verify-otp?' + params.toString();
			}
			var params = new URLSearchParams();
			params.set('app_id', appID);
			var redirectField = document.querySelector('[data-auth-login] [name="redirect"]');
			if (redirectField && redirectField.value) params.set('redirect', redirectField.value);
			return '/auth/login?' + params.toString();
		}

		var appID = (function () {
			var field = document.querySelector('input[name="app_id"]');
			return field ? field.value : '';
		})();

		function stepFromLocation() {
			return window.location.pathname.indexOf('/auth/verify-otp') === 0 ? 'verify' : 'login';
		}

		window.addEventListener('popstate', function () {
			var name = stepFromLocation();
			showStep(name, false);
			if (name === 'login') resetCaptcha();
		});

		function startRetryCountdown() {
			if (!retryButton) return;
			window.clearTimeout(countdownTimer);
			function tick() {
				var left = expiresAt - Math.floor(Date.now() / 1000);
				if (left > 0) {
					retryButton.disabled = true;
					retryButton.classList.add('is-waiting');
					if (retryCountdown) {
						retryCountdown.textContent = ' (' + toFa(pad2(Math.floor(left / 60))) + ':' + toFa(pad2(left % 60)) + ')';
					}
					countdownTimer = window.setTimeout(tick, 1000);
					return;
				}
				retryButton.disabled = false;
				retryButton.classList.remove('is-waiting');
				if (retryCountdown) retryCountdown.textContent = '';
			}
			tick();
		}
		startRetryCountdown();

		function resetCaptcha() {
			if (window.arcaptcha && typeof window.arcaptcha.reset === 'function') {
				try { window.arcaptcha.reset(); } catch (error) { /* widget not ready */ }
			}
		}

		function setField(form, name, value) {
			if (!form) return;
			var field = form.elements[name];
			if (field && value !== undefined && value !== null) field.value = value;
		}

		function applyChallenge(payload) {
			if (!payload) return;
			var verifyForm = document.querySelector('[data-auth-verify]');
			setField(verifyForm, 'identifier', payload.identifier);
			setField(verifyForm, 'otp_hash', payload.otp_hash);
			setField(verifyForm, 'expires_at', payload.expires_at);
			var retryForm = document.querySelector('[data-auth-retry]');
			setField(retryForm, 'identifier', payload.identifier);

			var subtitle = document.querySelector('[data-verify-subtitle]');
			if (subtitle) {
				var template = subtitle.getAttribute('data-template') || '';
				subtitle.textContent = template.replace('{identifier}', payload.identifier || '');
			}
			if (verifyStep) verifyStep.setAttribute('data-expires', String(payload.expires_at || 0));
			expiresAt = parseInt(payload.expires_at, 10) || 0;
			startRetryCountdown();
		}

		function setBusy(form, busy) {
			if (!form) return;
			var button = form.querySelector('button[type="submit"]');
			if (!button) return;
			button.disabled = busy;
			button.classList.toggle('is-busy', busy);
		}

		function postForm(form, url) {
			return fetch(url, {
				method: 'POST',
				body: new FormData(form),
				credentials: 'same-origin',
				headers: { 'X-Aman-SPA': '1', 'Accept': 'application/json' }
			}).then(function (response) {
				return response.json();
			}).catch(function () {
				return { ok: false, error: 'ارتباط با سرور برقرار نشد. دوباره تلاش کنید.' };
			});
		}

		var loginForm = document.querySelector('[data-auth-login]');
		var verifyForm = document.querySelector('[data-auth-verify]');
		var retryForm = document.querySelector('[data-auth-retry]');

		if (loginForm && window.fetch && window.FormData) {
			loginForm.addEventListener('submit', function (event) {
				event.preventDefault();
				setBusy(loginForm, true);
				postForm(loginForm, loginForm.getAttribute('action')).then(function (result) {
					if (!result.ok) {
						setBusy(loginForm, false);
						resetCaptcha();
						showToast('error', result.error);
						return;
					}
					applyChallenge(result);
					resetCaptcha();
					setBusy(loginForm, false);
					showStep('verify', true);
				});
			});
		}

		if (verifyForm && window.fetch && window.FormData) {
			verifyForm.addEventListener('submit', function (event) {
				event.preventDefault();
				setBusy(verifyForm, true);
				postForm(verifyForm, verifyForm.getAttribute('action')).then(function (result) {
					if (!result.ok) {
						setBusy(verifyForm, false);
						showToast('error', result.error);
						var first = verifyForm.querySelector('.pin-input');
						if (first) first.focus();
						return;
					}
					window.location.assign(result.redirect || '/');
				});
			});
		}

		if (retryForm && window.fetch && window.FormData) {
			retryForm.addEventListener('submit', function (event) {
				event.preventDefault();
				if (retryButton && retryButton.disabled) return;
				setBusy(retryForm, true);
				postForm(retryForm, retryForm.getAttribute('action')).then(function (result) {
					setBusy(retryForm, false);
					if (!result.ok) {
						showToast('error', result.error);
						return;
					}
					applyChallenge(result);
					clearPins(document.querySelector('[data-pin]'));
					var first = verifyForm ? verifyForm.querySelector('.pin-input') : null;
					if (first) first.focus();
					showStep('verify', true);
					showToast('success', 'کد تازه ارسال شد.');
				});
			});
		}

		var backLink = document.querySelector('[data-auth-back]');
		if (backLink) {
			backLink.addEventListener('click', function (event) {
				event.preventDefault();
				resetCaptcha();
				if (loginForm) setBusy(loginForm, false);
				showStep('login', true);
			});
		}
	})();
</script>
</body>
</html>`))
		return nil
	})
}

func banners(data pages.AuthData) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		mode := data.BannerMode
		if mode == "" {
			mode = "none"
		}
		if len(data.Banners) == 0 {
			// The split layout keeps its branded side panel even without a
			// banner, so the page never collapses into an off-center card.
			if data.Layout == "split" {
				w.Write([]byte(`<div class="banner-fallback"><span class="banner-fallback-mark">`))
				pages.WriteEscaped(w, firstLetter(data.App.Name))
				w.Write([]byte(`</span><strong>`))
				pages.WriteEscaped(w, data.App.Name)
				w.Write([]byte(`</strong><small>ورود امن و یکپارچه</small></div>`))
			}
			return nil
		}
		for _, banner := range data.Banners {
			duration := banner.Duration
			if duration <= 0 {
				duration = 6
			}
			w.Write([]byte(`<div class="banner" data-banner data-duration="`))
			pages.WriteEscaped(w, strconv.Itoa(duration))
			w.Write([]byte(`">`))
			if banner.URL != "" {
				w.Write([]byte(`<span class="banner-media" data-skeleton><span class="media-skeleton" aria-hidden="true"></span><img src="`))
				pages.WriteEscaped(w, banner.URL)
				w.Write([]byte(`" alt="`))
				pages.WriteEscaped(w, strings.TrimSpace(banner.Text))
				w.Write([]byte(`" loading="lazy" decoding="async" /></span>`))
			}
			if banner.Text != "" {
				w.Write([]byte(`<div class="banner-text">`))
				pages.WriteEscaped(w, banner.Text)
				w.Write([]byte(`</div>`))
			}
			if banner.Link != "" {
				w.Write([]byte(`<a class="banner-link" href="`))
				pages.WriteEscaped(w, banner.Link)
				w.Write([]byte(`">مشاهده</a>`))
			}
			if mode == "progress" {
				w.Write([]byte(`<span class="banner-progress" aria-hidden="true"><span class="banner-progress-fill" style="--banner-duration:`))
				pages.WriteEscaped(w, strconv.Itoa(duration))
				w.Write([]byte(`s"></span></span>`))
			}
			if mode == "time" {
				w.Write([]byte(`<span class="banner-time" data-banner-time></span>`))
			}
			w.Write([]byte(`</div>`))
		}
		return nil
	})
}

func toast(kind, message string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		w.Write([]byte(`<div class="toast toast-`))
		pages.WriteEscaped(w, kind)
		w.Write([]byte(`" role="alert"><span>`))
		pages.WriteEscaped(w, message)
		w.Write([]byte(`</span><button type="button" class="toast-close" aria-label="بستن">×</button></div>`))
		return nil
	})
}
