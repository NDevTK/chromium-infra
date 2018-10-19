'use strict';

/**
 * `<mr-edit-field>`
 *
 * A single edit input for a fieldDef + the values of the field.
 *
 */
class MrEditField extends Polymer.Element {
  static get is() {
    return 'mr-edit-field';
  }

  static get properties() {
    return {
      // TODO(zhangtiff): Redesign this a bit so we don't need two separate
      // ways of specifying "type" for a field. Right now, "type" is mapped to
      // the Monorail custom field types whereas "acType" includes additional
      // data types such as components, and labels.
      // String specifying what kind of autocomplete to add to this field.
      acType: String,
      delimiter: {
        type: String,
        value: ',',
      },
      type: String,
      multi: {
        type: Boolean,
        value: false,
      },
      name: String,
      initialValues: {
        type: Array,
        value: () => [],
      },
      // For enum fields, the possible options that you have. Each entry is a
      // label type with an additional optionName field added.
      options: {
        type: Array,
        value: () => [],
      },
      _acType: {
        type: String,
        computed: '_computeAcType(acType, type)',
      },
      _initialValue: {
        type: String,
        computed: '_computeInitialValue(initialValues)',
      },
      // Set to true if a field uses a standard input instead of any sort of
      // fancier edit type.
      _fieldIsBasic: {
        type: Boolean,
        computed: '_computeFieldIsBasic(type)',
        value: true,
      },
    };
  }

  focus() {
    if (this._fieldIsBasic) {
      this._getInput().focus();
    }
  }

  reset() {
    const input = this._getInput();
    if (this._fieldIsBasic) {
      input.value = this._initialValue;
    }
    if (this._fieldIsEnum(this.type)) {
      if (this.multi) {
        Polymer.dom(this.root).querySelectorAll('.enum-input').forEach(
          (checkbox) => {
            checkbox.checked = this._optionInValues(
              this.initialValues, checkbox.value);
          }
        );
      } else {
        const options = Array.from(input.querySelectorAll('option'));
        input.selectedIndex = options.findIndex((option) => {
          return this._computeIsSelected(this._initialValue,
            option.value);
        });
      }
    }
  }

  getValuesAdded() {
    if (!this.multi && !this.getValue().length) return [];
    return fltHelpers.arrayDifference(this.getValues(), this.initialValues,
      this._equalsIgnoreCase);
  }

  getValuesRemoved() {
    if (!this.multi && this.getValue().length > 0) return [];
    return fltHelpers.arrayDifference(this.initialValues, this.getValues(),
      this._equalsIgnoreCase);
  }

  getValues() {
    const val = this._getInput().value;
    if (this.multi) {
      if (this._fieldIsEnum(this.type)) {
        const checkboxes = Array.from(Polymer.dom(this.root).querySelectorAll(
          '.enum-input'));
        return checkboxes.filter((c) => c.checked).map((c) => c.value.trim());
      } else {
        let valueList = val.split(this.delimiter);
        valueList = valueList.map((s) => (s.trim()));
        valueList = valueList.filter((s) => (s.length > 0));
        return valueList;
      }
    }
    return [val.trim()];
  }

  getValue() {
    return this._getInput().value.trim();
  }

  _equalsIgnoreCase(a, b) {
    return a.toLowerCase() === b.toLowerCase();
  }

  // TODO(zhangtiff): We want to gradually make this list longer as we handle
  // all custom input cases.
  _computeFieldIsBasic(type) {
    return !(this._fieldIsEnum(type));
  }

  _fieldIsDate(type) {
    return type === fieldTypes.DATE_TYPE;
  }

  _fieldIsEnum(type) {
    return type === fieldTypes.ENUM_TYPE;
  }

  _fieldIsInt(type) {
    return type === fieldTypes.INT_TYPE;
  }

  _fieldIsStr(type) {
    return type === fieldTypes.STR_TYPE;
  }

  _fieldIsUser(type) {
    return type === fieldTypes.USER_TYPE;
  }

  _fieldIsUrl(type) {
    return type === fieldTypes.URL_TYPE;
  }

  _computeInitialValue(values) {
    return values.join(',');
  }

  _getInput() {
    if (this._fieldIsEnum(this.type) && !this.multi) {
      return Polymer.dom(this.root).querySelector('#editSelect');
    }
    return Polymer.dom(this.root).querySelector('#editInput');
  }

  _optionInValues(values, optionName) {
    return values.includes(optionName);
  }

  _computeIsSelected(initialValue, optionName) {
    return initialValue === optionName;
  }

  _computeAcType(acType, type) {
    return acType || (type === fieldTypes.USER_TYPE ? 'member' : '');
  }

  _computeDomAutocomplete(acType) {
    if (acType) return 'off';
    return '';
  }
}

customElements.define(MrEditField.is, MrEditField);
