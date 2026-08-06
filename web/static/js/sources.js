// "Sync" button on the Sources page (web/pages/sources.html) - POSTs, then redirects to /sources
// with a flash toast built from the JSON response.
async function syncSources(button) {
  button.disabled = true;

  let status = 0;
  let data = null;

  try {
    const res = await fetch("/sources/sync", { method: "POST" });
    status = res.status;
    if (res.ok) {
      data = await res.json();
    }
  } catch (e) {
    status = 0;
  }

  let params;
  if (status === 200 && data) {
    params = new URLSearchParams({ toast: "Sync complete", message: syncSummary(data) });
    if (data.sourcesFailed > 0) params.set("variant", "warning");
  } else if (status === 403) {
    params = new URLSearchParams({
      toast: "Admin only",
      message: "This action is only available for the admin account.",
      variant: "warning",
    });
  } else {
    params = new URLSearchParams({
      toast: "Couldn't sync sources",
      message: "Something went wrong. Please try again.",
      variant: "error",
    });
  }

  window.location.href = "/sources?" + params.toString();
}

function syncSummary(data) {
  const articles = plural(data.articlesInserted, "article");
  const synced = plural(data.sourcesSynced, "source");
  let summary = `${data.articlesInserted} new ${articles} from ${data.sourcesSynced} ${synced}.`;

  if (data.sourcesFailed > 0) {
    summary += ` ${data.sourcesFailed} ${plural(data.sourcesFailed, "source")} failed to sync.`;
  }

  return summary;
}

function plural(n, word) {
  return n === 1 ? word : word + "s";
}

document.addEventListener("DOMContentLoaded", function () {
  const button = document.getElementById("sync-sources-button");
  if (!button) return;

  button.addEventListener("click", function () {
    syncSources(button);
  });
});
