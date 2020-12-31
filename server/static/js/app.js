const process = require('process');

const $ = require('jquery');
const LazyLoad = require('vanilla-lazyload');
const varsnap = require('varsnap');

varsnap.updateConfig({
  varsnap: process.env.VARSNAP,
  env: process.env.ENVIRONMENT,
  producerToken: process.env.VARSNAP_PRODUCER_TOKEN,
  consumerToken: process.env.VARSNAP_CONSUMER_TOKEN,
});

const lazyLoadInstance = new LazyLoad({});
let pendingRequest = undefined;

function showPost(postID) {
  $.getJSON(
    "/postdata/" + postID,
    function processPostResult(data) {
      setResults("");
      addResults(data);
    }
  );
}

function getQuery() {
  const query = $("#query").val();
  return query;
}

function setResults(html) {
  $("#results").html(html);
}

function updateResults(query, offset) {
  if (pendingRequest) {
    pendingRequest.abort();
  }
  pendingRequest = $.getJSON(
    "/search",
    {
      query: query,
      offset: offset,
    },
    function processQueryResult(data) {
      setResults("");
      saveQuery(query, data);
      updateURL(query);
      addResults(data);
      window.scrollTo(0, 0);
    }
  );
  return pendingRequest;
}
// Cannot serialize and compare jquery request
// updateResults = varsnap(updateResults);

function saveQuery(query, data) {
  let dataHTML = '';
  dataHTML += '<input type="hidden" id="query" value="' + query + '">';
  dataHTML += '<input type="hidden" id="paginateCount" value="' + data.data.length + '">';
  dataHTML += '<input type="hidden" id="offset" value="' + data.offset + '">';
  dataHTML += '<input type="hidden" id="totalResults" value="' + data.totalResults + '">';
  $("#data").html(dataHTML);
  return dataHTML;
}
saveQuery = varsnap(saveQuery);

function updateURL(query) {
  let url = "/";
  if (query !== undefined && query !== "") {
    url += "?query=" + query;
  }
  const urlPath = window.location.pathname.split('/');
  if (urlPath[1] === 'post') {
    history.pushState({}, "Reaction Pics", url);
  } else {
    history.replaceState({}, "Reaction Pics", url);
  }
  return url;
}
// Security issue when running in headless browser
// updateURL = varsnap(updateURL);

function addResults(data) {
  let resultHTML = '';
  for (let x=0; x<data.data.length; x++) {
    const post = data.data[x];
    resultHTML += addResult(post);
  }
  if (data.data.length + data.offset < data.totalResults) {
    resultHTML += '<a class="btn btn-primary" href="#" id="paginateNext">';
    resultHTML += 'Next Page <span class="glyphicon glyphicon-menu-right" aria-hidden="true"></span>';
    resultHTML += '</a>';
  }
  setResults(resultHTML);
  $("#paginateNext").click(paginateNext);
  lazyLoadInstance.update();
  return resultHTML;
}
addResults = varsnap(addResults);

function paginateNext() {
  let offset = parseInt($("#offset").val(), 10);
  offset += parseInt($("#paginateCount").val(), 10);
  updateResults(getQuery(), offset);
}

function addResult(postData) {
  let postHTML = '<div class="result">';
  postHTML += '<h2>';
  if (postData.url) postHTML += '<a href="' + postData.internalURL + '">';
  postHTML += postData.title;
  if (postData.url) postHTML += '</a></h2>';
  if (postData.image) postHTML += '<p><img data-src="' + postData.image + '" class="result-img lazy" /></p>';
  if (postData.likes) {
      postHTML += '<p><a href="#" id="likes" class="btn btn-success disabled">';
      postHTML += postData.likes + ' <span class="glyphicon glyphicon-thumbs-up" aria-hidden="true"></span>';
      postHTML += '</a></p>';
      postHTML += '<p><a href="' + postData.url + '">Original</a></p>';
  }
  postHTML += '</div>';
  return postHTML;
}
addResult = varsnap(addResult);

function stats() {
  return $.getJSON(
    "/stats.json",
    function processStats(data) {
      const line = "Currently indexing " + data.postCount + " posts";
      $("#indexStat").text(line);
    }
  );
}
// Cannot serialize and compare jquery request
// stats = varsnap(stats);

function getParameterByName(url, name) {
    name = name.replace(/[\[\]]/g, "\\$&");
    const regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)");
    const results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}
getParameterByName = varsnap(getParameterByName);

$(function() {
  const query = getParameterByName(window.location.href, 'query');
  if (query !== undefined && query !== '') {
    $("#query").val(query);
  }
  $("#query").on('input', function(){updateResults(getQuery())});
  const urlPath = window.location.pathname.split('/');
  if (urlPath[1] === 'post') {
    showPost(urlPath[2]);
  } else {
    updateResults(getQuery());
  }
  stats();
});
