// Toast banners - generic, not specific to source delete. showToast() is the primitive; anything
// that redirects can hand off a message via ?toast=&message=&variant= (see
// web/static/js/source_form.js's deleteSource()) - on load we show it once and strip it from the
// URL so a refresh doesn't repeat it. Rendered inline at the top of whichever page's own
// #toast-container (see web/components/shared/toast_container.html), not a fixed floating overlay.
const TOAST_LIFETIME_MS = 3500;
const TOAST_FADE_MS = 300;

const TOAST_VARIANTS = {
  success: {
    bg: "bg-emerald-50 dark:bg-emerald-500/10",
    border: "border-emerald-500",
    iconBg: "bg-emerald-500",
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4"><polyline points="20 6 9 17 4 12"></polyline></svg>',
  },
  info: {
    bg: "bg-blue-50 dark:bg-blue-500/10",
    border: "border-blue-500",
    iconBg: "bg-blue-500",
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>',
  },
  warning: {
    bg: "bg-amber-50 dark:bg-amber-500/10",
    border: "border-amber-500",
    iconBg: "bg-amber-500",
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4"><path d="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>',
  },
  error: {
    bg: "bg-red-50 dark:bg-red-500/10",
    border: "border-red-500",
    iconBg: "bg-red-500",
    icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>',
  },
};

const TOAST_ICON_CLOSE =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>';

function showToast(title, message, variant) {
  const container = document.getElementById("toast-container");
  if (!container) return;

  const style = TOAST_VARIANTS[variant] || TOAST_VARIANTS.success;

  const toast = document.createElement("div");
  toast.className =
    "flex items-start gap-3 rounded-lg border-l-4 p-4 shadow-sm transition-opacity duration-300 " +
    style.bg + " " + style.border;

  const icon = document.createElement("span");
  icon.className = "flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-white " + style.iconBg;
  icon.innerHTML = style.icon;

  const body = document.createElement("div");
  body.className = "min-w-0 flex-1";

  const titleEl = document.createElement("div");
  titleEl.className = "text-sm font-semibold text-gray-900 dark:text-gray-100";
  titleEl.textContent = title;

  const messageEl = document.createElement("div");
  messageEl.className = "mt-0.5 text-sm text-gray-600 dark:text-gray-300";
  messageEl.textContent = message;

  body.append(titleEl, messageEl);

  const closeButton = document.createElement("button");
  closeButton.type = "button";
  closeButton.setAttribute("aria-label", "Dismiss");
  closeButton.className = "shrink-0 text-gray-400 hover:text-gray-700 dark:hover:text-gray-200";
  closeButton.innerHTML = TOAST_ICON_CLOSE;

  toast.append(icon, body, closeButton);
  container.appendChild(toast);

  let dismissed = false;

  const dismiss = function () {
    if (dismissed) return;
    dismissed = true;

    toast.classList.add("opacity-0");
    setTimeout(function () {
      toast.remove();
    }, TOAST_FADE_MS);
  };

  closeButton.addEventListener("click", dismiss);
  setTimeout(dismiss, TOAST_LIFETIME_MS);
}

document.addEventListener("DOMContentLoaded", function () {
  const params = new URLSearchParams(window.location.search);
  const title = params.get("toast");
  if (!title) return;

  showToast(title, params.get("message") || "", params.get("variant") || "success");

  params.delete("toast");
  params.delete("message");
  params.delete("variant");

  const rest = params.toString();
  history.replaceState(null, "", window.location.pathname + (rest ? "?" + rest : ""));
});
