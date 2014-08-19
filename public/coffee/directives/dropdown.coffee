angular = require 'angular'

angular.module('app').directive 'dropdown', ->
  body = angular.element document.body

  restrict: 'A'
  scope:
    dropdown: '='
  link: (scope, elem, attrs) ->
    children = elem.children()
    menu = children.eq(1)

    children.eq(0).on 'click', (e) ->
      e.preventDefault()
      e.stopPropagation()

      menu.toggleClass 'active'

    body.on 'click', ->
      menu.removeClass 'active'