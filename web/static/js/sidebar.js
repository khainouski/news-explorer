// Sidebar collapse/expand - plain JS, no HTMX needed. Collapsed by default; persists the user's
// choice to localStorage like theme.js. Toggles a "sidebar-collapsed" class on <html> (styled in
// base.html) rather than an inline style on #sidebar directly, so it can run synchronously in
// <head> - before #sidebar itself is even parsed - avoiding a flash of the wrong state, the same
// reason theme.js toggles a class instead of setting style directly.
const SIDEBAR_ICON_LEFT =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-3.5 w-3.5"><polyline points="15 18 9 12 15 6"></polyline></svg>';
const SIDEBAR_ICON_RIGHT =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-3.5 w-3.5"><polyline points="9 18 15 12 9 6"></polyline></svg>';

(function () {
  const collapsed = localStorage.getItem("sidebarCollapsed") !== "false";
  document.documentElement.classList.toggle("sidebar-collapsed", collapsed);
})();

function updateSidebarToggleIcon(collapsed) {
  const icon = document.getElementById("sidebar-toggle-icon");
  if (icon) icon.innerHTML = collapsed ? SIDEBAR_ICON_RIGHT : SIDEBAR_ICON_LEFT;
}

function toggleSidebar() {
  const collapsed = !document.documentElement.classList.contains("sidebar-collapsed");
  document.documentElement.classList.toggle("sidebar-collapsed", collapsed);
  localStorage.setItem("sidebarCollapsed", String(!collapsed));
  updateSidebarToggleIcon(collapsed);
}

document.addEventListener("DOMContentLoaded", function () {
  updateSidebarToggleIcon(document.documentElement.classList.contains("sidebar-collapsed"));
});
