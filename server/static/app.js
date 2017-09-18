var pendingRequest = undefined;

function showPost(postID) {
  $.getJSON(
    "/postdata/" + postID,
    function processPostResult(post) {
      clearResults();
      addResults([post]);
    }
  );
}

function updateResults(offset) {
  var query = $("#query").val();
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
      clearResults();
      saveQuery(query, data);
      updateURL(query);
      addResults(data);
      window.scrollTo(0, 0);
    }
  );
}

function clearResults() {
  $("#results").html("");
}

function saveQuery(query, data) {
  $("#data").html("");
  $("#data").append('<input type="hidden" id="query" value="' + query + '">');
  $("#data").append('<input type="hidden" id="paginateCount" value="' + data.data.length + '">');
  $("#data").append('<input type="hidden" id="offset" value="' + data.offset + '">');
  $("#data").append('<input type="hidden" id="totalResults" value="' + data.totalResults + '">');
}

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
}

function addResults(data) {
  for (var x=0; x<data.data.length; x++) {
    var post = data.data[x];
    addResult(post);
  }
  if (data.data.length + data.offset < data.totalResults) {
    var paginateHTML = '<a href="javascript:paginateNext()">';
    paginateHTML += 'Next Page <span class="glyphicon glyphicon-menu-right" aria-hidden="true"></span>';
    paginateHTML += '</a>';
    $("#results").append(paginateHTML);
  }
  $('img.result-img').lazyload({
    effect: "fadeIn",
    threshold: 1000,
    skip_invisible: true
  });
}

function paginateNext() {
  var offset = parseInt($("#offset").val(), 10);
  offset += parseInt($("#paginateCount").val(), 10);
  updateResults(offset);
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
  $("#results").append(postHTML);
}

function stats() {
  $.getJSON(
    "/stats.json",
    function processStats(data) {
      var line = "Currently indexing " + data.postCount + " posts";
      $("#indexStat").text(line);
    }
  );
}

function getParameterByName(name) {
    var url = window.location.href;
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

$(function() {
  var query = getParameterByName('query');
  if (query !== undefined && query !== '') {
    $("#query").val(query);
  }
  $("#query").on('input', function(){updateResults()});
  var urlPath = window.location.pathname.split('/');
  if (urlPath[1] === 'post') {
    showPost(urlPath[2]);
  } else {
    updateResults();
  }
  stats();
});
