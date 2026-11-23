package pages

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"

	"github.com/a-h/templ"
)

// AdminLayout is the Persian, right-to-left shell of the admin panel: brand
// sidebar, sticky top bar and a content area. It is wrapped around every admin
// page by the admin HTTP handler, which is also what links the stylesheet.
func AdminLayout(title, active string, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		w.Write([]byte(`<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta name="theme-color" content="#0b1220" />
	<link rel="icon" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%3E%3Crect width='32' height='32' rx='9' fill='%230b1220'/%3E%3Cpath d='M16 6.5l7.5 3.2v5.1c0 4.2-3 7.8-7.5 10.2-4.5-2.4-7.5-6-7.5-10.2V9.7L16 6.5z' fill='%233b82f6'/%3E%3C/svg%3E" />
	<title>`))
		pages.WriteEscaped(w, title)
		w.Write([]byte(` · Aman CAS</title>
	<link rel="stylesheet" href="/static/css/admin.css?v=11" />
</head>
<body>
	<div class="shell">
		<div class="sidebar-backdrop" data-sidebar-close></div>
		<aside class="sidebar" id="admin-sidebar">
			<a class="brand" href="/admin">
				<span class="brand-copy">
					<strong>امان</strong>
					<small>کنسول مدیریت</small>
				</span>
			</a>

			<nav class="nav" aria-label="ناوبری مدیریت">
`))
		navLink(w, "/admin", "داشبورد", iconDashboard, active)
		navLink(w, "/admin/apps", "برنامه‌ها", iconApps, active)
		navLink(w, "/admin/sessions", "نشست‌ها", iconSessions, active)
		navLink(w, "/admin/guide", "راهنما", iconGuide, active)
		navLink(w, "/admin/developers", "توسعه‌دهندگان", iconCode, active)
		w.Write([]byte(`
			</nav>
		</aside>

		<div class="workspace">
			<header class="topbar">
				<button class="icon-btn menu-btn" type="button" data-sidebar-toggle aria-label="باز کردن منو" aria-controls="admin-sidebar">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M4 7h16M4 12h16M4 17h16"/></svg>
				</button>
				<div class="topbar-copy">
					<span class="crumb">کنسول مدیریت</span>
				</div>
			</header>

			<main class="content">
`))
		if err := content.Render(ctx, w); err != nil {
			return err
		}
		w.Write([]byte(`
			</main>
			<footer class="page-footer">
				<span>نسخه ۱.۰</span>
			</footer>
		</div>
	</div>

	<dialog class="confirm-dialog" id="confirm-dialog">
		<form method="dialog" class="confirm-card">
			<h3 class="confirm-title">تایید عملیات</h3>
			<p class="confirm-message" data-confirm-message>آیا از انجام این عملیات مطمئن هستید؟</p>
			<div class="confirm-actions">
				<button type="button" class="btn btn-ghost" data-confirm-cancel>انصراف</button>
				<button type="button" class="btn btn-danger" data-confirm-accept>تایید و ادامه</button>
			</div>
		</form>
	</dialog>

	<script>
		(function () {
			var body = document.body;

			// Toasts replace inline alerts. Server-rendered errors/success land
			// here, and the URL is cleaned so a refresh does not repeat them.
			document.querySelectorAll('[data-toast-stack]').forEach(function (stack) {
				stack.querySelectorAll('.toast').forEach(function (toast) {
					function dismiss() {
						toast.classList.add('is-leaving');
						window.setTimeout(function () { toast.remove(); }, 240);
					}
					var close = toast.querySelector('.toast-close');
					if (close) close.addEventListener('click', dismiss);
					window.setTimeout(dismiss, 6500);
				});
				window.setTimeout(function () {
					if (!stack.children.length) stack.remove();
				}, 7000);
			});
			if (window.history && window.history.replaceState) {
				var url = new URL(window.location.href);
				if (url.searchParams.has('error') || url.searchParams.has('success')) {
					url.searchParams.delete('error');
					url.searchParams.delete('success');
					window.history.replaceState({}, '', url.pathname + (url.searchParams.toString() ? '?' + url.searchParams.toString() : ''));
				}
			}

			// Any field with a hint gets a small help icon with a popover, so
			// the form stays compact while every input keeps its explanation.
			document.querySelectorAll('.field').forEach(function (field) {
				var hint = field.querySelector('.hint');
				var label = field.querySelector('label');
				if (!hint || !label) return;
				var text = (hint.textContent || '').trim();
				if (!text) return;
				var tip = document.createElement('span');
				tip.className = 'help-tip';
				tip.tabIndex = 0;
				tip.setAttribute('role', 'button');
				tip.setAttribute('aria-label', 'راهنمای این فیلد');
				tip.textContent = '?';
				var bubble = document.createElement('span');
				bubble.className = 'help-bubble';
				bubble.textContent = text;
				tip.appendChild(bubble);
				label.appendChild(tip);
				hint.remove();
			});

			// Keep the hex label next to each color picker in sync while the
			// picker is open, so the value is readable without saving.
			document.querySelectorAll('input[type="color"]').forEach(function (input) {
				var field = input.closest('.color-field');
				var value = field ? field.querySelector('.color-value') : null;
				if (!value) return;
				input.addEventListener('input', function () {
					value.textContent = input.value;
				});
			});
			var toggle = document.querySelector('[data-sidebar-toggle]');
			var close = document.querySelector('[data-sidebar-close]');
			if (toggle) {
				toggle.addEventListener('click', function () {
					body.classList.toggle('sidebar-open');
				});
			}
			if (close) {
				close.addEventListener('click', function () {
					body.classList.remove('sidebar-open');
				});
			}

			// Declarative dependencies. Any input anywhere in the document can
			// show, hide, enable or disable another block, so conditional
			// sections stay in sync with checkboxes without per-page code.
			function syncDependents() {
				document.querySelectorAll('[data-show-when]').forEach(function (node) {
					var control = document.querySelector(node.getAttribute('data-show-when'));
					if (control) node.hidden = !control.checked;
				});
				document.querySelectorAll('[data-hide-when]').forEach(function (node) {
					var control = document.querySelector(node.getAttribute('data-hide-when'));
					if (control) node.hidden = control.checked;
				});
				document.querySelectorAll('[data-enabled-when]').forEach(function (node) {
					var control = document.querySelector(node.getAttribute('data-enabled-when'));
					if (!control) return;
					node.classList.toggle('is-disabled', !control.checked);
					node.querySelectorAll('input, select, textarea, button').forEach(function (field) {
						field.disabled = !control.checked;
					});
				});
				document.querySelectorAll('[data-disabled-when]').forEach(function (node) {
					var control = document.querySelector(node.getAttribute('data-disabled-when'));
					if (!control) return;
					node.classList.toggle('is-disabled', control.checked);
					node.querySelectorAll('input, select, textarea, button').forEach(function (field) {
						field.disabled = control.checked;
					});
				});
			}
			document.addEventListener('change', syncDependents);
			document.addEventListener('input', syncDependents);
			syncDependents();

			// Filters apply themselves as soon as their value changes; there is
			// no separate "apply" button.
			document.querySelectorAll('form[data-auto-filter]').forEach(function (form) {
				var timer = null;
				form.addEventListener('change', function () { form.requestSubmit(); });
				form.querySelectorAll('input[type="search"], input[type="text"]').forEach(function (input) {
					input.addEventListener('input', function () {
						window.clearTimeout(timer);
						timer = window.setTimeout(function () { form.requestSubmit(); }, 450);
					});
				});
			});

			// Banner editor: add, remove and reorder rows by dragging the handle.
			var bannerList = document.querySelector('[data-banner-list]');
			var bannerAdd = document.querySelector('[data-banner-add]');
			var bannerTemplate = document.querySelector('[data-banner-template]');
			if (bannerList && bannerAdd && bannerTemplate) {
				bannerAdd.addEventListener('click', function () {
					var node = bannerTemplate.content.firstElementChild.cloneNode(true);
					bannerList.appendChild(node);
					var first = node.querySelector('input, select');
					if (first) first.focus();
				});
				var dragRow = null;
				bannerList.addEventListener('click', function (event) {
					var button = event.target.closest('[data-banner-remove]');
					if (!button) return;
					var row = button.closest('[data-banner-row]');
					if (row) row.remove();
				});
				bannerList.addEventListener('dragstart', function (event) {
					var handle = event.target.closest('[data-drag-handle]');
					if (!handle) return;
					dragRow = handle.closest('[data-banner-row]');
					if (!dragRow) return;
					dragRow.classList.add('is-dragging');
					if (event.dataTransfer) {
						event.dataTransfer.effectAllowed = 'move';
						event.dataTransfer.setData('text/plain', 'banner');
					}
				});
				bannerList.addEventListener('dragover', function (event) {
					if (!dragRow) return;
					event.preventDefault();
					var target = event.target.closest('[data-banner-row]');
					if (!target || target === dragRow) return;
					var rect = target.getBoundingClientRect();
					var after = event.clientY > rect.top + rect.height / 2;
					bannerList.insertBefore(dragRow, after ? target.nextSibling : target);
				});
				bannerList.addEventListener('drop', function (event) {
					if (dragRow) event.preventDefault();
				});
				bannerList.addEventListener('dragend', function () {
					if (dragRow) dragRow.classList.remove('is-dragging');
					dragRow = null;
				});
			}

			// Session detail drawers.
			document.querySelectorAll('[data-drawer-open]').forEach(function (button) {
				button.addEventListener('click', function () {
					var drawer = document.getElementById(button.getAttribute('data-drawer-open'));
					if (drawer) drawer.showModal();
				});
			});
			document.querySelectorAll('[data-drawer-close]').forEach(function (button) {
				button.addEventListener('click', function () {
					var drawer = button.closest('dialog');
					if (drawer) drawer.close();
				});
			});
			document.querySelectorAll('dialog.drawer').forEach(function (drawer) {
				drawer.addEventListener('click', function (event) {
					if (event.target === drawer) drawer.close();
				});
			});

			// Replace window.confirm with a proper modal dialog.
			var dialog = document.getElementById('confirm-dialog');
			var message = dialog ? dialog.querySelector('[data-confirm-message]') : null;
			var accept = dialog ? dialog.querySelector('[data-confirm-accept]') : null;
			var cancel = dialog ? dialog.querySelector('[data-confirm-cancel]') : null;
			var pendingForm = null;

			if (dialog && accept && cancel) {
				document.querySelectorAll('form[data-confirm]').forEach(function (form) {
					form.addEventListener('submit', function (event) {
						if (form.dataset.confirmed === 'yes') return;
						event.preventDefault();
						pendingForm = form;
						if (message) message.textContent = form.getAttribute('data-confirm') || 'آیا از انجام این عملیات مطمئن هستید؟';
						dialog.showModal();
					});
				});
				accept.addEventListener('click', function () {
					if (!pendingForm) return;
					pendingForm.dataset.confirmed = 'yes';
					dialog.close();
					pendingForm.requestSubmit();
					pendingForm = null;
				});
				cancel.addEventListener('click', function () {
					pendingForm = null;
					dialog.close();
				});
				dialog.addEventListener('cancel', function () { pendingForm = null; });
			}

			// Copy buttons on documentation code blocks.
			document.querySelectorAll('[data-copy]').forEach(function (button) {
				button.addEventListener('click', function () {
					var block = button.closest('.code-block');
					var code = block ? block.querySelector('code') : null;
					if (!code) return;
					var text = code.innerText;
					function done() {
						var previous = button.textContent;
						button.textContent = 'کپی شد';
						window.setTimeout(function () { button.textContent = previous; }, 1600);
					}
					function fallback() {
						var area = document.createElement('textarea');
						area.value = text;
						area.setAttribute('readonly', 'readonly');
						area.style.position = 'fixed';
						area.style.opacity = '0';
						document.body.appendChild(area);
						area.select();
						try { document.execCommand('copy'); done(); } catch (error) { /* ignore */ }
						document.body.removeChild(area);
					}
					if (navigator.clipboard && navigator.clipboard.writeText) {
						navigator.clipboard.writeText(text).then(done).catch(fallback);
					} else {
						fallback();
					}
				});
			});

		})();
	</script>
</body>
</html>`))
		return nil
	})
}

