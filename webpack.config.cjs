const path = require("path");
const Dotenv = require("dotenv-webpack");

const isProduction = process.env.NODE_ENV == "production";

const config = {
  entry: "./server/static/js/app.js",
  output: {
    path: path.resolve(__dirname, "server", "static"),
    filename: 'app.js',
  },
  resolve: {
    alias: {
      "axios/dist/browser/axios.cjs": path.resolve(__dirname, "node_modules/axios/dist/browser/axios.cjs"),
    },
  },
  plugins: [
    // Add your plugins here
    // Learn more about plugins from https://webpack.js.org/configuration/plugins/
    new Dotenv(),
  ],
  module: {
    rules: [
      {
        test: /\.(eot|svg|ttf|woff|woff2|png|jpg|gif)$/i,
        type: "asset",
      },

      // Add your rules for custom modules here
      // Learn more about loaders from https://webpack.js.org/loaders/
    ],
  },
};

module.exports = () => {
  if (isProduction) {
    config.mode = "production";
  } else {
    config.mode = "development";
  }
  return config;
};
