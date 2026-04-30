// "Delete Source" on the Edit Source page (web/pages/source_form.html). Deletes immediately, no
// confirmation step, then always lands on /sources - success or not - with a flash toast shown
// there (web/static/js/toast.js reads ?toast=&message=&variant=).
async function deleteSource(id, name) {
  let ok = false;

  try {
    const res = await fetch("/sources/" + id, { method: "DELETE" });
    ok = res.ok;
  } catch (e) {
    ok = false;
  }

  const params = ok
    ? new URLSearchParams({ toast: "Source deleted", message: name + " was successfully deleted." })
    : new URLSearchParams({
        toast: "Couldn't delete source",
        message: "Something went wrong. Please try again.",
        variant: "error",
      });

  window.location.href = "/sources?" + params.toString();
}

document.addEventListener("DOMContentLoaded", function () {
  const button = document.getElementById("delete-source-button");
  if (!button) return;

  button.addEventListener("click", function () {
    deleteSource(button.dataset.sourceId, button.dataset.sourceName);
  });
});
