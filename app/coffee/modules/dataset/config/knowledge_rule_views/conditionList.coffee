
ConditionViewer = require './conditionViewer'
ConditionForm = require './conditionForm'

# # # # #

# Swaps models at supplied indicies & resets collection
# TODO - rather than swapping indicies we should just PUSH the list accordingly.
swapIndicies = (collection, oldIndex, newIndex) ->

  # Swaps model indicies
  collection.at(oldIndex).set('order', newIndex)
  collection.at(newIndex).set('order', oldIndex)

# # # # #

class ConditionEmpty extends Mn.LayoutView
  template: require './templates/condition_empty'
  className: 'list-group-item list-group-item-warning'

# # # # #

class ConditionChild extends Mn.LayoutView
  template: require './templates/condition_child'
  className: 'list-group-item'

  behaviors:
    SelectableChild: {}

  events:
    'sortable:end': 'onSortableEnd'

  modelEvents:
    'change:order': 'onOrderChange'

  onSortableEnd: (e, ev) ->
    # TODO - return if ev.oldIndex / newIndex == undefined
    # TODO - don't SWAP indicies. Rather, we should INSERT AT INDEX
    swapIndicies(@model.collection, ev.oldIndex, ev.newIndex)

  onOrderChange: ->
    # @model.save()

# # # # #

class ConditionList extends Mn.CompositeView
  className: 'list-group'
  template: require './templates/condition_list'
  childView: ConditionChild
  emptyView: ConditionEmpty

  onRender: ->

    # Initializes Sortable container
    Sortable.create @el,
      handle: '.handle'
      animation: 250
      onEnd: (e) => $(e.item).trigger('sortable:end', e)

# # # # #

class ConditionLayout extends Mn.LayoutView
  className: 'row'
  template: require './templates/condition_list_layout'

  # TODO - do we need these events?
  collectionEvents:
    'add':    'render'
    'remove': 'render'

  regions:
    listRegion:   '[data-region=list]'
    detailRegion: '[data-region=detail]'

  onRender: ->
    @showConditionList()

  showConditionList: ->
    listView = new ConditionList({ collection: @collection })
    listView.on 'childview:selected', (view) => @showConditionViewer(view.model)
    @listRegion.show(listView)
    @collection.at(0)?.trigger('selected')

  showConditionForm: (model) ->
    conditionForm = new ConditionForm({ model: model })
    conditionForm.on 'cancel', (view) => @showConditionViewer(view.model)
    conditionForm.on 'submit', (view) =>
      console.log 'SUBIMTTED'
      console.log view.model
      @showConditionList()
      # @showConditionViewer(view.model)

    @detailRegion.show conditionForm

  showConditionViewer: (model) ->
    conditionViewer = new ConditionViewer({ model: model })
    conditionViewer.on 'edit', (view) => @showConditionForm(view.model)
    @detailRegion.show conditionViewer

# # # # #

module.exports = ConditionLayout