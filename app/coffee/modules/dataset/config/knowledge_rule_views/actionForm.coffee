
class ActionTypeForm extends Mn.LayoutView
  className: 'row'

  templateMap:
    block:    require './templates/action_type_block'
    replace:  require './templates/action_type_replace'
    literal:  require './templates/action_type_literal'
    clone:    require './templates/action_type_clone'
    index_from_split: require './templates/action_type_index_from_split'

  initialize: (options) ->
    @model.set('action', @options.actionType)

  templateHelpers: ->
    return { sourceOptions: @options.sourceOptions }

  getTemplate: ->
    return @templateMap[@options.actionType]

  onRender: ->
    Backbone.Syphon.deserialize( @, @model.attributes )

# # # # #

class ActionForm extends Mn.LayoutView
  className: 'row'
  template: require './templates/action_form'

  ui:
    actionSelect: '[data-select=action]'

  events:
    'click @ui.actionSelect': 'actionSelected'

  regions:
    actionTypeRegion: '[data-region=action-type]'

  availableActions: [
    { action: 'literal',          icon: 'fa-quote-right',   text: 'Literal', default: true }
    { action: 'block',            icon: 'fa-hand-stop-o',   text: 'Blocking' }
    { action: 'replace',          icon: 'fa-strikethrough', text: 'Replace' }
    { action: 'clone',            icon: 'fa-copy',          text: 'Clone' }
    { action: 'index_from_split', icon: 'fa-code-fork',     text: 'Split and Index' }
  ]

  # TODO - Format Lowercase
  # TODO - Format Uppercase

  templateHelpers: ->
    return { isNew: @options.isNew, availableActions: @availableActions }

  onRender: ->
    @renderDefaultTypeForm()

  actionSelected: (e) ->
    el = $(e.currentTarget)
    actionType = el.data('action')
    el.addClass('active').siblings('.btn').removeClass('active')
    el.blur()
    @showActionTypeForm(actionType)

  showActionTypeForm: (actionType) ->
    @actionTypeRegion.show new ActionTypeForm({ model: @model, actionType: actionType, sourceOptions: @options.sourceOptions })

  # TODO - show the correct form for EDITING definitions that already know their actions
  renderDefaultTypeForm: ->

    # Shows default if @options.isNew
    if @options.isNew

      # Isolates default actionType
      defaultAction = _.findWhere(@availableActions, { default: true })
      @showActionTypeForm(defaultAction.action)

    # Shows the correct actionTypeForm view while editing
    else
      @showActionTypeForm(@model.get('action'))

# # # # #

module.exports = ActionForm
