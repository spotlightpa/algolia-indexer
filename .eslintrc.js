module.exports = {
  root: true,
  env: {
    es6: true,
    browser: true,
    node: true,
  },
  parser: "babel-eslint",
  extends: ["plugin:prettier/recommended", "eslint:recommended", "prettier"],
  plugins: ["prettier"],
  // add your custom rules here
  rules: {
    "no-console": "warn",
  },
};
