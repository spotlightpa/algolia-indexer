import "alpinejs";

import formatTel from "../utils/format-tel.js";
import searchPeople from "../utils/search-people.js";

window.spl = Object.assign({}, window.spl, {
  formatTel,
  searchPeople,
});

// Redirect admin emails to admin
const routes = /(confirmation|invite|recovery|email_change)_token=([^&]+)/g;

if (window.location.hash.match(routes)) {
  window.location.replace(
    window.location.origin + "/admin/" + window.location.hash
  );
}
