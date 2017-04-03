
class ApplyRulesView extends Mn.LayoutView
  template: require './templates/apply_rules'
  className: 'card card-block'

  behaviors:
    CancelButton: {}
    # SubmitButton: {}

    Confirmations:
      message:      'Are you sure you want to apply rules to this dataset?'
      confirmIcon:  'fa-check-circle-o'
      confirmText:  'Apply Rules'
      confirmCss:   'btn-success'

    Flashes:
      success:
        message:  'Successfully applied rules to dataset.'

  ui:
    confirmationTrigger:  '[data-click=submit]'

  onCancel: ->
    @trigger 'cancel'

  onConfirmed: ->
    # @disableSubmit()
    @disableCancel()

    @model.fetchDatapoints().then (datapoints) =>

      @model.fetchKnowledgeRules().then (rules) =>

        console.log 'APPLYING RULES...'

        rules.applyRules(datapoints).then () =>
          console.log 'RULES APPLIED!'

          console.log 'regenerateFacets...'
          @model.regenerateFacets().then () =>

            console.log 'regenerateFacets DONE!'

            # Shows success message
            @flashSuccess()

            # Triggers 'success' event
            @trigger('success')

# # # # #

module.exports = ApplyRulesView
