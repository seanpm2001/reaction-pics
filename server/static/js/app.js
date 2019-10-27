var process = require('process');

var $ = require('jquery');
window.jQuery = $; // Hack to get jquery-lazyload to bind to jQuery
require('jquery-lazyload');
var varsnap = require('varsnap');

varsnap.config = {
  varsnap: process.env.VARSNAP,
  env: process.env.ENVIRONMENT,
  producerToken: process.env.VARSNAP_PRODUCER_TOKEN,
  consumerToken: process.env.VARSNAP_CONSUMER_TOKEN,
};

var pendingRequest = undefined;

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
  var query = $("#query").val();
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
updateResults = varsnap(updateResults);

function saveQuery(query, data) {
  var dataHTML = '';
  dataHTML += '<input type="hidden" id="query" value="' + query + '">';
  dataHTML += '<input type="hidden" id="paginateCount" value="' + data.data.length + '">';
  dataHTML += '<input type="hidden" id="offset" value="' + data.offset + '">';
  dataHTML += '<input type="hidden" id="totalResults" value="' + data.totalResults + '">';
  $("#data").html(dataHTML);
  return dataHTML;
}
saveQuery = varsnap(saveQuery);

function updateURL(query) {
  var url = "/";
  if (query !== undefined && query !== "") {
    url += "?query=" + query;
  }
  var urlPath = window.location.pathname.split('/');
  if (urlPath[1] === 'post') {
    history.pushState({}, "Reaction Pics", url);
  } else {
    history.replaceState({}, "Reaction Pics", url);
  }
  return url;
}
updateURL = varsnap(updateURL);

function addResults(data) {
  var resultHTML = '';
  for (var x=0; x<data.data.length; x++) {
    var post = data.data[x];
    resultHTML += addResult(post);
  }
  if (data.data.length + data.offset < data.totalResults) {
    resultHTML += '<a href="javascript:paginateNext()">';
    resultHTML += 'Next Page <span class="glyphicon glyphicon-menu-right" aria-hidden="true"></span>';
    resultHTML += '</a>';
  }
  setResults(resultHTML);
  $('img.result-img').lazyload({
    effect: "fadeIn",
    threshold: 1000,
    skip_invisible: true
  });
  return resultHTML;
}
addResults = varsnap(addResults);

function paginateNext() {
  var offset = parseInt($("#offset").val(), 10);
  offset += parseInt($("#paginateCount").val(), 10);
  updateResults(getQuery(), offset);
}

function addResult(postData) {
  var postHTML = '<div class="result">';
  postHTML += '<h2>';
  if (postData.url) postHTML += '<a href="' + postData.internalURL + '">';
  postHTML += postData.title;
  if (postData.url) postHTML += '</a></h2>';
  if (postData.image) postHTML += '<p><img data-original="' + postData.image + '" class="result-img" /></p>';
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
      var line = "Currently indexing " + data.postCount + " posts";
      $("#indexStat").text(line);
    }
  );
}
stats = varsnap(stats);

function getParameterByName(url, name) {
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}
getParameterByName = varsnap(getParameterByName);

$(function() {
  var query = getParameterByName(window.location.href, 'query');
  if (query !== undefined && query !== '') {
    $("#query").val(query);
  }
  $("#query").on('input', function(){updateResults(getQuery())});
  var urlPath = window.location.pathname.split('/');
  if (urlPath[1] === 'post') {
    showPost(urlPath[2]);
  } else {
    updateResults(getQuery());
  }
  stats();
});
