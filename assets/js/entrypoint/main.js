import "alpinejs";

import searchPeople from "../utils/search-people.js";

window.spl = Object.assign({}, window.spl, {
  searchPeople,
});

// Redirect admin emails to admin
const routes = /(confirmation|invite|recovery|email_change)_token=([^&]+)/g;

if (window.location.hash.match(routes)) {
  window.location.replace(
    window.location.origin + "/admin/" + window.location.hash
  );
}
