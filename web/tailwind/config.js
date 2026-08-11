module.exports = {
  content: ['./web/**/*.html'],
  darkMode: 'class',
  // Badge colors (internal/controller/http/source/add.go's badgeColorPalette, plus what the
  // seed migration inserts directly) come from the DB, not literal class names in any .html
  // file - Tailwind can't find them by scanning, so they need a safelist.
  safelist: [
    { pattern: /^bg-(slate|gray|red|orange|amber|green|emerald|teal|sky|blue|indigo|purple|fuchsia|pink)-(500|600|700|900)$/ },
  ],
};
