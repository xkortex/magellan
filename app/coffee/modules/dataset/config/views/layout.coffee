FacetLayout = require '../facet_views/layout'
RuleLayout  = require '../knowledge_rule_views/layout'

# # # # #

class ConfigLayoutView extends require 'hn_views/lib/nav'
  className: 'container-fluid'
  # template: require './templates/layout'

  navItems: [
    { icon: 'fa-list',            text: 'Facets',             trigger: 'facets' }
    { icon: 'fa-university',      text: 'Knowledge Rules',    trigger: 'knowledge', default: true }
    { icon: 'fa-window-maximize', text: 'Viewer Rules',       trigger: 'viewer' }
  ]

  navEvents:
    'facets':     'facetConfig'
    'knowledge':  'knowledgeConfig'
    'viewer':     'viewerConfig'

  # navOptions:
  #   stacked: true

  facetConfig: ->
    @model.fetchFacets().then (facetCollection) =>
      @contentRegion.show new FacetLayout({ collection: facetCollection })

  knowledgeConfig: ->
    @model.fetchKnowledgeRules().then (ruleCollection) =>
      @contentRegion.show new RuleLayout({ collection: ruleCollection })

  viewerConfig: ->
    console.log 'VIEWER CONFIG'

# # # # #

module.exports = ConfigLayoutView