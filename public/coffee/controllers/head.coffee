angular = require 'angular'

angular.module('app').controller 'HeadCtrl', ($scope, $rootScope) ->
  $rootScope.$on 'titleChanged', (event, title) ->
    $scope.title = if title then "#{title} | maji.moe" else 'maji.moe'