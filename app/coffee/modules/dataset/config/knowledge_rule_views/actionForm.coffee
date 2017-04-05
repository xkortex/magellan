
class ActionTypeForm extends Mn.LayoutView
  className: 'row'

  templateMap:
    block:    require './templates/action_type_block'
    replace:  require './templates/action_type_replace'

  getTemplate: ->
    return @templateMap[@options.actionType]

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
    { action: 'block',    icon: 'fa-hand-stop-o', text: 'Blocking', default: true }
    { action: 'replace',  icon: 'fa-quote-right', text: 'Replace' }
  ]

  templateHelpers: ->
    return { isNew: @options.isNew, availableActions: @availableActions }

  onRender: ->
    Backbone.Syphon.deserialize( @, @model.attributes )
    @renderDefaultTypeForm()

  actionSelected: (e) ->
    el = $(e.currentTarget)
    actionType = el.data('action')
    el.addClass('active').siblings('.btn').removeClass('active')
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

    else
      console.log @model.get('action')
      console.log 'SHOW ACTION TYPE'
      console.log @model.attributes

# # # # #

module.exports = ActionForm