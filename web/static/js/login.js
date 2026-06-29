// Password show/hide toggles - one button can control any password input by id (see
// web/pages/login.html and web/pages/change_password.html: data-toggle-password="<input id>").
document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll("[data-toggle-password]").forEach(function (button) {
    const input = document.getElementById(button.dataset.togglePassword);
    if (!input) return;

    button.addEventListener("click", function () {
      const showing = input.type === "text";
      input.type = showing ? "password" : "text";
      button.querySelector('[data-eye-icon="off"]').classList.toggle("hidden", !showing);
      button.querySelector('[data-eye-icon="on"]').classList.toggle("hidden", showing);
    });
  });
});
