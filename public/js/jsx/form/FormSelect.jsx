define([
'kernel',
'react'
],
function( core, React ){
    var FormSelect = React.createClass({
        getDefaultProps: function () {
            return {
                values: [],
                change: $.noop
            }
        },
        setValues: function (values) {
            var html = values.map(function( t ){
                return '<option value="'+t.value+'">'+t.text+'</option>'
            }).join('');
            $(this.getDOMNode()).html(html);
        },
        getValue: function () {
            return this.getDOMNode().value;
        },
        restore: function(){
            this.setValue(this.props.value);
        },
        setValue: function (value) {
            value = 'undefined' == typeof value ? 'null' : value;
            $(this.getDOMNode()).val(value);
        },
        render: function () {
            return (
                <select
                    name={this.props.name}
                    onChange={this.props.change.bind(this.props.overload, this)}
                    className="form-control input-sm"
                    defaultValue={this.props.value}>
                    <option value="null">please select</option>
                    {this.props.values.map(function (result) {
                        return <option value={result.value}>{result.text}</option>
                    }, this)}
                </select>
            )
        }
    });

    return FormSelect;
});