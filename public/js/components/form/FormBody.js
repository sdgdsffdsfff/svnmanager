define([
'kernel',
'react',
'components/form/FormText',
'components/form/FormSelect'
],
function( core, React, FormText, FormSelect ){

    var FormBody = React.createClass({displayName: "FormBody",
        getRef: function( name ){
            if( name in this.refs ){
                return this.refs[name]
            }
            if( typeof name === 'undefined' ){
                return this.refs;
            }
            return null;
        },
        render: function(){
            return (
                React.createElement("form", {className: "form-black"}, 
                    this.props.fields.map(function (field, i) {
                        if (!field.name) return;

                        var input,
                            label = null,
                            type = typeof field.type == 'undefined' ? 'text' : field.type;

                        if (this.props.useLabel) {
                            label = React.createElement("label", {className: this.props.inline ? 'col-sm-3 control-label' : ''}, field.label || field.placeholder, ":");
                        }

                        if (/text|radio|hidden|file|email/.test(type)) {
                            input = React.createElement(FormText, React.__spread({},  field, 
                                {parent: this, 
                                overload: this.props.overload, 
                                ref: field.name, 
                                useLabel: this.props.useLabel})
                            );
                            if( type == 'hidden' ){
                                label = false;
                            }
                        } else if (type == 'select') {
                            input = React.createElement(FormSelect, React.__spread({},  field, 
                                {parent: this, 
                                overload: this.props.overload, 
                                ref: field.name, 
                                useLabel: this.props.useLabel})
                            );
                        }

                        if (this.props.useLabel && this.props.inline) {
                            input = React.createElement("div", {className: "col-sm-9"}, input);
                        }

                        return (
                            React.createElement("div", {className: label === false ? 'hidden' : 'form-group', key: i}, 
                                label, 
                                input
                            )
                        )
                    }, this)
                )
            )
        }
    });

    return FormBody;
});