func navLink(w io.Writer, href, label, icon, active string) {
	class := "nav-link"
	current := ""
	if active == href || (href == "/admin" && active == "/admin") {
		class += " is-active"
		current = ` aria-current="page"`
	}
	w.Write([]byte(`<a class="` + class + `" href="` + href + `"` + current + `>` + icon + `<span>` + label + `</span></a>`))
}

func flash(errorMessage, successMessage string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if errorMessage == "" && successMessage == "" {
			return nil
		}
		w.Write([]byte(`<div class="toast-stack" data-toast-stack aria-live="polite">`))
		if errorMessage != "" {
			w.Write([]byte(`<div class="toast toast-error" role="alert">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="9"/><path d="M12 8v4.5M12 16h.01" stroke-linecap="round"/></svg>
				<span>`))
			pages.WriteEscaped(w, errorMessage)
			w.Write([]byte(`</span><button type="button" class="toast-close" aria-label="بستن">×</button></div>`))
		}
		if successMessage != "" {
			message := successMessage
			switch successMessage {
			case "created":
				message = "برنامه با موفقیت ساخته شد."
			case "updated":
				message = "تغییرات برنامه ذخیره شد."
			case "deleted":
				message = "برنامه حذف شد."
			case "revoked":
				message = "نشست ابطال شد."
			}
			w.Write([]byte(`<div class="toast toast-success" role="status">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="9"/><path d="M8.5 12.2l2.4 2.4 4.6-4.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
				<span>`))
			pages.WriteEscaped(w, message)
			w.Write([]byte(`</span><button type="button" class="toast-close" aria-label="بستن">×</button></div>`))
		}
		w.Write([]byte(`</div>`))
		return nil
	})
}

