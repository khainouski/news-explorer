// Theme toggle - plain JS, no HTMX needed (no server round-trip, just localStorage + a class on
// <html>). Runs synchronously in <head>, before <body> paints, so there's no flash of the wrong
// theme. Tailwind's `dark:` variant classes react to the "dark" class automatically
// (tailwind.config = { darkMode: 'class' } in base.html).
(function () {
  const stored = localStorage.getItem("theme");
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const isDark = stored ? stored === "dark" : prefersDark;

  document.documentElement.classList.toggle("dark", isDark);
})();

function toggleTheme() {
  const isDark = document.documentElement.classList.toggle("dark");
  localStorage.setItem("theme", isDark ? "dark" : "light");
}
