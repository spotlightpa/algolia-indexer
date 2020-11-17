import searchAPI from "../utils/search-api.js";
import { debouncer } from "../utils/timers.js";

function normalize(obj) {
  return {
    path: obj["path"] || "",
    name: obj["name"] || "",
    full_name: obj["full_name"] || "",
    bio: obj["bio"] || "",
    last_name: obj["last_name"] || "",
    location: obj["location"] || "",
    role: obj["role"] || "",
    expertise: obj["expertise"] || [],
    keywords: obj["keywords"] || [],
  };
}

export default function searchPeople() {
  return {
    query: window.history.state?.searchQuery || "",
    results: window.history.state?.searchResult || null,
    error: null,
    isLoading: false,

    init() {
      const bouncedSearch = debouncer({ milliseconds: 500 }, () =>
        this.search()
      );
      this.$watch("query", (query) => {
        this.isLoading = !!query;
        this.storeHistory();
        bouncedSearch();
      });
    },

    search() {
      searchAPI(this.query)
        .then((results) => {
          this.error = null;
          if (results) {
            this.results = results;
          }
          this.storeHistory();
        })
        .catch((error) => {
          this.isLoading = false;
          this.error = error;
        })
        .finally(() => {
          this.isLoading = false;
        });
    },

    get people() {
      if (!this.results || !this.results.hits) {
        return [];
      }
      return this.results.hits.map(normalize);
    },

    get resultsText() {
      let nHits = this.results?.nbHits ?? 0;
      if (!nHits) {
        return "No search results.";
      }
      if (nHits === 1) {
        return "Got one search result.";
      }
      let nStories = this.results?.hits?.length ?? 0;
      let more = nHits > nStories ? `Showing first ${nStories}.` : "";
      return `Got ${nHits} search results. ${more}`;
    },

    storeHistory() {
      let searchQuery = "" + this.query;
      let searchResult = JSON.parse(JSON.stringify(this.results));

      window.history.replaceState(
        {
          searchQuery,
          searchResult,
        },
        ""
      );
    },

    clear() {
      this.results = null;
      this.query = "";
    },
  };
}