const (
	iconDashboard = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><rect x="3.5" y="3.5" width="7" height="9" rx="1.6"/><rect x="13.5" y="3.5" width="7" height="5" rx="1.6"/><rect x="13.5" y="11.5" width="7" height="9" rx="1.6"/><rect x="3.5" y="15.5" width="7" height="5" rx="1.6"/></svg>`
	iconApps      = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round"><rect x="3.5" y="3.5" width="7" height="7" rx="1.6"/><rect x="13.5" y="3.5" width="7" height="7" rx="1.6"/><rect x="3.5" y="13.5" width="7" height="7" rx="1.6"/><rect x="13.5" y="13.5" width="7" height="7" rx="1.6"/></svg>`
	iconSessions  = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3.5l7 2.8v5.2c0 4.1-2.9 7.7-7 9.5-4.1-1.8-7-5.4-7-9.5V6.3l7-2.8z"/><path d="M9.2 12.1l2 2 3.8-4"/></svg>`
	iconGuide     = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><path d="M4.5 5.5A2 2 0 0 1 6.5 3.5H19v17H6.5a2 2 0 0 1-2-2z"/><path d="M8 7.5h7M8 11h7M8 14.5h4"/></svg>`
	iconCode      = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><path d="M8.5 8.5L5 12l3.5 3.5M15.5 8.5L19 12l-3.5 3.5M13 6l-2 12"/></svg>`
)
