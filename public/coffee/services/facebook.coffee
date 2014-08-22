angular = require 'angular'

angular.module('app').factory 'Facebook', ($window, $q, Config) ->
  deferred = $q.defer()

  $window.fbAsyncInit = ->
    FB.init
      appId: Config.facebook.app_id
      cookie: true
      version: 'v2.1'

    deferred.resolve FB

  `(function(d, s, id) {
    var js, fjs = d.getElementsByTagName(s)[0];
    if (d.getElementById(id)) return;
    js = d.createElement(s); js.id = id;
    js.src = "//connect.facebook.net/en_US/sdk.js";
    fjs.parentNode.insertBefore(js, fjs);
  }(document, 'script', 'facebook-jssdk'))`

  deferred.promise