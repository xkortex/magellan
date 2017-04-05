
class ConstraintForm extends Mn.LayoutView
  className: 'row'
  template: require './templates/constraint_form'

  # Defines operationOptions
  # TODO - regex match, Starts With (case sensitive), Ends With (case sensitive), Exact Match (case sensitive)
  operationOptions: [
    ['exact_match',             'Exact Match']
    ['starts_with',             'Starts With']
    ['ends_with',               'Ends With']
    ['contains',                'Contains']
    ['contains_case_sensitive', 'Contains (case sensitive)']
  ]

  templateHelpers: ->
    return helpers =
      sourceOptions:    @options.sourceOptions
      operationOptions: @operationOptions

  onRender: ->
    Backbone.Syphon.deserialize( @, @model.attributes )

# # # # #

module.exports = ConstraintForm