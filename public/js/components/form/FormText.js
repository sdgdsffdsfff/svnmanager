define([
'kernel',
'react'
],
function( core, React ){
    var FormText = React.createClass({displayName: "FormText",
        getDefaultProps: function () {
            return {
                type: 'text',
                name: '',
                required: false,
                value: '',
                placeholder: '',
                keyup: $.noop
            }
        },
        clear: function () {
            this.getDOMNode().value = '';
        },
        isEmpty: function () {
            return this.getValue() == ''
        },
        restore: function () {
            this.getDOMNode().value = this.props.value;
        },
        setValue: function( value ){
            this.getDOMNode().value = value;
        },
        getValue: function () {
            return this.getDOMNode().value.trim()
        },
        focus: function () {
            this.getDOMNode().focus();
        },
        render: function () {
            return (
                React.createElement("input", {
                    type: this.props.type, 
                    name: this.props.name, 
                    className: "form-control", 
                    onKeyup: this.props.keyup.bind(this.props.overload, this), 
                    placeholder: this.props.useLabel ? '' : this.props.placeholder, 
                    defaultValue: this.props.value})
            )
        }
    });

    return FormText;
});