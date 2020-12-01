# Spotlight PA Diverse Sources Database

A resource for Pennsylvania journalists to improve representation and diversify perspectives in their coverage.

## License

All content copyright Spotlight PA. Code available under the MIT license. Photos used with permission of subjects. Information contained in this database is self-reported by participants and should be verified before publication.

## Installation

Project requires Yarn, Hugo, and optionally Go. See netlify.toml file for Node and Hugo versions. Go version specified in .go-version file.

To setup, run `yarn`.

To develop locally, run `hugo serve` and open a web browser to http://localhost:1313/.

## Architecture

The site uses a [THANG Stack](https://twitter.com/carlmjohnson/status/1327090078578053120) architecture:

- [**T**ailwind CSS](https://tailwindcss.com): Provides basic CSS architecture/theming
- [**H**ugo](https://gohugo.io): Site builder
- [**A**lpine.js](https://github.com/alpinejs/alpine): JavaScript micro-framework
- [**N**etlify CMS](https://www.netlifycms.org): Allows editors to change pages without coding
- [**G**o](https://golang.org): Handles miscellaneous tasks

Site search is powered by [Algolia](https://www.algolia.com). On deploy, a search index JSON file is built at /searchindex.json and sent to Algolia by a small Go script. To work, the Algolia script API key must have permission to create/drop an index because it creates a temporary table, sends all the data to the temporary table, then swaps it in.

Email addresses are Base64 encoded to prevent casual scraping.

The site was not made with reuse in mind, but it shouldn't be so hard. Just rip out the content files, rewrite nav.html and footer.html to remove references to Spotlight PA, and change the base URL and Google Analytics key in config.toml. Contact webmaster@spotlightpa.org with questions.
