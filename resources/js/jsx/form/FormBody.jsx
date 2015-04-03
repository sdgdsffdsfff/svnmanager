define([
'kernel',
'react',
'components/form/FormText',
'components/form/FormSelect'
],
function( core, React, FormText, FormSelect ){

    var FormBody = React.createClass({
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
                <form className="form-black">
                    {this.props.fields.map(function (field, i) {
                        if (!field.name) return;

                        var input,
                            label = null,
                            type = typeof field.type == 'undefined' ? 'text' : field.type;

                        if (this.props.useLabel) {
                            label = <label className={this.props.inline ? 'col-sm-3 control-label' : ''}>{field.label || field.placeholder}:</label>;
                        }

                        if (/text|radio|hidden|file|email/.test(type)) {
                            input = <FormText {...field}
                                parent={this}
                                overload={this.props.overload}
                                ref={field.name}
                                useLabel={this.props.useLabel}
                            />;
                            if( type == 'hidden' ){
                                label = false;
                            }
                        } else if (type == 'select') {
                            input = <FormSelect {...field}
                                parent={this}
                                overload={this.props.overload}
                                ref={field.name}
                                useLabel={this.props.useLabel}
                            />;
                        }

                        if (this.props.useLabel && this.props.inline) {
                            input = <div className="col-sm-9">{input}</div>;
                        }

                        return (
                            <div className={label === false ? 'hidden' : 'form-group'} key={i}>
                                {label}
                                {input}
                            </div>
                        )
                    }, this)}
                </form>
            )
        }
    });

    return FormBody;
});