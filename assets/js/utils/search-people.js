import searchAPI from "../utils/search-api.js";
import { debouncer } from "../utils/timers.js";

function normalize(obj) {
  return {
    path: obj["path"] || "",
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
    query: "",
    results: null,
    error: null,
    hasSearched: false,
    isLoading: false,

    init() {
      const bouncedSearch = debouncer({ milliseconds: 500 }, () =>
        searchAPI(this.query)
          .then((results) => {
            this.error = null;
            if (results) {
              this.results = results;
            }
          })
          .catch((error) => {
            this.isLoading = false;
            this.error = error;
          })
          .finally(() => {
            this.isLoading = false;
          })
      );

      this.$watch("query", () => {
        this.hasSearched = true;
        this.isLoading = true;
        bouncedSearch();
      });

      if (window.history.state?.searchResult) {
        this.results = window.history.state.searchResult;
        this.query = window.history.state.searchQuery;
      }
    },

    get people() {
      if (!this.results || !this.results.hits) {
        return [];
      }
      window.history.replaceState(
        {
          searchQuery: this.query,
          searchResult: this.results,
        },
        ""
      );
      return this.results.hits.map(normalize);
    },

    get resultsText() {
      let nHits = this.results?.nbHits ?? 0;
      if (!nHits) {
        return "No search results";
      }
      if (nHits === 1) {
        return "Got one search result.";
      }
      let nStories = this.results?.hits?.length ?? 0;
      let more = nHits > nStories ? `Showing first ${nStories}.` : "";
      return `Got ${nHits} search results. ${more}`;
    },
  };
}
