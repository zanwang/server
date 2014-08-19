angular = require 'angular'

angular.module('app').controller 'SidebarCtrl', ($scope, $rootScope, $state) ->
  $scope.state = $state.current.name

  $scope.links = [
    { name : 'Domains', state : 'app.domains' }
    { name : 'Settings', state : 'app.settings' }
    { name : 'Support', state : 'app.support' }
  ]

  $rootScope.$on '$stateChangeSuccess', (event, toState, toParams, fromState, fromParams) ->
    $scope.state = toState.name