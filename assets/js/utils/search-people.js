import searchAPI from "../utils/search-api.js";
import { debouncer } from "../utils/timers.js";

function normalize(obj) {
  return {
    name: obj.name || "",
  };
}

export default function searchPeople() {
  return {
    query: "",
    results: null,
    error: null,
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
        this.isLoading = true;
        bouncedSearch();
      });
    },

    get people() {
      if (!this.results || !this.results.hits) {
        return [];
      }
      return this.results.hits.map(normalize);
    },

    get resultsCount() {
      return this.results?.nbHits ?? 0;
    },

    get resultsText() {
      let nHits = this.resultsCount;
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
