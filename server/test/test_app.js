const expect = require('chai').expect;
const process = require('process');
const varsnap = require('varsnap');

varsnap.updateConfig({
  varsnap: process.env.VARSNAP,
  env: process.env.ENVIRONMENT,
  producerToken: process.env.VARSNAP_PRODUCER_TOKEN,
  consumerToken: process.env.VARSNAP_CONSUMER_TOKEN,
});

require('../static/js/app.js');

context('Varsnap', function() {
  this.timeout(30 * 1000);
  it('runs with production', async function() {
    const status = await varsnap.runTests();
    expect(status).to.be.true;
  });
});
