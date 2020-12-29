const fs = require('fs');
const path = require('path');

const browserify = require('browserify');
require('dotenv').config();

const inputFile = path.join(__dirname, '..', 'server', 'static', 'js', 'app.js');
const outputFile = path.join(__dirname, '..', 'server', 'static', 'app.js');

browserify(inputFile, {debug: true})
  .transform('unassertify', {global: true})
  .transform('envify')
  .plugin('common-shakeify')
  .transform('babelify',  {presets: ['@babel/preset-env']})
  .transform('uglifyify', {
    mangle: false,
    toplevel: true,
    keep_fnames: true,
    keep_classnames: true,
  })
  .bundle()
  .pipe(fs.createWriteStream(outputFile));
