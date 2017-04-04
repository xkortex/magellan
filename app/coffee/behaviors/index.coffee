
# Marionette Behavior Manifest
module.exports =
  BindInputs:         require 'hn_behaviors/lib/bindInputs'
  CancelButton:       require './cancelButton'
  ClickableRelations: require './clickableRelations'
  Confirmations:      require './confirmations'
  CopyToClipboard:    require 'hn_behaviors/lib/copyToClipboard'
  Flashes:            require 'hn_behaviors/lib/flashes'
  FormSerialize:      require 'hn_behaviors/lib/formSerialize'
  ModelEvents:        require 'hn_behaviors/lib/modelEvents'
  SelectableChild:    require './selectableChild'
  SubmitButton:       require 'hn_behaviors/lib/submitButton'
  Tooltips:           require 'hn_behaviors/lib/tooltips'
  SortableList:       require './sortableList'
  SortableChild:      require './sortableChild'

